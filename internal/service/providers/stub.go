package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/de4et/flight-booking/internal/model/sro"
	"github.com/de4et/flight-booking/internal/model/trip"
)

type StubGDS struct {
	secs int
}

func NewStubGDS(secs int) *StubGDS {
	return &StubGDS{
		secs: secs,
	}
}

func (gds *StubGDS) Search(ctx context.Context, s sro.SRO) (*trip.Trips, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	select {
	case <-time.After(time.Duration(gds.secs) * time.Second):
		ts := trip.NewTrips()
		for i := range 15 {
			ts.AddTrip(trip.Trip{
				RID:     "RIDDD",
				TID:     "TIDD",
				SID:     "SSIIIDDD",
				CacheID: fmt.Sprintf("best_cached_id_ever_%d", i),
				Provider: trip.Provider{
					Name:              "MyMan",
					GDS:               "luchiy",
					GDSServer:         "111.111.bs.3",
					OfficeID:          "tam-to",
					ValidatingCarrier: "chto?",
				},
				Segments: []trip.TripSegment{},
				Prices:   trip.TripPrices{},
				Rules:    trip.FareRules{},
				Metadata: trip.TripMetadata{},
				Booking:  trip.TripBooking{},
				SRO:      &s,
			})
		}
		return ts, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (gds *StubGDS) GetAvailability() bool {
	return true
}
