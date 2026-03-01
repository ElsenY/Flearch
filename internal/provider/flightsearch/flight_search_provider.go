package flightsearchprovider

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/flearch/internal/constant"
	flightsearchmodel "github.com/flearch/internal/model/flightsearch"
)

//go:embed api_mock/*.json
var mockFS embed.FS

type FlightSearch interface {
	FindAllAirAsiaFlights(ctx context.Context, flightSearchRequest *AirAsiaFlightSearchRequest) ([]flightsearchmodel.Flight, error)
	FindAllBatikAirFlights(ctx context.Context, flightSearchRequest *BatikAirFlightSearchRequest) ([]flightsearchmodel.Flight, error)
	FindAllGarudaIndonesiaFlights(ctx context.Context, flightSearchRequest *GarudaIndonesiaFlightSearchRequest) ([]flightsearchmodel.Flight, error)
	FindAllLionAirFlights(ctx context.Context, flightSearchRequest *LionAirFlightSearchRequest) ([]flightsearchmodel.Flight, error)
}

func NewFlightSearch() FlightSearch {
	return &flightSearch{}
}

type flightSearch struct{}

func (fs *flightSearch) FindAllAirAsiaFlights(ctx context.Context, flightSearchRequest *AirAsiaFlightSearchRequest) ([]flightsearchmodel.Flight, error) {

	// mock get data from "api", usually we want to send filtering here, but since its mock json, we filter in the loop below
	data, err := mockFS.ReadFile("api_mock/airasia_search_response.json")
	if err != nil {
		return nil, fmt.Errorf("read airasia mock: %w", err)
	}

	var resp AirAsiaResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse airasia response: %w", err)
	}

	filtered := make([]AirAsiaFlight, 0, len(resp.Flights))
	for _, f := range resp.Flights {
		if flightSearchRequest.Filter(f) {
			filtered = append(filtered, f)
		}
	}

	if flightSearchRequest.SortBy != nil {
		SortAirAsiaFlights(filtered, *flightSearchRequest.SortBy)
	}

	if flightSearchRequest.Limit != nil && flightSearchRequest.Page != nil {
		filtered = filtered[(*flightSearchRequest.Limit)*(*flightSearchRequest.Page) : (*flightSearchRequest.Limit)*(*flightSearchRequest.Page)+*flightSearchRequest.Limit]
	}

	// mock delay (50-150ms)
	delay := time.Duration(50+rand.Intn(101)) * time.Millisecond
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(delay):
	}

	// mock failure
	if rand.Intn(101) < 10 {
		return nil, errors.New("Failed to get data from AirAsia")
	}

	flights := make([]flightsearchmodel.Flight, 0, len(filtered))
	for _, f := range filtered {
		depTime, err := time.Parse(time.RFC3339, f.DepartTime)
		if err != nil {
			return nil, fmt.Errorf("parse departure time %q: %w", f.DepartTime, err)
		}

		arrTime, err := time.Parse(time.RFC3339, f.ArriveTime)
		if err != nil {
			return nil, fmt.Errorf("parse arrival time %q: %w", f.ArriveTime, err)
		}

		baggateNote := strings.Split(f.BaggageNote, ", ")

		// based on mock data, total duration is already calculated layover hours, no need to add it again
		durationMin := int(math.Round(f.DurationHrs * 60))

		flights = append(flights, flightsearchmodel.Flight{
			ID:              f.FlightCode + "_AirAsia",
			Provider:        "AirAsia",
			AirlineName:     f.Airline,
			AirlineCode:     strings.Split(f.FlightCode, "")[0] + strings.Split(f.FlightCode, "")[1],
			FlightNumber:    f.FlightCode,
			Origin:          flightsearchmodel.Airport{Code: f.FromAirport, City: constant.AirportCityMapper[f.FromAirport]},
			Destination:     flightsearchmodel.Airport{Code: f.ToAirport, City: constant.AirportCityMapper[f.ToAirport]},
			DepartureTime:   depTime,
			ArrivalTime:     arrTime,
			DurationMinutes: durationMin,
			Stops:           len(f.Stops),
			Price:           flightsearchmodel.Price{Amount: f.PriceIDR, Currency: "IDR"},
			AvailableSeats:  f.Seats,
			CabinClass:      f.CabinClass,
			Baggage: flightsearchmodel.Baggage{
				CarryOn: baggateNote[0],
				Checked: baggateNote[1],
			},
		})
	}

	return flights, nil
}

