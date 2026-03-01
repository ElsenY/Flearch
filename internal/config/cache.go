package config

import (
	"log"
	"os"
	"strconv"
)

type CacheConfig struct {
	Capacity    int
	SampleCount int
}

func NewCacheConfig() CacheConfig {
	capacity, err := strconv.Atoi(os.Getenv("CACHE_CAPACITY"))
	if err != nil {
		log.Println("error converting CACHE_CAPACITY to int", err)
		return CacheConfig{
			Capacity:    10,
			SampleCount: 10,
		}
	}

	sampleCount, err := strconv.Atoi(os.Getenv("CACHE_SAMPLE_COUNT"))
	if err != nil {
		log.Println("error converting CACHE_SAMPLE_COUNT to int", err)
		return CacheConfig{
			Capacity:    10,
			SampleCount: 10,
		}
	}

	return CacheConfig{
		Capacity:    capacity,
		SampleCount: sampleCount,
	}
}
