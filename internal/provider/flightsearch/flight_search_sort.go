package flightsearchprovider

import (
	"sort"
	"strings"
)

func SortAirAsiaFlights(flights []AirAsiaFlight, sortBy string) {
	descending := false
	if strings.HasPrefix(sortBy, "-") {
		sortBy = strings.TrimPrefix(sortBy, "-")
		descending = true
	}

	switch sortBy {
	case "price":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].PriceIDR > flights[j].PriceIDR
			} else {
				return flights[i].PriceIDR < flights[j].PriceIDR
			}
		})
	case "duration":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].DurationHrs > flights[j].DurationHrs
			} else {
				return flights[i].DurationHrs < flights[j].DurationHrs
			}
		})
	case "departure_time":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].DepartTime > flights[j].DepartTime
			} else {
				return flights[i].DepartTime < flights[j].DepartTime
			}
		})
	case "arrival_time":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].ArriveTime > flights[j].ArriveTime
			} else {
				return flights[i].ArriveTime < flights[j].ArriveTime
			}
		})
	}
}

func SortBatikAirFlights(flights []BatikAirFlight, sortBy string) {
	descending := false
	if strings.HasPrefix(sortBy, "-") {
		sortBy = strings.TrimPrefix(sortBy, "-")
		descending = true
	}

	switch sortBy {
	case "price":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].Fare.TotalPrice > flights[j].Fare.TotalPrice
			} else {
				return flights[i].Fare.TotalPrice < flights[j].Fare.TotalPrice
			}
		})
	case "duration":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].DurationMinutes > flights[j].DurationMinutes
			} else {
				return flights[i].DurationMinutes < flights[j].DurationMinutes
			}
		})
	case "departure_time":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].DepartureDateTime > flights[j].DepartureDateTime
			} else {
				return flights[i].DepartureDateTime < flights[j].DepartureDateTime
			}
		})
	case "arrival_time":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].ArrivalDateTime > flights[j].ArrivalDateTime
			} else {
				return flights[i].ArrivalDateTime < flights[j].ArrivalDateTime
			}
		})
	}
}

func SortGarudaFlights(flights []GarudaFlight, sortBy string) {
	descending := false
	if strings.HasPrefix(sortBy, "-") {
		sortBy = strings.TrimPrefix(sortBy, "-")
		descending = true
	}

	switch sortBy {
	case "price":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].Price.Amount > flights[j].Price.Amount
			} else {
				return flights[i].Price.Amount < flights[j].Price.Amount
			}
		})
	case "duration":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].DurationMinutes > flights[j].DurationMinutes
			} else {
				return flights[i].DurationMinutes < flights[j].DurationMinutes
			}
		})
	case "departure_time":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].Departure.Time > flights[j].Departure.Time
			} else {
				return flights[i].Departure.Time < flights[j].Departure.Time
			}
		})
	case "arrival_time":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].Arrival.Time > flights[j].Arrival.Time
			} else {
				return flights[i].Arrival.Time < flights[j].Arrival.Time
			}
		})
	}
}

func SortLionAirFlights(flights []LionAirFlight, sortBy string) {
	descending := false
	if strings.HasPrefix(sortBy, "-") {
		sortBy = strings.TrimPrefix(sortBy, "-")
		descending = true
	}

	switch sortBy {
	case "price":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].Pricing.Total > flights[j].Pricing.Total
			} else {
				return flights[i].Pricing.Total < flights[j].Pricing.Total
			}
		})
	case "duration":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].FlightTime > flights[j].FlightTime
			} else {
				return flights[i].FlightTime < flights[j].FlightTime
			}
		})
	case "departure_time":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].Schedule.Departure > flights[j].Schedule.Departure
			} else {
				return flights[i].Schedule.Departure < flights[j].Schedule.Departure
			}
		})
	case "arrival_time":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].Schedule.Arrival > flights[j].Schedule.Arrival
			} else {
				return flights[i].Schedule.Arrival < flights[j].Schedule.Arrival
			}
		})
	}
}
