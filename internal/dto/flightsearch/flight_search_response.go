package flightsearchdto

type FlightSearchResponse struct {
	SearchCriteria SearchCriteria `json:"search_criteria"`
	Metadata       Metadata       `json:"metadata"`
	Flights        []FlightDTO    `json:"flights"`
}

type SearchCriteria struct {
	Origin          *string `json:"origin"`
	Destination     *string `json:"destination"`
	DepartureDate   *string `json:"departure_date"`
	ReturnDate      *string `json:"return_date"`
	Passengers      *int    `json:"passengers"`
	CabinClass      *string `json:"cabin_class"`
	MinPrice        *int    `json:"min_price"`
	MaxPrice        *int    `json:"max_price"`
	Stops           *int    `json:"stops"`
	DurationMinutes *int    `json:"duration_minutes"`
	Airline         *string `json:"airline"`
	ArrivalDate     *string `json:"arrival_date"`
	Limit           *int    `json:"limit"`
	Page            *int    `json:"page"`
	SortBy          *string `json:"sort_by"`
}

type Metadata struct {
	TotalResults       int  `json:"total_results"`
	ProvidersQueried   int  `json:"providers_queried"`
	ProvidersSucceeded int  `json:"providers_succeeded"`
	ProvidersFailed    int  `json:"providers_failed"`
	SearchTimeMs       int  `json:"search_time_ms"`
	CacheHit           bool `json:"cache_hit"`
}

type FlightDTO struct {
	ID                      string                     `json:"id"`
	Provider                string                     `json:"provider"`
	Airline                 AirlineDTO                 `json:"airline"`
	FlightNumber            string                     `json:"flight_number"`
	Departure               EndpointDTO                `json:"departure"`
	Arrival                 EndpointDTO                `json:"arrival"`
	Duration                DurationDTO                `json:"duration"`
	Stops                   int                        `json:"stops"`
	Price                   PriceDTO                   `json:"price"`
	AvailableSeats          int                        `json:"available_seats"`
	CabinClass              string                     `json:"cabin_class"`
	Aircraft                *string                    `json:"aircraft"`
	Amenities               []string                   `json:"amenities"`
	Baggage                 BaggageDTO                 `json:"baggage"`
	BestPriceSameFlight     BestPriceSameFlightDTO     `json:"best_price_same_flight"`
	BestAmenitiesSameFlight BestAmenitiesSameFlightDTO `json:"best_amenities_same_flight"`
}

type BestPriceSameFlightDTO struct {
	Price      int64  `json:"price"`
	FlightCode string `json:"flight_code"`
}

type BestAmenitiesSameFlightDTO struct {
	AmenitiesCount int    `json:"amenities_count"`
	FlightCode     string `json:"flight_code"`
}
type AirlineDTO struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type EndpointDTO struct {
	Airport   string `json:"airport"`
	City      string `json:"city"`
	Datetime  string `json:"datetime"`
	Timestamp int64  `json:"timestamp"`
}

type DurationDTO struct {
	TotalMinutes int    `json:"total_minutes"`
	Formatted    string `json:"formatted"`
}

type PriceDTO struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type BaggageDTO struct {
	CarryOn string `json:"carry_on"`
	Checked string `json:"checked"`
}
