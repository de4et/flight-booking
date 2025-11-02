package trip

import (
	"flight-booking/internal/model/sro"
	"time"
)

type Trip struct {
	RID      string        `json:"rid"`
	TID      string        `json:"tid"`
	SID      string        `json:"sid"`
	CacheID  string        `json:"cacheId"`
	Provider Provider      `json:"provider"`
	Segments []TripSegment `json:"segments"`
	Prices   TripPrices    `json:"prices"`
	Rules    FareRules     `json:"rules"`
	Metadata TripMetadata  `json:"metadata"`
	Booking  TripBooking   `json:"booking"`
	SRO      *sro.SRO      `json:"sro"`
}

type TripSegment struct {
	FlightNumber     string      `json:"flightNumber"`
	Carrier          string      `json:"carrier"`
	OperatingCarrier string      `json:"operatingCarrier"`
	Departure        FlightPoint `json:"departure"`
	Arrival          FlightPoint `json:"arrival"`
	DurationMinutes  int         `json:"durationMinutes"`
	CabinClass       TravelClass `json:"cabinClass"`
	FareCode         string      `json:"fareCode"`
	Baggage          BaggageInfo `json:"baggage"`
	Meal             string      `json:"meal"`
	Aircraft         string      `json:"aircraft"`
	StopTimeMinutes  int         `json:"stopTimeMinutes"`
	Direction        int         `json:"direction"`
}

type TravelClass string

const (
	TravelClassE TravelClass = "E"
	TravelClassB TravelClass = "B"
	TravelClassF TravelClass = "F"
	TravelClassW TravelClass = "W"
)

type FlightPoint struct {
	Airport  string    `json:"airport"`
	Terminal string    `json:"terminal,omitempty"`
	Time     time.Time `json:"time"`
}

type BaggageInfo struct {
	Pieces int    `json:"pieces"`
	Weight int    `json:"weight"`
	Type   string `json:"type"`
}

type TripPrices struct {
	Price                  float64            `json:"price"`
	SearchPrice            float64            `json:"searchPrice"`
	PriceFare              float64            `json:"priceFare"`
	ProviderServiceFee     float64            `json:"providerServiceFee"`
	ProviderTaxesAmount    float64            `json:"providerTaxesAmount"`
	ProviderCurrency       string             `json:"providerCurrency"`
	MinAllowablePrice      float64            `json:"minAllowablePrice"`
	TkpTax                 float64            `json:"tkpTax"`
	SpecTax                float64            `json:"specTax"`
	PricerInfo             PricerInfo         `json:"pricerInfo"`
	BagsPrice              float64            `json:"bagsPrice"`
	PassengersPriceDetails map[string]float64 `json:"passengersPriceDetails"`
}

type PricerInfo struct {
	Markup              float64 `json:"markup"`
	Commission          float64 `json:"commission"`
	PartnerAffiliateFee float64 `json:"partnerAffiliateFee"`
	CashbackRate        float64 `json:"cashbackRate"`
}

type Provider struct {
	Name              string `json:"name"`
	GDS               string `json:"gds"`
	GDSServer         string `json:"gdsServer"`
	OfficeID          string `json:"officeId"`
	ValidatingCarrier string `json:"validatingCarrier"`
}

type FareRules struct {
	IsRefund       bool    `json:"isRefund"`
	IsExchangeable bool    `json:"isExchangeable"`
	RefundAmount   float64 `json:"refundAmount"`
	ExchangeFee    float64 `json:"exchangeFee"`
	Penalty        float64 `json:"penalty"`
}

type TripMetadata struct {
	FlightType         string     `json:"flightType"`
	IsVtrip            bool       `json:"isVtrip"`
	VtripComboID       string     `json:"vtripComboId,omitempty"`
	RouteDuration      int        `json:"routeDuration"`
	NumTransfers       int        `json:"numTransfers"`
	HasBaggage         bool       `json:"hasBaggage"`
	HasLuggage         bool       `json:"hasLuggage"`
	FareFamily         FareFamily `json:"fareFamily"`
	TariffType         string     `json:"tariffType"`
	AgeThreshold       int        `json:"ageThreshold,omitempty"`
	IsVirtualInterline bool       `json:"isVirtualInterline"`
}

type FareFamily struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	MarketingName string `json:"marketingName"`
	HasFareFamily bool   `json:"hasFareFamily"`
}

type TripBooking struct {
	ExpiresAt                     time.Time `json:"expiresAt"`
	TicketingTimeLimit            time.Time `json:"ticketingTimeLimit"`
	ProviderRecommendationLimit   time.Time `json:"providerRecommendationLimit"`
	ProviderRecommendationCreated time.Time `json:"providerRecommendationCreated"`
	CountOfBlanks                 int       `json:"countOfBlanks"`
	BookingWithPartialDataAllowed bool      `json:"bookingWithPartialDataAllowed"`
	BookingActualizationAllowed   bool      `json:"bookingActualizationAllowed"`
}

func (t *Trip) GetPrice() float64 { return t.Prices.Price }
func (t *Trip) HasBaggage() bool  { return t.Metadata.HasBaggage }
func (t *Trip) GetSRO() *sro.SRO  { return t.SRO }
func (t *Trip) GetForwardSegments() []TripSegment {
	var segs []TripSegment
	for _, s := range t.Segments {
		if s.Direction == 0 {
			segs = append(segs, s)
		}
	}
	return segs
}
