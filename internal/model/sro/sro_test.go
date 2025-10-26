package sro_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"flight-booking/internal/model/sro"
)

func TestFromToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		want    *sro.SRO
		wantErr bool
	}{
		{
			name:  "Simple OW",
			token: "AKV40000OWE1000000091MOWLED20241015",
			want: &sro.SRO{
				Segments: []sro.Segment{
					{From: "MOW", To: "LED", Date: time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC)},
				},
				Passengers: sro.Passengers{ADT: 1, INS: true},
				Class:      sro.TravelClassE,
				Type:       sro.RouteTypeOW,
				ChannelToken: sro.ChannelToken{
					PartnerCode: "AKV4",
					SourceCode:  "0000",
				},
				Filters: sro.Filters{MaxStops: 9},
			},
			wantErr: false,
		},
		{
			name:  "Round Trip RT",
			token: "AKV40000RTE2000000010MOWLED20241015LEDMOW20241020",
			want: &sro.SRO{
				Segments: []sro.Segment{
					{From: "MOW", To: "LED", Date: time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC)},
					{From: "LED", To: "MOW", Date: time.Date(2024, 10, 20, 0, 0, 0, 0, time.UTC)},
				},
				Passengers: sro.Passengers{ADT: 2},
				Class:      sro.TravelClassE,
				Type:       sro.RouteTypeRT,
				ChannelToken: sro.ChannelToken{
					PartnerCode: "AKV4",
					SourceCode:  "0000",
				},
				Filters: sro.Filters{MaxStops: 1},
			},
			wantErr: false,
		},
		{
			name:  "Direct Only + No baggage",
			token: "AKV40000OWE1000001110MOWLED20241015",
			want: &sro.SRO{
				Segments: []sro.Segment{
					{From: "MOW", To: "LED", Date: time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC)},
				},
				Passengers: sro.Passengers{ADT: 1},
				Class:      sro.TravelClassE,
				Type:       sro.RouteTypeOW,
				ChannelToken: sro.ChannelToken{
					PartnerCode: "AKV4",
					SourceCode:  "0000",
				},
				Filters: sro.Filters{
					IsDirectOnly:    true,
					WithBaggageOnly: true,
					MaxStops:        1,
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid token (too short)",
			token:   "AKV40000OWE1",
			wantErr: true,
		},
		{
			name:    "Empty string",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sro.FromToken(tt.token)

			if (err != nil) != tt.wantErr {
				t.Fatalf("FromToken() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				assertEqualSRO(t, tt.want, got)
			}
		})
	}
}

func TestSRO_GetToken(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		token string
		want  string
	}{
		{
			name:  "simple",
			token: "AKV40000OWE1000001110MOWLED20241015",
			want:  "AKV40000OWE1000001110MOWLED20241015",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sro, err := sro.FromToken(tt.token)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := sro.GetToken()
			if tt.want != got {
				t.Errorf("GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func assertEqualSRO(t *testing.T, want, got *sro.SRO) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("SRO mismatch (-want +got):\n%s", diff)
	}
}
