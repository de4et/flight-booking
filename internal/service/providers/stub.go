package providers

import (
	"context"
	"flight-booking/internal/model/sro"
	"flight-booking/internal/model/trip"
	"time"
)

type StubGDS struct {
	secs int
}

func NewStubGDS(secs int) *StubGDS {
	return &StubGDS{
		secs: secs,
	}
}

func (gds *StubGDS) Search(ctx context.Context, sro sro.SRO) (*trip.Trips, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	select {
	case <-time.After(time.Duration(gds.secs) * time.Second):
		ts := trip.NewTrips()
		ts.AddTrip(trip.Trip{})
		return ts, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (gds *StubGDS) GetAvailability() bool {
	return true
}
