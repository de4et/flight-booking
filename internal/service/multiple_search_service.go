package service

import (
	"context"
	"fmt"
	"sync"

	"flight-booking/internal/model/sro"
	"flight-booking/internal/model/trip"
)

type provider interface {
	Search(context.Context, sro.SRO) *trip.Trips
	GetAvailability() bool
}

type MultipleSearchService struct {
	// cache SROCache
	providers []provider
}

func NewMultipleSearchService() *MultipleSearchService {
	return &MultipleSearchService{
		providers: make([]provider, 0),
	}
}

func (svc *MultipleSearchService) AddProviderService(p provider) {
	svc.providers = append(svc.providers, p)
}

func (svc *MultipleSearchService) SearchByToken(ctx context.Context, token string) (*trip.Trips, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	sro, err := sro.FromToken(token)
	if err != nil {
		return nil, fmt.Errorf("couldn't get SRO from token: %v", err)
	}

	ts := svc.searchParallel(ctx, *sro)
	return ts, nil
}

func (svc *MultipleSearchService) searchParallel(ctx context.Context, sro sro.SRO) *trip.Trips {
	ts := trip.NewTrips()
	outCh := make(chan *trip.Trips)
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
		ts.Merge(v)
	}

	return ts
}

func (svc *MultipleSearchService) searchByProvider(ctx context.Context, wg *sync.WaitGroup, p provider, outCh chan *trip.Trips, sro sro.SRO) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	case outCh <- p.Search(ctx, sro):
		return
	}
}
