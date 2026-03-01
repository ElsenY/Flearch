package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type ServerConfig struct {
	URL  string
	Port string
	Host string
}

func NewServerConfig() ServerConfig {
	return ServerConfig{
		URL:  os.Getenv("SERVER_URL"),
		Port: os.Getenv("SERVER_PORT"),
		Host: fmt.Sprintf("%s:%s", os.Getenv("SERVER_URL"), os.Getenv("SERVER_PORT")),
	}
}

type AppConfig struct {
	RequestTimeoutMs int
}

func NewAppConfig() AppConfig {
	requestTimeoutMs, err := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT_MS"))
	if err != nil {
		log.Println("error converting REQUEST_TIMEOUT_MS to int", err)
		requestTimeoutMs = 10000
	}

	return AppConfig{
		RequestTimeoutMs: requestTimeoutMs,
	}
}
