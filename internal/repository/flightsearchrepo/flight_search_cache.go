package flightsearchrepo

import (
	"context"
	"encoding/json"

	flightsearchmodel "github.com/flearch/internal/model/flightsearch"
	"github.com/flearch/internal/repository"
)

type FlightSearchCache interface {
	Get(ctx context.Context, key string) (value *flightsearchmodel.FlightSearch, err error)
	Set(ctx context.Context, key string, value *flightsearchmodel.FlightSearch) (err error)
}

type flightSearchCache struct {
	cache repository.Cache
}

func NewFlightSearchCache(cache repository.Cache) FlightSearchCache {
	return &flightSearchCache{cache: cache}
}

func (fc *flightSearchCache) Get(ctx context.Context, key string) (*flightsearchmodel.FlightSearch, error) {
	value, err := fc.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var flightSearch flightsearchmodel.FlightSearch
	err = json.Unmarshal([]byte(value), &flightSearch)
	if err != nil {
		return nil, err
	}

	return &flightSearch, nil
}

func (fc *flightSearchCache) Set(ctx context.Context, key string, value *flightsearchmodel.FlightSearch) error {
	err := fc.cache.Set(ctx, key, *value)
	if err != nil {
		return err
	}
	return nil
}
