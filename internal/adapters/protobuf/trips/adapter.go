package trips

import (
	"time"

	sroAdapter "flight-booking/internal/adapters/protobuf/sro"

	"flight-booking/internal/model/sro"
	"flight-booking/internal/model/trip"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// TripToProto converts a domain Trip to protobuf Trip
func TripToProto(t *trip.Trip) *Trip {
	if t == nil {
		return nil
	}

	return &Trip{
		Rid:      t.RID,
		Tid:      t.TID,
		Sid:      t.SID,
		CacheId:  t.CacheID,
		Provider: providerToProto(t.Provider),
		Segments: segmentsToProto(t.Segments),
		Prices:   pricesToProto(t.Prices),
		Rules:    rulesToProto(t.Rules),
		Metadata: metadataToProto(t.Metadata),
		Booking:  bookingToProto(t.Booking),
		Sro:      sroToProto(t.SRO), // Use the SRO adapter
	}
}

// ProtoToTrip converts a protobuf Trip to domain Trip
func ProtoToTrip(p *Trip) *trip.Trip {
	if p == nil {
		return nil
	}

	return &trip.Trip{
		RID:      p.GetRid(),
		TID:      p.GetTid(),
		SID:      p.GetSid(),
		CacheID:  p.GetCacheId(),
		Provider: protoToProvider(p.GetProvider()),
		Segments: protoToSegments(p.GetSegments()),
		Prices:   protoToPrices(p.GetPrices()),
		Rules:    protoToRules(p.GetRules()),
		Metadata: protoToMetadata(p.GetMetadata()),
		Booking:  protoToBooking(p.GetBooking()),
		SRO:      protoToSRO(p.GetSro()), // Use the SRO adapter
	}
}

// Collection adapters (unchanged)
func TripsToProto(ts *trip.Trips) *Trips {
	if ts == nil {
		return nil
	}

	tripArray := ts.ToArray()
	protoTrips := make([]*Trip, 0, len(tripArray))

	for i := range tripArray {
		protoTrips = append(protoTrips, TripToProto(&tripArray[i]))
	}

	return &Trips{
		Trips: protoTrips,
	}
}

func ProtoToTrips(collection *Trips) *trip.Trips {
	if collection == nil {
		return trip.NewTrips()
	}

	trips := trip.NewTrips()
	for _, protoTrip := range collection.GetTrips() {
		if domainTrip := ProtoToTrip(protoTrip); domainTrip != nil {
			trips.AddTrip(*domainTrip)
		}
	}
	return trips
}

// Helper conversion functions (all the existing ones remain the same)
func providerToProto(p trip.Provider) *Provider {
	return &Provider{
		Name:              p.Name,
		Gds:               p.GDS,
		GdsServer:         p.GDSServer,
		OfficeId:          p.OfficeID,
		ValidatingCarrier: p.ValidatingCarrier,
	}
}

func protoToProvider(p *Provider) trip.Provider {
	if p == nil {
		return trip.Provider{}
	}
	return trip.Provider{
		Name:              p.GetName(),
		GDS:               p.GetGds(),
		GDSServer:         p.GetGdsServer(),
		OfficeID:          p.GetOfficeId(),
		ValidatingCarrier: p.GetValidatingCarrier(),
	}
}

func segmentsToProto(segments []trip.TripSegment) []*TripSegment {
	protoSegments := make([]*TripSegment, 0, len(segments))
	for _, segment := range segments {
		protoSegments = append(protoSegments, &TripSegment{
			FlightNumber:     segment.FlightNumber,
			Carrier:          segment.Carrier,
			OperatingCarrier: segment.OperatingCarrier,
			Departure:        flightPointToProto(segment.Departure),
			Arrival:          flightPointToProto(segment.Arrival),
			DurationMinutes:  int32(segment.DurationMinutes),
			CabinClass:       string(segment.CabinClass),
			FareCode:         segment.FareCode,
			Baggage:          baggageToProto(segment.Baggage),
			Meal:             segment.Meal,
			Aircraft:         segment.Aircraft,
			StopTimeMinutes:  int32(segment.StopTimeMinutes),
			Direction:        int32(segment.Direction),
		})
	}
	return protoSegments
}

func protoToSegments(protoSegments []*TripSegment) []trip.TripSegment {
	segments := make([]trip.TripSegment, 0, len(protoSegments))
	for _, protoSegment := range protoSegments {
		if protoSegment != nil {
			segments = append(segments, trip.TripSegment{
				FlightNumber:     protoSegment.GetFlightNumber(),
				Carrier:          protoSegment.GetCarrier(),
				OperatingCarrier: protoSegment.GetOperatingCarrier(),
				Departure:        protoToFlightPoint(protoSegment.GetDeparture()),
				Arrival:          protoToFlightPoint(protoSegment.GetArrival()),
				DurationMinutes:  int(protoSegment.GetDurationMinutes()),
				CabinClass:       trip.TravelClass(protoSegment.GetCabinClass()),
				FareCode:         protoSegment.GetFareCode(),
				Baggage:          protoToBaggage(protoSegment.GetBaggage()),
				Meal:             protoSegment.GetMeal(),
				Aircraft:         protoSegment.GetAircraft(),
				StopTimeMinutes:  int(protoSegment.GetStopTimeMinutes()),
				Direction:        int(protoSegment.GetDirection()),
			})
		}
	}
	return segments
}

func flightPointToProto(fp trip.FlightPoint) *FlightPoint {
	return &FlightPoint{
		Airport:  fp.Airport,
		Terminal: fp.Terminal,
		Time:     timestamppb.New(fp.Time),
	}
}

func protoToFlightPoint(fp *FlightPoint) trip.FlightPoint {
	if fp == nil {
		return trip.FlightPoint{}
	}
	var t time.Time
	if fp.GetTime() != nil {
		t = fp.GetTime().AsTime()
	}
	return trip.FlightPoint{
		Airport:  fp.GetAirport(),
		Terminal: fp.GetTerminal(),
		Time:     t,
	}
}

func baggageToProto(b trip.BaggageInfo) *BaggageInfo {
	return &BaggageInfo{
		Pieces: int32(b.Pieces),
		Weight: int32(b.Weight),
		Type:   b.Type,
	}
}

func protoToBaggage(b *BaggageInfo) trip.BaggageInfo {
	if b == nil {
		return trip.BaggageInfo{}
	}
	return trip.BaggageInfo{
		Pieces: int(b.GetPieces()),
		Weight: int(b.GetWeight()),
		Type:   b.GetType(),
	}
}

func pricesToProto(p trip.TripPrices) *TripPrices {
	passengersPriceDetails := make(map[string]float64)
	for k, v := range p.PassengersPriceDetails {
		passengersPriceDetails[k] = v
	}

	return &TripPrices{
		Price:                  p.Price,
		SearchPrice:            p.SearchPrice,
		PriceFare:              p.PriceFare,
		ProviderServiceFee:     p.ProviderServiceFee,
		ProviderTaxesAmount:    p.ProviderTaxesAmount,
		ProviderCurrency:       p.ProviderCurrency,
		MinAllowablePrice:      p.MinAllowablePrice,
		TkpTax:                 p.TkpTax,
		SpecTax:                p.SpecTax,
		PricerInfo:             pricerInfoToProto(p.PricerInfo),
		BagsPrice:              p.BagsPrice,
		PassengersPriceDetails: passengersPriceDetails,
	}
}

func protoToPrices(p *TripPrices) trip.TripPrices {
	if p == nil {
		return trip.TripPrices{}
	}

	passengersPriceDetails := make(map[string]float64)
	for k, v := range p.GetPassengersPriceDetails() {
		passengersPriceDetails[k] = v
	}

	return trip.TripPrices{
		Price:                  p.GetPrice(),
		SearchPrice:            p.GetSearchPrice(),
		PriceFare:              p.GetPriceFare(),
		ProviderServiceFee:     p.GetProviderServiceFee(),
		ProviderTaxesAmount:    p.GetProviderTaxesAmount(),
		ProviderCurrency:       p.GetProviderCurrency(),
		MinAllowablePrice:      p.GetMinAllowablePrice(),
		TkpTax:                 p.GetTkpTax(),
		SpecTax:                p.GetSpecTax(),
		PricerInfo:             protoToPricerInfo(p.GetPricerInfo()),
		BagsPrice:              p.GetBagsPrice(),
		PassengersPriceDetails: passengersPriceDetails,
	}
}

func pricerInfoToProto(pi trip.PricerInfo) *PricerInfo {
	return &PricerInfo{
		Markup:              pi.Markup,
		Commission:          pi.Commission,
		PartnerAffiliateFee: pi.PartnerAffiliateFee,
		CashbackRate:        pi.CashbackRate,
	}
}

func protoToPricerInfo(pi *PricerInfo) trip.PricerInfo {
	if pi == nil {
		return trip.PricerInfo{}
	}
	return trip.PricerInfo{
		Markup:              pi.GetMarkup(),
		Commission:          pi.GetCommission(),
		PartnerAffiliateFee: pi.GetPartnerAffiliateFee(),
		CashbackRate:        pi.GetCashbackRate(),
	}
}

func rulesToProto(r trip.FareRules) *FareRules {
	return &FareRules{
		IsRefund:       r.IsRefund,
		IsExchangeable: r.IsExchangeable,
		RefundAmount:   r.RefundAmount,
		ExchangeFee:    r.ExchangeFee,
		Penalty:        r.Penalty,
	}
}

func protoToRules(r *FareRules) trip.FareRules {
	if r == nil {
		return trip.FareRules{}
	}
	return trip.FareRules{
		IsRefund:       r.GetIsRefund(),
		IsExchangeable: r.GetIsExchangeable(),
		RefundAmount:   r.GetRefundAmount(),
		ExchangeFee:    r.GetExchangeFee(),
		Penalty:        r.GetPenalty(),
	}
}

func metadataToProto(m trip.TripMetadata) *TripMetadata {
	return &TripMetadata{
		FlightType:         m.FlightType,
		IsVtrip:            m.IsVtrip,
		VtripComboId:       m.VtripComboID,
		RouteDuration:      int32(m.RouteDuration),
		NumTransfers:       int32(m.NumTransfers),
		HasBaggage:         m.HasBaggage,
		HasLuggage:         m.HasLuggage,
		FareFamily:         fareFamilyToProto(m.FareFamily),
		TariffType:         m.TariffType,
		AgeThreshold:       int32(m.AgeThreshold),
		IsVirtualInterline: m.IsVirtualInterline,
	}
}

func protoToMetadata(m *TripMetadata) trip.TripMetadata {
	if m == nil {
		return trip.TripMetadata{}
	}
	return trip.TripMetadata{
		FlightType:         m.GetFlightType(),
		IsVtrip:            m.GetIsVtrip(),
		VtripComboID:       m.GetVtripComboId(),
		RouteDuration:      int(m.GetRouteDuration()),
		NumTransfers:       int(m.GetNumTransfers()),
		HasBaggage:         m.GetHasBaggage(),
		HasLuggage:         m.GetHasLuggage(),
		FareFamily:         protoToFareFamily(m.GetFareFamily()),
		TariffType:         m.GetTariffType(),
		AgeThreshold:       int(m.GetAgeThreshold()),
		IsVirtualInterline: m.GetIsVirtualInterline(),
	}
}

func fareFamilyToProto(ff trip.FareFamily) *FareFamily {
	return &FareFamily{
		Type:          ff.Type,
		Name:          ff.Name,
		MarketingName: ff.MarketingName,
		HasFareFamily: ff.HasFareFamily,
	}
}

func protoToFareFamily(ff *FareFamily) trip.FareFamily {
	if ff == nil {
		return trip.FareFamily{}
	}
	return trip.FareFamily{
		Type:          ff.GetType(),
		Name:          ff.GetName(),
		MarketingName: ff.GetMarketingName(),
		HasFareFamily: ff.GetHasFareFamily(),
	}
}

func bookingToProto(b trip.TripBooking) *TripBooking {
	return &TripBooking{
		ExpiresAt:                     timeToProto(b.ExpiresAt),
		TicketingTimeLimit:            timeToProto(b.TicketingTimeLimit),
		ProviderRecommendationLimit:   timeToProto(b.ProviderRecommendationLimit),
		ProviderRecommendationCreated: timeToProto(b.ProviderRecommendationCreated),
		CountOfBlanks:                 int32(b.CountOfBlanks),
		BookingWithPartialDataAllowed: b.BookingWithPartialDataAllowed,
		BookingActualizationAllowed:   b.BookingActualizationAllowed,
	}
}

func protoToBooking(b *TripBooking) trip.TripBooking {
	if b == nil {
		return trip.TripBooking{}
	}
	return trip.TripBooking{
		ExpiresAt:                     protoToTime(b.GetExpiresAt()),
		TicketingTimeLimit:            protoToTime(b.GetTicketingTimeLimit()),
		ProviderRecommendationLimit:   protoToTime(b.GetProviderRecommendationLimit()),
		ProviderRecommendationCreated: protoToTime(b.GetProviderRecommendationCreated()),
		CountOfBlanks:                 int(b.GetCountOfBlanks()),
		BookingWithPartialDataAllowed: b.GetBookingWithPartialDataAllowed(),
		BookingActualizationAllowed:   b.GetBookingActualizationAllowed(),
	}
}

// Helper functions for time conversion
func timeToProto(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

func protoToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func sroToProto(s *sro.SRO) *sroAdapter.SRO {
	return sroAdapter.SROToProto(s)
}

func protoToSRO(p *sroAdapter.SRO) *sro.SRO {
	return sroAdapter.ProtoToSRO(p)
}