func (fs *flightSearch) FindAllBatikAirFlights(ctx context.Context, flightSearchRequest *BatikAirFlightSearchRequest) ([]flightsearchmodel.Flight, error) {
	data, err := mockFS.ReadFile("api_mock/batik_air_search_response.json")
	if err != nil {
		return nil, fmt.Errorf("read batik air mock: %w", err)
	}

	var resp BatikAirResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse batik air response: %w", err)
	}

	const batikTimeLayout = "2006-01-02T15:04:05-0700"

	// pre-process DurationMinutes so filter and sort use the corrected value
	for i := range resp.Results {
		depTime, err := time.Parse(batikTimeLayout, resp.Results[i].DepartureDateTime)
		if err != nil {
			return nil, fmt.Errorf("parse departure time %q: %w", resp.Results[i].DepartureDateTime, err)
		}
		arrTime, err := time.Parse(batikTimeLayout, resp.Results[i].ArrivalDateTime)
		if err != nil {
			return nil, fmt.Errorf("parse arrival time %q: %w", resp.Results[i].ArrivalDateTime, err)
		}
		durationMin := int(arrTime.Sub(depTime).Minutes())
		resp.Results[i].DurationMinutes = durationMin + 5
	}

	filtered := make([]BatikAirFlight, 0, len(resp.Results))
	for _, f := range resp.Results {
		if flightSearchRequest.Filter(f) {
			filtered = append(filtered, f)
		}
	}

	if flightSearchRequest.SortBy != nil {
		SortBatikAirFlights(filtered, *flightSearchRequest.SortBy)
	}

	if flightSearchRequest.Limit != nil && flightSearchRequest.Page != nil {
		filtered = filtered[(*flightSearchRequest.Limit)*(*flightSearchRequest.Page) : (*flightSearchRequest.Limit)*(*flightSearchRequest.Page)+*flightSearchRequest.Limit]
	}

	// mock delay (200-400ms)
	delay := time.Duration(200+rand.Intn(200)) * time.Millisecond
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(delay):
	}

	flights := make([]flightsearchmodel.Flight, 0, len(filtered))
	for _, f := range filtered {
		depTime, err := time.Parse(batikTimeLayout, f.DepartureDateTime)
		if err != nil {
			return nil, fmt.Errorf("parse departure time %q: %w", f.DepartureDateTime, err)
		}
		arrTime, err := time.Parse(batikTimeLayout, f.ArrivalDateTime)
		if err != nil {
			return nil, fmt.Errorf("parse arrival time %q: %w", f.ArrivalDateTime, err)
		}
		durationMin := int(arrTime.Sub(depTime).Minutes())
		baggateNote := strings.Split(f.BaggageInfo, ", ")

		var services []string
		for _, s := range f.OnboardServices {
			services = append(services, strings.ToLower(s))
		}

		flights = append(flights, flightsearchmodel.Flight{
			ID:              f.FlightNumber + "_BatikAir",
			Provider:        "Batik Air",
			AirlineName:     f.AirlineName,
			AirlineCode:     f.AirlineIATA,
			FlightNumber:    f.FlightNumber,
			Origin:          flightsearchmodel.Airport{Code: f.Origin, City: constant.AirportCityMapper[f.Origin]},
			Destination:     flightsearchmodel.Airport{Code: f.Destination, City: constant.AirportCityMapper[f.Destination]},
			DepartureTime:   depTime,
			ArrivalTime:     arrTime,
			DurationMinutes: durationMin,
			Stops:           f.NumberOfStops,
			Price:           flightsearchmodel.Price{Amount: f.Fare.TotalPrice, Currency: f.Fare.CurrencyCode},
			AvailableSeats:  f.SeatsAvailable,
			CabinClass:      constant.CabinClassMapper[f.Fare.Class],
			Aircraft:        f.AircraftModel,
			Amenities:       services,
			Baggage: flightsearchmodel.Baggage{
				CarryOn: baggateNote[0],
				Checked: baggateNote[1],
			},
		})
	}

	return flights, nil
}

