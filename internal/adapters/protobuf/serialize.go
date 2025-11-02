package protobuf

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/de4et/flight-booking/internal/adapters/protobuf/trips"
	"github.com/de4et/flight-booking/internal/model/trip"
)

type TripsSerializer struct{}

func NewTripsSerializer() *TripsSerializer {
	return &TripsSerializer{}
}

func (s *TripsSerializer) SerializeTrips(ts *trip.Trips) ([]byte, error) {
	if ts == nil {
		return nil, fmt.Errorf("trips cannot be nil")
	}

	if ts.IsEmpty() {
		return []byte{}, nil
	}

	protoCollection := trips.TripsToProto(ts)
	if protoCollection == nil {
		return nil, fmt.Errorf("failed to convert trips to protobuf")
	}

	data, err := proto.Marshal(protoCollection)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trips: %w", err)
	}

	return data, nil
}

func (s *TripsSerializer) DeserializeTrips(data []byte) (*trip.Trips, error) {
	if len(data) == 0 {
		return trip.NewTrips(), nil
	}

	var protoTrips trips.Trips
	if err := proto.Unmarshal(data, &protoTrips); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trips: %w", err)
	}

	trips := trips.ProtoToTrips(&protoTrips)
	return trips, nil
}
