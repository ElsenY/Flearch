package sorthelper

import (
	"sort"
	"strings"

	flightsearchmodel "github.com/flearch/internal/model/flightsearch"
)

func SortFlights(flights []flightsearchmodel.Flight, sortBy string) {
	descending := false
	if strings.HasPrefix(sortBy, "-") {
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
				return flights[i].DepartureTime.After(flights[j].DepartureTime)
			} else {
				return flights[i].DepartureTime.Before(flights[j].DepartureTime)
			}
		})
	case "arrival_time":
		sort.Slice(flights, func(i, j int) bool {
			if descending {
				return flights[i].ArrivalTime.After(flights[j].ArrivalTime)
			} else {
				return flights[i].ArrivalTime.Before(flights[j].ArrivalTime)
			}
		})
	}
}
