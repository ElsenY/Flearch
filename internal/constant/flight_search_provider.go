package constant

var CabinClassMapper = map[string]string{
	"Y":       "economy",
	"ECONOMY": "economy",
}

var AirportCityMapper = map[string]string{
	"CGK": "Jakarta",
	"DPS": "Denpasar",
	"SUB": "Surabaya",
	"UPG": "Makassar",
}

var FlightSearchSortKeys = map[string]string{
	"price":          "price",
	"duration":       "duration",
	"departure_time": "departure_time",
	"arrival_time":   "arrival_time",
}
