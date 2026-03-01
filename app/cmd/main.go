package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/flearch/internal/config"
	"github.com/flearch/internal/controller"
	"github.com/flearch/internal/repository"
	"github.com/flearch/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func runHttpServer(ctx context.Context,
	errChan chan error,
	serverConfig *config.ServerConfig,
	validator *validator.Validate,
	services *service.ServiceBootstrap,
	appConfig *config.AppConfig,
) {
	log.Println("Starting HTTP server...")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit:                1024 * 1024 * 60, // 60MB
		DisableStartupMessage:    true,
		EnableSplittingOnParsers: true,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
	})

	// provide health check
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	controllers := controller.NewController(app, validator, services, appConfig)
	controllers.Routes()

	go func() {
		log.Printf("App running on address %s", serverConfig.Host)
		if err := app.Listen(serverConfig.Host); err != nil {
			log.Printf("Failed running app: %v", err)
			errChan <- err
		}
	}()

	// Wait for context cancellation to trigger shutdown
	<-ctx.Done()
	log.Println("Shutting down HTTP server...")

	if err := app.Shutdown(); err != nil {
		log.Printf("HTTP server was not successfully shutdown. Forcing exit: %v", err)
		return
	}

	log.Println("HTTP server was successfully shutdown")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, set up default values")
		os.Setenv("SERVER_URL", "localhost")
		os.Setenv("SERVER_PORT", "8080")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error)
	defer close(errChan)

	config := config.NewConfig()
	validator := validator.New()
	cache := repository.NewCache(config.CacheConfig.Capacity, config.CacheConfig.SampleCount)
	services := service.NewServiceBootstrap(cache)
	go runHttpServer(ctx, errChan, &config.ServerConfig, validator, services, &config.AppConfig)

	select {
	case err := <-errChan:
		log.Println(err.Error())
		os.Exit(1)
	case <-ctx.Done():
		log.Println("Cancelled by context")
		os.Exit(0)
	}
}