func (fs *flightSearch) FindAllGarudaIndonesiaFlights(ctx context.Context, flightSearchRequest *GarudaIndonesiaFlightSearchRequest) ([]flightsearchmodel.Flight, error) {
	data, err := mockFS.ReadFile("api_mock/garuda_indonesia_search_response.json")
	if err != nil {
		return nil, fmt.Errorf("read garuda mock: %w", err)
	}

	var resp GarudaResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse garuda response: %w", err)
	}

	// pre-process segments data so filter and sort use corrected values
	for i := range resp.Flights {
		if resp.Flights[i].Segments != nil {
			segmentsTimeSpent := 0
			for _, segment := range resp.Flights[i].Segments {
				segmentsTimeSpent += segment.DurationMinutes
				segmentsTimeSpent += segment.LayoverMinutes
			}
			resp.Flights[i].DurationMinutes = segmentsTimeSpent
			resp.Flights[i].Arrival.City = resp.Flights[i].Segments[len(resp.Flights[i].Segments)-1].Arrival.Airport
			resp.Flights[i].Arrival.Time = resp.Flights[i].Segments[len(resp.Flights[i].Segments)-1].Arrival.Time
			resp.Flights[i].Stops = len(resp.Flights[i].Segments) - 1
		}
	}

	filtered := make([]GarudaFlight, 0, len(resp.Flights))
	for _, f := range resp.Flights {
		if flightSearchRequest.Filter(f) {
			filtered = append(filtered, f)
		}
	}

	if flightSearchRequest.SortBy != nil {
		SortGarudaFlights(filtered, *flightSearchRequest.SortBy)
	}

	if flightSearchRequest.Limit != nil && flightSearchRequest.Page != nil {
		filtered = filtered[(*flightSearchRequest.Limit)*(*flightSearchRequest.Page) : (*flightSearchRequest.Limit)*(*flightSearchRequest.Page)+*flightSearchRequest.Limit]
	}

	// mock delay (50-100ms)
	delay := time.Duration(50+rand.Intn(51)) * time.Millisecond
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(delay):
	}

	flights := make([]flightsearchmodel.Flight, 0, len(filtered))
	for _, f := range filtered {
		depTime, err := time.Parse(time.RFC3339, f.Departure.Time)
		if err != nil {
			return nil, fmt.Errorf("parse departure time %q: %w", f.Departure.Time, err)
		}
		arrTime, err := time.Parse(time.RFC3339, f.Arrival.Time)
		if err != nil {
			return nil, fmt.Errorf("parse arrival time %q: %w", f.Arrival.Time, err)
		}

		stops := f.Stops
		if len(f.Segments) > 1 {
			stops = len(f.Segments) - 1
		}

		baggageCarryOn := fmt.Sprintf("%d piece", f.Baggage.CarryOn)
		baggageChecked := fmt.Sprintf("%d piece", f.Baggage.Checked)
		if f.Baggage.CarryOn != 1 {
			baggageCarryOn += "s"
		}
		if f.Baggage.Checked != 1 {
			baggageChecked += "s"
		}

		flights = append(flights, flightsearchmodel.Flight{
			ID:              f.FlightID + "_GarudaIndonesia",
			Provider:        "Garuda Indonesia",
			AirlineName:     f.Airline,
			AirlineCode:     f.AirlineCode,
			FlightNumber:    f.FlightID,
			Origin:          flightsearchmodel.Airport{Code: f.Departure.Airport, City: f.Departure.City},
			Destination:     flightsearchmodel.Airport{Code: f.Arrival.Airport, City: f.Arrival.City},
			DepartureTime:   depTime,
			ArrivalTime:     arrTime,
			DurationMinutes: f.DurationMinutes,
			Stops:           stops,
			Price:           flightsearchmodel.Price{Amount: f.Price.Amount, Currency: f.Price.Currency},
			AvailableSeats:  f.AvailableSeats,
			CabinClass:      f.FareClass,
			Aircraft:        f.Aircraft,
			Amenities:       f.Amenities,
			Baggage: flightsearchmodel.Baggage{
				CarryOn: baggageCarryOn,
				Checked: baggageChecked,
			},
		})
	}

	return flights, nil
}

