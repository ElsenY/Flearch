package flightsearchsvc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	flightsearchdto "github.com/flearch/internal/dto/flightsearch"
	helper "github.com/flearch/internal/helper/retry"
	sorthelper "github.com/flearch/internal/helper/sort"
	flightsearchmodel "github.com/flearch/internal/model/flightsearch"
	flightsearchprovider "github.com/flearch/internal/provider/flightsearch"
	"github.com/flearch/internal/repository/flightsearchrepo"
)

type FlightSearchService interface {
	SearchFlights(ctx context.Context, flightSearchRequest *flightsearchdto.FlightSearchRequest) (*flightsearchmodel.FlightSearch, error)
}

type flightSearchService struct {
	flightSearchProvider flightsearchprovider.FlightSearch
	flightSearchCache    flightsearchrepo.FlightSearchCache
}

func NewFlightSearchService(flightSearchProvider flightsearchprovider.FlightSearch, flightSearchCache flightsearchrepo.FlightSearchCache) FlightSearchService {
	return &flightSearchService{flightSearchProvider: flightSearchProvider, flightSearchCache: flightSearchCache}
}

type providerFetcher func(ctx context.Context, req *flightsearchdto.FlightSearchRequest, flightSearchProvider flightsearchprovider.FlightSearch) ([]flightsearchmodel.Flight, error)

func (fs *flightSearchService) SearchFlights(ctx context.Context, flightSearchRequest *flightsearchdto.FlightSearchRequest) (*flightsearchmodel.FlightSearch, error) {
	b, err := json.Marshal(flightSearchRequest)
	if err != nil {
		return nil, err
	}

	cacheKey := sha256.Sum256(b)
	cacheKeyString := hex.EncodeToString(cacheKey[:])

	// check if there's cached data
	cacheValue, err := fs.flightSearchCache.Get(ctx, cacheKeyString)
	if err != nil {
		if !errors.Is(err, errors.New("key not found")) {
			log.Println("error getting cache value", err)
		}
	}

	if cacheValue != nil {
		cacheValue.Metadata.CacheHit = true
		return cacheValue, nil
	}

	// fetch flights from availableproviders
	fetchers := []providerFetcher{
		fetchAirAsia,
		fetchBatikAir,
		fetchGarudaIndonesia,
		fetchLionAir,
	}

	var (
		mu          sync.Mutex
		wg          sync.WaitGroup
		flights     []flightsearchmodel.Flight
		errMessages []string
	)

	wg.Add(len(fetchers))
	for _, fetch := range fetchers {
		go func(fn providerFetcher) {
			defer wg.Done()

			result, err := fn(ctx, flightSearchRequest, fs.flightSearchProvider)

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errMessages = append(errMessages, err.Error())
				return
			}
			flights = append(flights, result...)
		}(fetch)
	}

	wg.Wait()

	// give warning if there's any provider failed
	for _, v := range errMessages {
		log.Println(v)
	}

	// if all providers failed, return error
	if len(errMessages) == len(fetchers) {
		return &flightsearchmodel.FlightSearch{}, errors.New("failed to fetch flights")
	}

	// compose metadata
	metadata := flightsearchmodel.FlightSearchMetadata{
		TotalResults:       len(flights),
		ProvidersQueried:   len(fetchers),
		ProvidersSucceeded: len(fetchers) - len(errMessages),
		ProvidersFailed:    len(errMessages),
		CacheHit:           false,
	}

	// sort & limit flights if requested
	if flightSearchRequest.SortBy != nil {
		sorthelper.SortFlights(flights, *flightSearchRequest.SortBy)
	}

	if flightSearchRequest.Limit != nil && flightSearchRequest.Page != nil {
		flights = flights[(*flightSearchRequest.Limit)*(*flightSearchRequest.Page) : (*flightSearchRequest.Limit)*(*flightSearchRequest.Page)+*flightSearchRequest.Limit]
	}

	// ranking best price & amenities for same flight route
	minPriceMap := map[string]flightsearchmodel.BestPriceSameFlight{}
	maxAmenitiesMap := map[string]flightsearchmodel.BestAmenitiesSameFlight{}
	for _, v := range flights {
		flightkey := v.Origin.Code + v.Destination.Code
		if _, ok := minPriceMap[flightkey]; !ok {
			minPriceMap[flightkey] = flightsearchmodel.BestPriceSameFlight{
				Price:      int64(v.Price.Amount),
				FlightCode: v.ID,
			}
		} else {
			if int64(v.Price.Amount) < minPriceMap[flightkey].Price {
				minPriceMap[flightkey] = flightsearchmodel.BestPriceSameFlight{
					Price:      int64(v.Price.Amount),
					FlightCode: v.ID,
				}
			}
		}

		if _, ok := maxAmenitiesMap[flightkey]; !ok {
			maxAmenitiesMap[flightkey] = flightsearchmodel.BestAmenitiesSameFlight{
				AmenitiesCount: len(v.Amenities),
				FlightCode:     v.ID,
			}
		} else {
			if len(v.Amenities) > maxAmenitiesMap[flightkey].AmenitiesCount {
				maxAmenitiesMap[flightkey] = flightsearchmodel.BestAmenitiesSameFlight{
					AmenitiesCount: len(v.Amenities),
					FlightCode:     v.ID,
				}
			}
		}
	}

	for i, v := range flights {
		flightkey := v.Origin.Code + v.Destination.Code
		if _, ok := minPriceMap[flightkey]; ok {
			flights[i].BestPriceSameFlight = minPriceMap[flightkey]
		}
		if _, ok := maxAmenitiesMap[flightkey]; ok {
			flights[i].BestAmenitiesSameFlight = maxAmenitiesMap[flightkey]
		}
	}

	res := &flightsearchmodel.FlightSearch{
		Flights:  flights,
		Metadata: metadata,
	}

	// only do cache if all providers succeeded
	if metadata.ProvidersFailed <= 0 {
		go fs.flightSearchCache.Set(ctx, cacheKeyString, res)
	}

	return res, nil
}

