package remote

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DucTran999/cachekit/config"
	cacheerr "github.com/DucTran999/cachekit/errors"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(config config.RedisConfig) (*redisCache, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate redis config: %w", err)
	}

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Address(),
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,

		// Set connection pool options
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &redisCache{client: rdb}, nil
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		// Match the error pattern used in the in-memory implementation for consistency
		if err == redis.Nil {
			return "", fmt.Errorf("%w: %s", cacheerr.ErrKeyNotFound, key)
		}

		return "", err
	}

	return val, nil
}

func (r *redisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if value == nil {
		return cacheerr.ErrSetNil
	}

	var b []byte
	var err error

	switch v := value.(type) {
	case encoding.BinaryMarshaler:
		b, err = v.MarshalBinary()
	case string:
		b = []byte(v)
	case []byte:
		b = v
	default:
		b, err = json.Marshal(v) // fallback to JSON
	}
	if err != nil {
		return fmt.Errorf("%w: %w", cacheerr.ErrSerializeValue, err)
	}

	return r.client.Set(ctx, key, string(b), expiration).Err()
}

func (r *redisCache) GetInto(ctx context.Context, key string, dest any) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("%w: %s", cacheerr.ErrKeyNotFound, key)
		}
		return err
	}

	// Unmarshal JSON into dest
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("%w: key=%q", cacheerr.ErrDecode, key)
	}

	return nil
}

func (r *redisCache) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *redisCache) Has(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (r *redisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *redisCache) Close() error {
	return r.client.Close()
}

func (r *redisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

func (r *redisCache) TTL(ctx context.Context, key string) (int64, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return -1, err // return the actual error for clarity
	}

	if ttl < 0 {
		// -2: key does not exist
		// -1: key exists but has no expiration
		if ttl == -2 {
			return -2, cacheerr.ErrKeyNotFound
		}

		return -1, nil
	}

	return int64(ttl.Seconds()), nil
}