func (fs *flightSearch) FindAllLionAirFlights(ctx context.Context, flightSearchRequest *LionAirFlightSearchRequest) ([]flightsearchmodel.Flight, error) {
	data, err := mockFS.ReadFile("api_mock/lion_air_search_response.json")
	if err != nil {
		return nil, fmt.Errorf("read lion air mock: %w", err)
	}

	var resp LionAirResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse lion air response: %w", err)
	}

	filtered := make([]LionAirFlight, 0, len(resp.Data.AvailableFlights))
	for _, f := range resp.Data.AvailableFlights {
		if flightSearchRequest.Filter(f) {
			filtered = append(filtered, f)
		}
	}

	if flightSearchRequest.SortBy != nil {
		SortLionAirFlights(filtered, *flightSearchRequest.SortBy)
	}

	if flightSearchRequest.Limit != nil && flightSearchRequest.Page != nil {
		filtered = filtered[(*flightSearchRequest.Limit)*(*flightSearchRequest.Page) : (*flightSearchRequest.Limit)*(*flightSearchRequest.Page)+*flightSearchRequest.Limit]
	}

	// mock delay (100-200ms)
	delay := time.Duration(100+rand.Intn(101)) * time.Millisecond
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(delay):
	}

	flights := make([]flightsearchmodel.Flight, 0, len(filtered))
	for _, f := range filtered {
		depLoc, err := time.LoadLocation(f.Schedule.DepartureTimezone)
		if err != nil {
			return nil, fmt.Errorf("load departure timezone %q: %w", f.Schedule.DepartureTimezone, err)
		}
		arrLoc, err := time.LoadLocation(f.Schedule.ArrivalTimezone)
		if err != nil {
			return nil, fmt.Errorf("load arrival timezone %q: %w", f.Schedule.ArrivalTimezone, err)
		}

		const localLayout = "2006-01-02T15:04:05"
		depTime, err := time.ParseInLocation(localLayout, f.Schedule.Departure, depLoc)
		if err != nil {
			return nil, fmt.Errorf("parse departure time %q: %w", f.Schedule.Departure, err)
		}
		arrTime, err := time.ParseInLocation(localLayout, f.Schedule.Arrival, arrLoc)
		if err != nil {
			return nil, fmt.Errorf("parse arrival time %q: %w", f.Schedule.Arrival, err)
		}

		stops := 0
		if !f.IsDirect {
			stops = f.StopCount
		}

		var amenities []string
		if f.Services.WifiAvailable {
			amenities = append(amenities, "wifi")
		}
		if f.Services.MealsIncluded {
			amenities = append(amenities, "meal")
		}

		flights = append(flights, flightsearchmodel.Flight{
			ID:              f.ID + "_LionAir",
			Provider:        "Lion Air",
			AirlineName:     f.Carrier.Name,
			AirlineCode:     f.Carrier.IATA,
			FlightNumber:    f.ID,
			Origin:          flightsearchmodel.Airport{Code: f.Route.From.Code, City: f.Route.From.City},
			Destination:     flightsearchmodel.Airport{Code: f.Route.To.Code, City: f.Route.To.City},
			DepartureTime:   depTime,
			ArrivalTime:     arrTime,
			DurationMinutes: f.FlightTime,
			Stops:           stops,
			Price:           flightsearchmodel.Price{Amount: f.Pricing.Total, Currency: f.Pricing.Currency},
			AvailableSeats:  f.SeatsLeft,
			CabinClass:      constant.CabinClassMapper[f.Pricing.FareType],
			Aircraft:        f.PlaneType,
			Amenities:       amenities,
			Baggage: flightsearchmodel.Baggage{
				CarryOn: f.Services.BaggageAllowance.Cabin,
				Checked: f.Services.BaggageAllowance.Hold,
			},
		})
	}

	return flights, nil
}
