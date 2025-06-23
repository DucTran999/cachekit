package cachekit

import (
	"context"
	"time"

	"github.com/DucTran999/cachekit/config"
	"github.com/DucTran999/cachekit/remote"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	GetInto(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Has(ctx context.Context, key string) (bool, error)

	Ping(ctx context.Context) error
	Close() error
}

func NewRedisCache(config config.RedisConfig) (Cache, error) {
	return remote.NewRedisCache(config)
}
