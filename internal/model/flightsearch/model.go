package flightsearch

import "time"

type FlightSearch struct {
	Flights  []Flight
	Metadata FlightSearchMetadata
}

type FlightSearchMetadata struct {
	TotalResults       int
	ProvidersQueried   int
	ProvidersSucceeded int
	ProvidersFailed    int
	SearchTimeMs       int
	CacheHit           bool
}

type Flight struct {
	ID                      string
	Provider                string
	AirlineName             string
	AirlineCode             string
	FlightNumber            string
	Origin                  Airport
	Destination             Airport
	DepartureTime           time.Time
	ArrivalTime             time.Time
	DurationMinutes         int
	Stops                   int
	Price                   Price
	AvailableSeats          int
	CabinClass              string
	Aircraft                string
	Amenities               []string
	Baggage                 Baggage
	BestPriceSameFlight     BestPriceSameFlight
	BestAmenitiesSameFlight BestAmenitiesSameFlight
}

type BestPriceSameFlight struct {
	Price      int64
	FlightCode string
}

type BestAmenitiesSameFlight struct {
	AmenitiesCount int
	FlightCode     string
}

type Airport struct {
	Code string
	City string
}

type Price struct {
	Amount   int64
	Currency string
}

type Baggage struct {
	CarryOn string
	Checked string
}
