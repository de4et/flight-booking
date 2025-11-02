package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/de4et/flight-booking/internal/logger"
	"github.com/de4et/flight-booking/internal/model/sro"
	"github.com/de4et/flight-booking/internal/model/trip"
)

var (
	ErrInvalidSRO = fmt.Errorf("invalid SRO")
	ErrNoCacheHit = errors.New("")
)

type provider interface {
	Search(context.Context, sro.SRO) (*trip.Trips, error)
	GetAvailability() bool
}

type cache interface {
	Get(context.Context, string) (*trip.Trips, error)
	Set(context.Context, string, *trip.Trips) error
}

type MultipleSearchService struct {
	cache     cache
	providers []provider
}

func NewMultipleSearchService(cache cache) *MultipleSearchService {
	return &MultipleSearchService{
		providers: make([]provider, 0),
		cache:     cache,
	}
}

func (svc *MultipleSearchService) AddProviderService(p provider) {
	svc.providers = append(svc.providers, p)
}

func (svc *MultipleSearchService) SearchByToken(ctx context.Context, token string) (*trip.Trips, error) {
	slog.DebugContext(ctx, "Starting searching token...")
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	ts, err := svc.cache.Get(ctx, token)
	if err == nil {
		slog.InfoContext(ctx, "Cache hit!")
		return ts, nil
	}

	if !errors.Is(err, ErrNoCacheHit) {
		slog.ErrorContext(ctx, "Failed calling cache", "error", err)
	}

	sro, err := sro.FromToken(token)
	if err != nil {
		return nil, ErrInvalidSRO
	}

	ctx = logger.WithContext(ctx, "sro.channeltoken", sro.ChannelToken)
	slog.DebugContext(ctx, "Sucessfully serialized sro from token")

	ts, err = svc.searchParallel(ctx, *sro)
	if err != nil {
		return nil, err
	}

	err = svc.cache.Set(ctx, token, ts)
	if err != nil {
		slog.DebugContext(ctx, "Couldn't set cache", "error", err)
	}
	return ts, nil
}

type searchResponse struct {
	tr  *trip.Trips
	err error
}

func (svc *MultipleSearchService) searchParallel(ctx context.Context, sro sro.SRO) (*trip.Trips, error) {
	ts := trip.NewTrips()
	outCh := make(chan searchResponse)
	wg := &sync.WaitGroup{}

	wg.Add(len(svc.providers))
	for i := range svc.providers {
		go svc.searchByProvider(ctx, wg, svc.providers[i], outCh, sro)
	}

	go func() {
		wg.Wait()
		close(outCh)
	}()

	for v := range outCh {
		if v.err != nil {
			continue
		}
		ts.Merge(v.tr)
	}

	return ts, nil
}

func (svc *MultipleSearchService) searchByProvider(ctx context.Context, wg *sync.WaitGroup, p provider, outCh chan searchResponse, sro sro.SRO) {
	defer wg.Done()

	resultCh := make(chan searchResponse)

	go func() {
		defer close(resultCh)
		v, err := p.Search(ctx, sro)
		resultCh <- searchResponse{v, err}
	}()

	select {
	case <-ctx.Done():
		outCh <- searchResponse{nil, ctx.Err()}
	case result := <-resultCh:
		outCh <- result
	}
}
