package config_test

import (
	"testing"

	"github.com/DucTran999/cachekit/config"
	cacheerr "github.com/DucTran999/cachekit/errors"
	"github.com/stretchr/testify/assert"
)

func TestRedisConfigValidate(t *testing.T) {
	type testcase struct {
		name        string
		cfg         config.RedisConfig
		expectedErr error
	}

	testTable := []testcase{
		{
			name: "invalid host",
			cfg: config.RedisConfig{
				Port: 7439,
			},
			expectedErr: cacheerr.ErrMissingHost,
		},
		{
			name: "invalid post",
			cfg: config.RedisConfig{
				Host: "localhost",
			},
			expectedErr: cacheerr.ErrInvalidPort,
		},
		{
			name: "invalid DB",
			cfg: config.RedisConfig{
				Host: "localhost",
				Port: 7439,
				DB:   -1,
			},
			expectedErr: cacheerr.ErrInvalidDB,
		},
		{
			name: "valid config",
			cfg: config.RedisConfig{
				Host: "localhost",
				Port: 6379,
				DB:   1,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestRedisAddress(t *testing.T) {
	expected := "localhost:6379"
	cfg := config.RedisConfig{
		Host: "localhost",
		Port: 6379,
		DB:   1,
	}

	assert.Equal(t, expected, cfg.Address())
}
