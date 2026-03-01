package repository

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"sync"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (value string, err error)
	Set(ctx context.Context, key string, value interface{}) (err error)
}

type CacheValue struct {
	Value    interface{}
	Metadata cacheValueMetadata
}

type cacheValueMetadata struct {
	CreatedAt      time.Time
	LastAccessedAt time.Time
	HitCount       int64
}

type cache struct {
	cache       map[string]CacheValue
	capacity    int
	sampleCount int
	mu          sync.RWMutex
}

func NewCache(capacity int, sampleCount int) Cache {
	return &cache{
		cache:       make(map[string]CacheValue),
		capacity:    capacity,
		sampleCount: sampleCount,
	}
}

func (fc *cache) Get(ctx context.Context, key string) (string, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	value, ok := fc.cache[key]
	if !ok {
		return "", errors.New("key not found")
	}

	fc.cache[key] = CacheValue{
		Value: value.Value,
		Metadata: cacheValueMetadata{
			CreatedAt:      value.Metadata.CreatedAt,
			LastAccessedAt: time.Now(),
			HitCount:       value.Metadata.HitCount + 1,
		},
	}

	return value.Value.(string), nil
}

func (fc *cache) Set(ctx context.Context, key string, value interface{}) (err error) {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	fc.mu.Lock()
	defer fc.mu.Unlock()

	if len(fc.cache) >= fc.capacity {
		var minHitCount int64 = math.MaxInt64
		var oldestAccess time.Time
		minKey := ""

		count := 0
		for k, v := range fc.cache {
			if count >= fc.sampleCount {
				break
			}

			if v.Metadata.HitCount < minHitCount ||
				(v.Metadata.HitCount == minHitCount && (oldestAccess.IsZero() || v.Metadata.LastAccessedAt.Before(oldestAccess))) {
				minHitCount = v.Metadata.HitCount
				oldestAccess = v.Metadata.LastAccessedAt
				minKey = k
			}

			count++
		}

		delete(fc.cache, minKey)
	}

	fc.cache[key] = CacheValue{
		Value: string(b),
		Metadata: cacheValueMetadata{
			CreatedAt:      time.Now(),
			LastAccessedAt: time.Now(),
			HitCount:       0,
		},
	}

	return nil
}
