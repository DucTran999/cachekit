# cachekit

[![Go Report Card](https://goreportcard.com/badge/github.com/DucTran999/cachekit)](https://goreportcard.com/report/github.com/DucTran999/cachekit)
[![Go](https://img.shields.io/badge/Go-1.23-blue?logo=go)](https://golang.org)
[![codecov](https://codecov.io/gh/DucTran999/cachekit/graph/badge.svg?token=5XBMMBKCPD)](https://codecov.io/gh/DucTran999/cachekit)
[![Known Vulnerabilities](https://snyk.io/test/github/ductran999/cachekit/badge.svg)](https://snyk.io/test/github/ductran999/cachekit)
[![License](https://img.shields.io/github/license/DucTran999/cachekit)](LICENSE)

---

**`cachekit`** is a unified caching interface for Go that supports both **in-memory** (e.g. Ristretto) and **remote** (e.g. Redis) backends.  
It provides a clean abstraction for setting, getting, and managing cache values using a common API.

---

## ✨ Features

- 🔄 Unified interface for remote/local caches
- 🔐 Redis support (via `go-redis/v9`)
- 🧠 In-memory support (via `dgraph-io/ristretto`)
- 🧪 Well-tested with real and mock backends
- ✅ Graceful error handling and type-safe GetInto()

---

## 📦 Installation

```bash
go get github.com/DucTran999/cachekit
```

## 🚀 Quick Start

Setup redis with docker compose

```bash
# run task start redis with docker compose
task testenv
```

```bash
# Also available script with make
make testenv
```

Example cache with redis. See more in `examples`

```go
package main

import (
	"log"

	"github.com/DucTran999/cachekit"
	"github.com/DucTran999/cachekit/config"
)

func main() {
	cfg := config.RedisConfig{
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
```

## 📜 License

This project is licensed under the MIT License.

## 🙌 Contributions

Contributions are welcome! Please open an issue or submit a pull request.
