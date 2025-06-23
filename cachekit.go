package cachekit

import (
	"context"
	"time"

	"github.com/DucTran999/cachekit/config"
	cacheerr "github.com/DucTran999/cachekit/errors"
	"github.com/DucTran999/cachekit/remote"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	GetInto(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Has(ctx context.Context, key string) (bool, error)

	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (int64, error)

	Ping(ctx context.Context) error
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

func NewRedisCache(config RedisConfig) (Cache, error) {
	return remote.NewRedisCache(config)
}
