package flightsearchprovider

import (
	"strings"
)

type AirAsiaFlightSearchRequest struct {
	Origin          *string `json:"from_airport"`
	Destination     *string `json:"destination"`
	DepartureDate   *string `json:"depart_time"`
	ArrivalDate     *string `json:"arrive_time"`
	Airline         *string `json:"airline"`
	DurationMinutes *int    `json:"duration_minutes"`
	ReturnDate      *string `json:"return_date"`
	Passengers      *int    `json:"passengers"`
	CabinClass      *string `json:"cabin_class"`
	MinPrice        *int    `json:"min_price"`
	MaxPrice        *int    `json:"max_price"`
	Stops           *int    `json:"stops"`
	SortBy          *string `json:"sort_by"`
	Limit           *int    `json:"limit"`
	Page            *int    `json:"page"`
}

func (req *AirAsiaFlightSearchRequest) Filter(flight AirAsiaFlight) bool {
	if req.Origin != nil && *req.Origin != flight.FromAirport {
		return false
	}
	if req.Destination != nil && *req.Destination != flight.ToAirport {
		return false
	}

	if req.CabinClass != nil && !strings.EqualFold(*req.CabinClass, flight.CabinClass) {
		return false
	}
	if req.Passengers != nil && *req.Passengers > flight.Seats {
		return false
	}

	if req.MinPrice != nil && int(flight.PriceIDR) < *req.MinPrice {
		return false
	}
	if req.MaxPrice != nil && int(flight.PriceIDR) > *req.MaxPrice {
		return false
	}
	if req.Stops != nil && len(flight.Stops) != *req.Stops {
		return false
	}
	if req.DurationMinutes != nil && int(flight.DurationHrs*60) != *req.DurationMinutes {
		return false
	}
	if req.Airline != nil && flight.Airline != *req.Airline {
		return false
	}
	if req.ArrivalDate != nil && flight.ArriveTime != *req.ArrivalDate {
		return false
	}
	if req.DepartureDate != nil && flight.DepartTime != *req.DepartureDate {
		return false
	}
	if req.ReturnDate != nil && flight.ArriveTime != *req.ReturnDate {
		return false
	}
	return true
}

type BatikAirFlightSearchRequest struct {
	Origin          *string `json:"from_airport"`
	Destination     *string `json:"to_airport"`
	DepartureDate   *string `json:"departure_date"`
	ArrivalDate     *string `json:"arrive_time"`
	ReturnDate      *string `json:"return_date"`
	MinPrice        *int    `json:"min_price"`
	MaxPrice        *int    `json:"max_price"`
	Stops           *int    `json:"stops"`
	DurationMinutes *int    `json:"duration_minutes"`
	Airline         *string `json:"airline"`
	Passengers      *int    `json:"passengers"`
	CabinClass      *string `json:"cabin_class"`
	SortBy          *string `json:"sort_by"`
	Limit           *int    `json:"limit"`
	Page            *int    `json:"page"`
}

func (req *BatikAirFlightSearchRequest) Filter(flight BatikAirFlight) bool {
	if req.Origin != nil && *req.Origin != flight.Origin {
		return false
	}
	if req.Destination != nil && *req.Destination != flight.Destination {
		return false
	}
	if req.CabinClass != nil && !strings.EqualFold(*req.CabinClass, flight.Fare.Class) {
		return false
	}
	if req.Passengers != nil && *req.Passengers > flight.SeatsAvailable {
		return false
	}

	if req.MinPrice != nil && int(flight.Fare.BasePrice) < *req.MinPrice {
		return false
	}
	if req.MaxPrice != nil && int(flight.Fare.BasePrice) > *req.MaxPrice {
		return false
	}
	if req.Stops != nil && flight.NumberOfStops != *req.Stops {
		return false
	}
	if req.DurationMinutes != nil && flight.DurationMinutes != *req.DurationMinutes {
		return false
	}
	if req.Airline != nil && flight.AirlineName != *req.Airline {
		return false
	}
	if req.ArrivalDate != nil && flight.ArrivalDateTime != *req.ArrivalDate {
		return false
	}
	if req.DepartureDate != nil && flight.DepartureDateTime != *req.DepartureDate {
		return false
	}
	if req.ReturnDate != nil && flight.ArrivalDateTime != *req.ReturnDate {
		return false
	}
	return true
}

