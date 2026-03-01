package service

import (
	flightsearchprovider "github.com/flearch/internal/provider/flightsearch"
	"github.com/flearch/internal/repository"
	"github.com/flearch/internal/repository/flightsearchrepo"
	fls "github.com/flearch/internal/service/flightsearch"
)

type ServiceBootstrap struct {
	FlightSearchService fls.FlightSearchService
}

func NewServiceBootstrap(cache repository.Cache) *ServiceBootstrap {
	return &ServiceBootstrap{
		FlightSearchService: fls.NewFlightSearchService(
			flightsearchprovider.NewFlightSearch(),
			flightsearchrepo.NewFlightSearchCache(cache),
		),
	}
}
