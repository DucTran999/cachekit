package remote_test

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/DucTran999/cachekit"
	"github.com/DucTran999/cachekit/config"
	cacheerr "github.com/DucTran999/cachekit/errors"
	"github.com/DucTran999/cachekit/test/mocks"
	"github.com/DucTran999/cachekit/test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type redisCacheUT struct {
	redisInst cachekit.Cache
}

func GetRedisInstance(t *testing.T) cachekit.Cache {
	err := utils.LoadEnv(".env.local")
	require.NoError(t, err)

	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	require.NoError(t, err)

	cfg := config.RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     port,
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	cache, err := cachekit.NewRedisCache(cfg)
	require.NoError(t, err)

	return cache
}

func TestRedisConnectFailed(t *testing.T) {
	t.Run("missing port in config", func(t *testing.T) {
		cfg := config.RedisConfig{
			Host: "localhost",
			DB:   1,
		}
		cache, err := cachekit.NewRedisCache(cfg)
		require.ErrorIs(t, err, cacheerr.ErrInvalidPort)

		assert.Nil(t, cache)
	})

	t.Run("mistake close connection", func(t *testing.T) {
		cache := GetRedisInstance(t)
		// Close be for ping
		cache.Close()

		err := cache.Ping(context.Background())
		require.NotNil(t, err)
	})
}

func TestRedis(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	ctx := context.Background()
	cache := GetRedisInstance(t)
	defer cache.Close()

	// Define test data
	user := &User{"daniel"}
	userJSON, err := json.Marshal(user)
	require.NoError(t, err)

	type badStruct struct {
		chanel chan int
	}

	testCases := []struct {
		name        string
		key         string
		value       any
		expectedErr error
	}{
		{"string", "test:redis:string", "hello", nil},
		{"int", "test:redis:int", 5, nil},
		{"struct", "test:redis:model", user, nil},
		{"bytes", "test:redis:bytes", userJSON, nil},
		{"binary", "test:redis:binary", mocks.BinaryVal{Data: "hello"}, nil},
		{"nil", "test:redis:nil", nil, cacheerr.ErrSetNil},
		{"non-serialize", "test:redis:serialize-error", mocks.BadBinary{}, cacheerr.ErrSerializeValue},
	}

	// Set test values
	for _, tc := range testCases {
		t.Run("Set "+tc.name, func(t *testing.T) {
			err := cache.Set(ctx, tc.key, tc.value, time.Minute)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}

	// Get & verify string value
	t.Run("Get string value", func(t *testing.T) {
		got, err := cache.Get(ctx, "test:redis:string")
		require.NoError(t, err)
		require.Equal(t, "hello", got)
	})

	t.Run("Get struct value", func(t *testing.T) {
		user := User{}
		err := cache.GetInto(ctx, "test:redis:model", &user)
		require.NoError(t, err)
		require.Equal(t, "daniel", user.Name)
	})

	t.Run("Get struct not found key", func(t *testing.T) {
		user := User{}
		err := cache.GetInto(ctx, "test:redis:model-not-found", &user)
		require.ErrorIs(t, err, cacheerr.ErrKeyNotFound)
	})

	t.Run("Get model unmarshal error", func(t *testing.T) {
		user := User{}
		err := cache.GetInto(ctx, "test:redis:string", &user)
		require.ErrorIs(t, err, cacheerr.ErrDecode)
	})

	// Has
	t.Run("Has string key", func(t *testing.T) {
		exists, err := cache.Has(ctx, "test:redis:string")
		require.NoError(t, err)
		require.True(t, exists)
	})

	// Delete
	t.Run("Delete string key", func(t *testing.T) {
		err := cache.Del(ctx, "test:redis:string")
		require.NoError(t, err)
	})

	// Get after delete
	t.Run("Get deleted key", func(t *testing.T) {
		_, err := cache.Get(ctx, "test:redis:string")
		require.Error(t, err)
	})
}

func TestRedisErrorWhileRunning(t *testing.T) {
	ctx := context.Background()
	cache := GetRedisInstance(t)
	cache.Close()

	key := "test:redis:key"
	value := "hello"

	// Set
	err := cache.Set(ctx, key, value, 1*time.Minute)
	require.NotNil(t, err)

	// Get
	_, err = cache.Get(ctx, key)
	require.NotNil(t, err)

	// Get Into
	err = cache.GetInto(ctx, key, struct{}{})
	require.NotNil(t, err)

	// Has
	_, err = cache.Has(ctx, key)
	require.NotNil(t, err)

	// Delete
	err = cache.Del(ctx, key)
	require.NotNil(t, err)

	// Get after delete
	_, err = cache.Get(ctx, key)
	require.Error(t, err)
}
