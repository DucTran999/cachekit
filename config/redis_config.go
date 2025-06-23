package config

import (
	"net"
	"strconv"

	cacheerr "github.com/DucTran999/cachekit/errors"
)

type RedisConfig struct {
	Host string
	Port int

	Username string
	Password string
	DB       int
}

func (c RedisConfig) Validate() error {
	if c.Host == "" {
		return cacheerr.ErrMissingHost
	}
	if c.Port <= 0 || c.Port > 65535 {
		return cacheerr.ErrInvalidPort
	}
	if c.DB < 0 {
		return cacheerr.ErrInvalidDB
	}

	return nil
}

func (c *RedisConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}
