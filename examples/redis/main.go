package main

import (
	"log"

	"github.com/DucTran999/cachekit"
)

func main() {
	cfg := cachekit.RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Username: "default",
		Password: "example",
		DB:       1,
	}

	cache, err := cachekit.NewRedisCache(cfg)
	if err != nil {
		log.Fatalf("[FATAL] failed to connect redis: %v", err)
	}

	log.Println("[INFO] redis connected")
	defer cache.Close()
}