type GarudaIndonesiaFlightSearchRequest struct {
	Origin          *string `json:"origin"`
	Destination     *string `json:"destination"`
	DepartureDate   *string `json:"departure_date"`
	ReturnDate      *string `json:"return_date"`
	Passengers      *int    `json:"passengers"`
	CabinClass      *string `json:"cabin_class"`
	MinPrice        *int    `json:"min_price"`
	MaxPrice        *int    `json:"max_price"`
	DurationMinutes *int    `json:"duration_minutes"`
	Airline         *string `json:"airline"`
	ArrivalDate     *string `json:"arrive_time"`
	Stops           *int    `json:"stops"`
	SortBy          *string `json:"sort_by"`
	Limit           *int    `json:"limit"`
	Page            *int    `json:"page"`
}

func (req *GarudaIndonesiaFlightSearchRequest) Filter(flight GarudaFlight) bool {
	if req.Origin != nil && *req.Origin != flight.Departure.Airport {
		return false
	}
	if req.Destination != nil && *req.Destination != flight.Arrival.Airport {
		return false
	}
	if req.CabinClass != nil && !strings.EqualFold(*req.CabinClass, flight.FareClass) {
		return false
	}
	if req.Passengers != nil && *req.Passengers > flight.AvailableSeats {
		return false
	}
	if req.MinPrice != nil && int(flight.Price.Amount) < *req.MinPrice {
		return false
	}
	if req.MaxPrice != nil && int(flight.Price.Amount) > *req.MaxPrice {
		return false
	}
	if req.DurationMinutes != nil && flight.DurationMinutes != *req.DurationMinutes {
		return false
	}
	if req.Airline != nil && flight.Airline != *req.Airline {
		return false
	}
	if req.ArrivalDate != nil && flight.Arrival.Time != *req.ArrivalDate {
		return false
	}
	if req.DepartureDate != nil && flight.Departure.Time != *req.DepartureDate {
		return false
	}
	if req.ReturnDate != nil && flight.Arrival.Time != *req.ReturnDate {
		return false
	}
	return true
}

type LionAirFlightSearchRequest struct {
	Origin          *string `json:"origin"`
	Destination     *string `json:"destination"`
	DepartureDate   *string `json:"departure_date"`
	ReturnDate      *string `json:"return_date"`
	Passengers      *int    `json:"passengers"`
	CabinClass      *string `json:"cabin_class"`
	MinPrice        *int    `json:"min_price"`
	MaxPrice        *int    `json:"max_price"`
	DurationMinutes *int    `json:"duration_minutes"`
	Airline         *string `json:"airline"`
	ArrivalDate     *string `json:"arrive_time"`
	Stops           *int    `json:"stops"`
	SortBy          *string `json:"sort_by"`
	Limit           *int    `json:"limit"`
	Page            *int    `json:"page"`
}

func (req *LionAirFlightSearchRequest) Filter(flight LionAirFlight) bool {
	if req.Origin != nil && *req.Origin != flight.Route.From.Code {
		return false
	}
	if req.Destination != nil && *req.Destination != flight.Route.To.Code {
		return false
	}
	if req.CabinClass != nil && !strings.EqualFold(*req.CabinClass, flight.Pricing.FareType) {
		return false
	}
	if req.Passengers != nil && *req.Passengers > flight.SeatsLeft {
		return false
	}
	if req.MinPrice != nil && int(flight.Pricing.Total) < *req.MinPrice {
		return false
	}
	if req.MaxPrice != nil && int(flight.Pricing.Total) > *req.MaxPrice {
		return false
	}
	if req.DurationMinutes != nil && flight.FlightTime != *req.DurationMinutes {
		return false
	}
	if req.Airline != nil && flight.Carrier.Name != *req.Airline {
		return false
	}
	if req.ArrivalDate != nil && flight.Schedule.Arrival != *req.ArrivalDate {
		return false
	}
	if req.DepartureDate != nil && flight.Schedule.Departure != *req.DepartureDate {
		return false
	}
	if req.ReturnDate != nil && flight.Schedule.Arrival != *req.ReturnDate {
		return false
	}
	return true
}
