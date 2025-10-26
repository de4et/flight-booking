package providers

import (
	"context"

	"flight-booking/internal/model/sro"
	"flight-booking/internal/model/trip"
)

type StubGDS struct{}

func NewStubGDS() *StubGDS {
	return &StubGDS{}
}

func (gds *StubGDS) Search(ctx context.Context, sro sro.SRO) *trip.Trips {
	if ctx.Err() != nil {
		return nil
	}
	ts := trip.NewTrips()
	return ts
}

func (gds *StubGDS) GetAvailability() bool {
	return true
}