// we declare internal func as var so we can mock it easier in unit tests
var fetchAirAsia = func(ctx context.Context, req *flightsearchdto.FlightSearchRequest, flightSearchProvider flightsearchprovider.FlightSearch) ([]flightsearchmodel.Flight, error) {
	var flights []flightsearchmodel.Flight
	err := helper.RetryRequest(ctx, func() error {
		result, err := flightSearchProvider.FindAllAirAsiaFlights(ctx, &flightsearchprovider.AirAsiaFlightSearchRequest{
			Origin:          req.Origin,
			Destination:     req.Destination,
			DepartureDate:   req.DepartureDate,
			ReturnDate:      req.ReturnDate,
			Passengers:      req.Passengers,
			CabinClass:      req.CabinClass,
			MinPrice:        req.MinPrice,
			MaxPrice:        req.MaxPrice,
			Stops:           req.Stops,
			DurationMinutes: req.DurationMinutes,
			Airline:         req.Airline,
			ArrivalDate:     req.ArrivalDate,
			Limit:           req.Limit,
			Page:            req.Page,
			SortBy:          req.SortBy,
		})

		if err != nil {
			return err
		}

		flights = result
		return nil
	}, 3, 1*time.Second)

	return flights, err
}

// we declare internal func as var so we can mock it easier in unit tests
var fetchBatikAir = func(ctx context.Context, req *flightsearchdto.FlightSearchRequest, flightSearchProvider flightsearchprovider.FlightSearch) ([]flightsearchmodel.Flight, error) {
	var flights []flightsearchmodel.Flight
	err := helper.RetryRequest(ctx, func() error {
		result, err := flightSearchProvider.FindAllBatikAirFlights(ctx, &flightsearchprovider.BatikAirFlightSearchRequest{
			Origin:          req.Origin,
			Destination:     req.Destination,
			DepartureDate:   req.DepartureDate,
			ReturnDate:      req.ReturnDate,
			Passengers:      req.Passengers,
			CabinClass:      req.CabinClass,
			MinPrice:        req.MinPrice,
			MaxPrice:        req.MaxPrice,
			Stops:           req.Stops,
			DurationMinutes: req.DurationMinutes,
			Airline:         req.Airline,
			ArrivalDate:     req.ArrivalDate,
			Limit:           req.Limit,
			Page:            req.Page,
			SortBy:          req.SortBy,
		})

		if err != nil {
			return err
		}

		flights = result
		return nil
	}, 3, 1*time.Second)

	return flights, err
}

// we declare internal func as var so we can mock it easier in unit tests
var fetchGarudaIndonesia = func(ctx context.Context, req *flightsearchdto.FlightSearchRequest, flightSearchProvider flightsearchprovider.FlightSearch) ([]flightsearchmodel.Flight, error) {
	var flights []flightsearchmodel.Flight
	err := helper.RetryRequest(ctx, func() error {
		result, err := flightSearchProvider.FindAllGarudaIndonesiaFlights(ctx, &flightsearchprovider.GarudaIndonesiaFlightSearchRequest{
			Origin:          req.Origin,
			Destination:     req.Destination,
			DepartureDate:   req.DepartureDate,
			ReturnDate:      req.ReturnDate,
			Passengers:      req.Passengers,
			CabinClass:      req.CabinClass,
			MinPrice:        req.MinPrice,
			MaxPrice:        req.MaxPrice,
			Stops:           req.Stops,
			DurationMinutes: req.DurationMinutes,
			Airline:         req.Airline,
			ArrivalDate:     req.ArrivalDate,
			Limit:           req.Limit,
			Page:            req.Page,
			SortBy:          req.SortBy,
		})

		if err != nil {
			return err
		}
		flights = result
		return nil
	}, 3, 1*time.Second)

	return flights, err
}

// we declare internal func as var so we can mock it easier in unit tests
var fetchLionAir = func(ctx context.Context, req *flightsearchdto.FlightSearchRequest, flightSearchProvider flightsearchprovider.FlightSearch) ([]flightsearchmodel.Flight, error) {
	var flights []flightsearchmodel.Flight
	err := helper.RetryRequest(ctx, func() error {
		result, err := flightSearchProvider.FindAllLionAirFlights(ctx, &flightsearchprovider.LionAirFlightSearchRequest{
			Origin:          req.Origin,
			Destination:     req.Destination,
			DepartureDate:   req.DepartureDate,
			ReturnDate:      req.ReturnDate,
			Passengers:      req.Passengers,
			CabinClass:      req.CabinClass,
			MinPrice:        req.MinPrice,
			MaxPrice:        req.MaxPrice,
			Stops:           req.Stops,
			DurationMinutes: req.DurationMinutes,
			Airline:         req.Airline,
			ArrivalDate:     req.ArrivalDate,
			Limit:           req.Limit,
			Page:            req.Page,
			SortBy:          req.SortBy,
		})

		if err != nil {
			return err
		}

		flights = result
		return nil
	}, 3, 1*time.Second)

	return flights, err
}
