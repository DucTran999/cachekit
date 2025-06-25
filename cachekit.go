package cachekit

import (
	"context"
	"time"

	"github.com/DucTran999/cachekit/config"
	cacheerr "github.com/DucTran999/cachekit/errors"
	"github.com/DucTran999/cachekit/remote"
)

// Cache defines a generic caching interface supporting basic key-value operations
// with optional TTL, key existence checking, and cache lifecycle management.
type Cache interface {
	// Get retrieves the value associated with the given key as a string.
	Get(ctx context.Context, key string) (string, error)

	// GetInto retrieves the value and unmarshals it into the provided destination.
	GetInto(ctx context.Context, key string, dest any) error

	// Set stores a value with the specified key and TTL (time to live).
	Set(ctx context.Context, key string, value any, ttl time.Duration) error

	// Del deletes one or more keys from the cache.
	Del(ctx context.Context, keys ...string) error

	// Has checks if the given key exists in the cache.
	Has(ctx context.Context, key string) (bool, error)

	// ExistingKeys returns a list of keys that currently exist among the provided ones.
	ExistingKeys(ctx context.Context, keys ...string) ([]string, error)

	// MissingKeys returns the subset of keys that do not exist.
	MissingKeys(ctx context.Context, keys ...string) ([]string, error)

	// Expire updates the TTL of the given key.
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// TTL returns the remaining TTL (in seconds) for the given key.
	TTL(ctx context.Context, key string) (int64, error)

	// FlushAll removes all keys from the cache.
	FlushAll(ctx context.Context) error
}

type RemoteCache interface {
	Cache

	// Ping checks the health/status of the cache backend.
	Ping(ctx context.Context) error

	// Close releases any resources held by the cache (e.g., closes connections).
	Close() error
}

type RedisConfig = config.RedisConfig

// Configuration-related errors
var (
	ErrMissingHost = cacheerr.ErrMissingHost
	ErrInvalidPort = cacheerr.ErrInvalidPort
	ErrInvalidDB   = cacheerr.ErrInvalidDB
)

// Operation/runtime errors
var (
	ErrKeyNotFound    = cacheerr.ErrKeyNotFound
	ErrDecode         = cacheerr.ErrDecode
	ErrSetNil         = cacheerr.ErrSetNil
	ErrSerializeValue = cacheerr.ErrSerializeValue
)

func NewRedisCache(config RedisConfig) (RemoteCache, error) {
	return remote.NewRedisCache(config)
}
