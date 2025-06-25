package local

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

type ristrettoCache struct {
	cache *ristretto.Cache[string, string]
}

// RistrettoConfig holds configuration for the in-memory Ristretto cache
type RistrettoConfig struct {
	NumCounters int64
	MaxCost     int64
	BufferItems int64
}

// DefaultRistrettoConfig returns sensible default values for Ristretto
func DefaultRistrettoConfig() RistrettoConfig {
	return RistrettoConfig{
		NumCounters: 1e7,     // 10M
		MaxCost:     1 << 30, // 1GB
		BufferItems: 64,
	}
}

func NewRistrettoCache(config ...RistrettoConfig) (*ristrettoCache, error) {
	cfg := DefaultRistrettoConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	c, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: cfg.NumCounters, // number of keys to track frequency
		MaxCost:     cfg.MaxCost,     // maximum cost of cache
		BufferItems: cfg.BufferItems, // number of keys per Get buffer
	})

	if err != nil {
		return nil, err
	}

	return &ristrettoCache{cache: c}, nil
}

func (r *ristrettoCache) Get(ctx context.Context, key string) (string, error) {
	val, found := r.cache.Get(key)
	if !found {
		return "", fmt.Errorf("key not found: %s", key)
	}

	return val, nil
}

func (r *ristrettoCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	strVal := fmt.Sprintf("%v", value)
	if strVal == "" {
		return fmt.Errorf("value cannot be empty")
	}

	// Use string length as a reasonable cost metric
	cost := int64(len(strVal))
	ok := r.cache.SetWithTTL(key, strVal, cost, expiration)
	if !ok {
		return fmt.Errorf("failed to set key: %s", key)
	}

	// Ensure value is visible immediately
	r.cache.Wait()

	return nil
}

func (r *ristrettoCache) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		r.cache.Del(key)
	}

	return nil
}

func (r *ristrettoCache) Ping(ctx context.Context) error {
	// Ristretto does not have a ping method, but we can check if the cache is initialized
	if r.cache == nil {
		return fmt.Errorf("cache is not initialized")
	}

	return nil
}

func (r *ristrettoCache) Close() error {
	r.cache.Close()
	return nil
}
