package flightsearchdto

type FlightSearchRequest struct {
	Origin          *string `json:"origin"`
	Destination     *string `json:"destination"`
	DepartureDate   *string `json:"departureDate"`
	ReturnDate      *string `json:"returnDate"`
	Passengers      *int    `json:"passengers"`
	CabinClass      *string `json:"cabinClass"`
	MinPrice        *int    `json:"min_price"`
	MaxPrice        *int    `json:"max_price"`
	Stops           *int    `json:"stops"`
	DurationMinutes *int    `json:"duration_minutes"`
	Airline         *string `json:"airline"`
	ArrivalDate     *string `json:"arrive_time"`
	SortBy          *string `json:"sort_by"`
	Limit           *int    `json:"limit"`
	Page            *int    `json:"page"`
}
