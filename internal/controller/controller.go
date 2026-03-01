package controller

import (
	"time"

	"github.com/flearch/internal/config"
	"github.com/flearch/internal/controller/flightsearch"
	"github.com/flearch/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Routes()
}

type controller struct {
	app                    *fiber.App
	appConfig              *config.AppConfig
	flightSearchController flightsearch.FlightSearchController
}

func NewController(app *fiber.App, validator *validator.Validate, services *service.ServiceBootstrap, appConfig *config.AppConfig) Controller {
	return &controller{
		app:                    app,
		appConfig:              appConfig,
		flightSearchController: flightsearch.NewFlightSearchController(services.FlightSearchService, validator, time.Duration(appConfig.RequestTimeoutMs)*time.Millisecond),
	}
}

func (c *controller) Routes() {
	api := c.app.Group("/api")

	c.flightSearchController.Routes(api)
}
