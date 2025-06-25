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

// Get retrieves the string value associated with the given key from Redis.
// Returns ErrKeyNotFound if the key does not exist.
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

// Set stores the given value in Redis under the specified key with an optional expiration.
// The value can be a string, []byte, encoding.BinaryMarshaler, or any type that can be JSON-marshaled.
// Returns ErrSetNil if the value is nil, or ErrSerializeValue if encoding fails.
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

// GetInto retrieves the cached value for the given key and unmarshals it into dest.
// Returns ErrKeyNotFound if the key doesn't exist, or ErrDecode if unmarshaling fails.
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

// Del deletes one or more keys from the Redis cache.
func (r *redisCache) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Has checks whether the given key exists in the Redis cache.
func (r *redisCache) Has(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// Ping checks the connection to the Redis server to ensure it is reachable.
func (r *redisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close shuts down the Redis client and releases any associated resources.
func (r *redisCache) Close() error {
	return r.client.Close()
}

// Expire sets a new expiration time for the given key in the Redis cache.
func (r *redisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// TTL returns the time-to-live (in seconds) for the given key in Redis.
// A return value of -2 indicates the key does not exist (with ErrKeyNotFound).
// A return value of -1 indicates the key exists but has no expiration set.
// Otherwise, it returns the remaining TTL in seconds.
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

// ExistingKeys checks which of the provided keys exist in the Redis cache.
// It returns a slice of keys that currently have associated values.
// Non-existent keys are filtered out using MGet and nil checking.
func (r *redisCache) ExistingKeys(ctx context.Context, keys ...string) ([]string, error) {
	vals, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	// Filter keys that exist
	keyExists := make([]string, 0, len(vals))
	for i, val := range vals {
		if val != nil {
			keyExists = append(keyExists, keys[i])
		}
	}

	return keyExists, nil
}

// FlushAll deletes all keys in the Redis cache across all databases.
// Use with caution as this will clear the entire Redis instance.
func (r *redisCache) FlushAll(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

// MissingKeys returns a list of keys that do not exist in the Redis cache.
// It uses MGet to check the existence of each key and filters out those with nil values.
func (r *redisCache) MissingKeys(ctx context.Context, keys ...string) ([]string, error) {
	vals, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	// Filter out keys that are missing (i.e., have nil value)
	var missingKeys []string
	for i, val := range vals {
		if val == nil {
			missingKeys = append(missingKeys, keys[i])
		}
	}

	return missingKeys, nil
}
