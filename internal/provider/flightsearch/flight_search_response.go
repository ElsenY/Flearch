package flightsearchprovider

// AirAsia response structs

type AirAsiaResponse struct {
	Status  string          `json:"status"`
	Flights []AirAsiaFlight `json:"flights"`
}

type AirAsiaFlight struct {
	FlightCode   string        `json:"flight_code"`
	Airline      string        `json:"airline"`
	FromAirport  string        `json:"from_airport"`
	ToAirport    string        `json:"to_airport"`
	DepartTime   string        `json:"depart_time"`
	ArriveTime   string        `json:"arrive_time"`
	DurationHrs  float64       `json:"duration_hours"`
	DirectFlight bool          `json:"direct_flight"`
	Stops        []AirAsiaStop `json:"stops,omitempty"`
	PriceIDR     int64         `json:"price_idr"`
	Seats        int           `json:"seats"`
	CabinClass   string        `json:"cabin_class"`
	BaggageNote  string        `json:"baggage_note"`
}

type AirAsiaStop struct {
	Airport         string `json:"airport"`
	WaitTimeMinutes int    `json:"wait_time_minutes"`
}

// Batik Air response structs

type BatikAirResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Results []BatikAirFlight `json:"results"`
}

type BatikAirFlight struct {
	FlightNumber      string               `json:"flightNumber"`
	AirlineName       string               `json:"airlineName"`
	AirlineIATA       string               `json:"airlineIATA"`
	Origin            string               `json:"origin"`
	Destination       string               `json:"destination"`
	DepartureDateTime string               `json:"departureDateTime"`
	ArrivalDateTime   string               `json:"arrivalDateTime"`
	TravelTime        string               `json:"travelTime"`
	NumberOfStops     int                  `json:"numberOfStops"`
	Connections       []BatikAirConnection `json:"connections,omitempty"`
	Fare              BatikAirFare         `json:"fare"`
	SeatsAvailable    int                  `json:"seatsAvailable"`
	AircraftModel     string               `json:"aircraftModel"`
	BaggageInfo       string               `json:"baggageInfo"`
	OnboardServices   []string             `json:"onboardServices"`
	DurationMinutes   int                  `json:"duration_minutes"`
}

type BatikAirConnection struct {
	StopAirport  string `json:"stopAirport"`
	StopDuration string `json:"stopDuration"`
}

type BatikAirFare struct {
	BasePrice    int64  `json:"basePrice"`
	Taxes        int64  `json:"taxes"`
	TotalPrice   int64  `json:"totalPrice"`
	CurrencyCode string `json:"currencyCode"`
	Class        string `json:"class"`
}

// Garuda Indonesia response structs

type GarudaResponse struct {
	Status  string         `json:"status"`
	Flights []GarudaFlight `json:"flights"`
}

type GarudaFlight struct {
	FlightID        string          `json:"flight_id"`
	Airline         string          `json:"airline"`
	AirlineCode     string          `json:"airline_code"`
	Departure       GarudaEndpoint  `json:"departure"`
	Arrival         GarudaEndpoint  `json:"arrival"`
	DurationMinutes int             `json:"duration_minutes"`
	Stops           int             `json:"stops"`
	Aircraft        string          `json:"aircraft"`
	Price           GarudaPrice     `json:"price"`
	Segments        []GarudaSegment `json:"segments,omitempty"`
	AvailableSeats  int             `json:"available_seats"`
	FareClass       string          `json:"fare_class"`
	Baggage         GarudaBaggage   `json:"baggage"`
	Amenities       []string        `json:"amenities,omitempty"`
}

type GarudaEndpoint struct {
	Airport  string `json:"airport"`
	City     string `json:"city"`
	Time     string `json:"time"`
	Terminal string `json:"terminal"`
}

type GarudaPrice struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type GarudaSegment struct {
	FlightNumber    string         `json:"flight_number"`
	Departure       GarudaSegPoint `json:"departure"`
	Arrival         GarudaSegPoint `json:"arrival"`
	DurationMinutes int            `json:"duration_minutes"`
	LayoverMinutes  int            `json:"layover_minutes,omitempty"`
}

type GarudaSegPoint struct {
	Airport string `json:"airport"`
	Time    string `json:"time"`
}

type GarudaBaggage struct {
	CarryOn int `json:"carry_on"`
	Checked int `json:"checked"`
}

// Lion Air response structs

type LionAirResponse struct {
	Success bool        `json:"success"`
	Data    LionAirData `json:"data"`
}

type LionAirData struct {
	AvailableFlights []LionAirFlight `json:"available_flights"`
}

type LionAirFlight struct {
	ID         string           `json:"id"`
	Carrier    LionAirCarrier   `json:"carrier"`
	Route      LionAirRoute     `json:"route"`
	Schedule   LionAirSchedule  `json:"schedule"`
	FlightTime int              `json:"flight_time"`
	IsDirect   bool             `json:"is_direct"`
	StopCount  int              `json:"stop_count,omitempty"`
	Layovers   []LionAirLayover `json:"layovers,omitempty"`
	Pricing    LionAirPricing   `json:"pricing"`
	SeatsLeft  int              `json:"seats_left"`
	PlaneType  string           `json:"plane_type"`
	Services   LionAirServices  `json:"services"`
}

type LionAirCarrier struct {
	Name string `json:"name"`
	IATA string `json:"iata"`
}

type LionAirRoute struct {
	From LionAirAirport `json:"from"`
	To   LionAirAirport `json:"to"`
}

type LionAirAirport struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

type LionAirSchedule struct {
	Departure         string `json:"departure"`
	DepartureTimezone string `json:"departure_timezone"`
	Arrival           string `json:"arrival"`
	ArrivalTimezone   string `json:"arrival_timezone"`
}

type LionAirLayover struct {
	Airport         string `json:"airport"`
	DurationMinutes int    `json:"duration_minutes"`
}

type LionAirPricing struct {
	Total    int64  `json:"total"`
	Currency string `json:"currency"`
	FareType string `json:"fare_type"`
}

type LionAirServices struct {
	WifiAvailable    bool           `json:"wifi_available"`
	MealsIncluded    bool           `json:"meals_included"`
	BaggageAllowance LionAirBaggage `json:"baggage_allowance"`
}

type LionAirBaggage struct {
	Cabin string `json:"cabin"`
	Hold  string `json:"hold"`
}
