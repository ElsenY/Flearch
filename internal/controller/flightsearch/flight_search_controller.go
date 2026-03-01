package flightsearch

import (
	"context"
	"fmt"
	"time"

	flightsearchdto "github.com/flearch/internal/dto/flightsearch"
	helper "github.com/flearch/internal/helper/dto"
	flightsearchmodel "github.com/flearch/internal/model/flightsearch"
	fls "github.com/flearch/internal/service/flightsearch"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type FlightSearchController interface {
	Routes(r fiber.Router)
}

type flightSearchController struct {
	timeout             time.Duration
	flightSearchService fls.FlightSearchService
	validator           *validator.Validate
}

func NewFlightSearchController(flightSearchService fls.FlightSearchService, validator *validator.Validate, timeout time.Duration) FlightSearchController {
	return &flightSearchController{timeout: timeout, flightSearchService: flightSearchService, validator: validator}
}

func (c *flightSearchController) Routes(r fiber.Router) {
	// /api/flight-search
	r = r.Group("/flight-search")

	r.Get("/", c.SearchFlights)
}

func (flc *flightSearchController) SearchFlights(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), flc.timeout)
	defer cancel()

	flightSearchRequest := &flightsearchdto.FlightSearchRequest{}

	if err := c.BodyParser(flightSearchRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := flc.validator.Struct(flightSearchRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	flightSearch, err := flc.flightSearchService.SearchFlights(ctx, flightSearchRequest)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return c.Status(fiber.StatusGatewayTimeout).JSON(fiber.Map{
				"error": "request timed out",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	deadline, _ := ctx.Deadline()
	resp := flightsearchdto.FlightSearchResponse{
		SearchCriteria: flightsearchdto.SearchCriteria{
			Origin:        flightSearchRequest.Origin,
			Destination:   flightSearchRequest.Destination,
			DepartureDate: flightSearchRequest.DepartureDate,
			Passengers:    flightSearchRequest.Passengers,
			CabinClass:    flightSearchRequest.CabinClass,
		},
		Metadata: flightsearchdto.Metadata{
			TotalResults:       flightSearch.Metadata.TotalResults,
			ProvidersQueried:   flightSearch.Metadata.ProvidersQueried,
			ProvidersSucceeded: flightSearch.Metadata.ProvidersSucceeded,
			ProvidersFailed:    flightSearch.Metadata.ProvidersFailed,
			SearchTimeMs:       int(flc.timeout.Milliseconds() - time.Until(deadline).Milliseconds()),
			CacheHit:           flightSearch.Metadata.CacheHit,
		},
		Flights: helper.MapToDTO(flightSearch.Flights, func(f flightsearchmodel.Flight) flightsearchdto.FlightDTO {
			return flightsearchdto.FlightDTO{
				ID:       f.ID,
				Provider: f.Provider,
				Airline: flightsearchdto.AirlineDTO{
					Name: f.AirlineName,
					Code: f.AirlineCode,
				},
				FlightNumber: f.FlightNumber,
				Departure: flightsearchdto.EndpointDTO{
					Airport:   f.Origin.Code,
					City:      f.Origin.City,
					Datetime:  f.DepartureTime.Format(time.RFC3339),
					Timestamp: f.DepartureTime.Unix(),
				},
				Arrival: flightsearchdto.EndpointDTO{
					Airport:   f.Destination.Code,
					City:      f.Destination.City,
					Datetime:  f.ArrivalTime.Format(time.RFC3339),
					Timestamp: f.ArrivalTime.Unix(),
				},
				Duration: flightsearchdto.DurationDTO{
					TotalMinutes: f.DurationMinutes,
					Formatted:    fmt.Sprintf("%dh %02dm", f.DurationMinutes/60, f.DurationMinutes%60),
				},
				Stops: f.Stops,
				Price: flightsearchdto.PriceDTO{
					Amount:   f.Price.Amount,
					Currency: f.Price.Currency,
				},
				AvailableSeats: f.AvailableSeats,
				CabinClass:     f.CabinClass,
				Aircraft:       &f.Aircraft,
				Amenities:      f.Amenities,
				Baggage: flightsearchdto.BaggageDTO{
					CarryOn: f.Baggage.CarryOn,
					Checked: f.Baggage.Checked,
				},
				BestPriceSameFlight: flightsearchdto.BestPriceSameFlightDTO{
					Price:      f.BestPriceSameFlight.Price,
					FlightCode: f.BestPriceSameFlight.FlightCode,
				},
				BestAmenitiesSameFlight: flightsearchdto.BestAmenitiesSameFlightDTO{
					AmenitiesCount: f.BestAmenitiesSameFlight.AmenitiesCount,
					FlightCode:     f.BestAmenitiesSameFlight.FlightCode,
				},
			}
		}),
	}

	return c.JSON(resp)
}
