package adapter

import (
	"net/url"
	"time"
)

// memory
const (
	DefaultMemoryId = "memory"
)

// http
const (
	DefaultHttpId = "http"
)

// maxMind
const (
	DefaultMaxMindId = "maxMind"
)

// memCache
const (
	DefaultMemCacheId = "memCache"
)

// mongo
const (
	DefaultMongoId = "mongo"
	ErrInvalidOpt  = "mongo: invalid option "
)

// redis
const (
	DefaultRedisId = "redis"
)

// rabbitMq
const (
	DefaultRabbitId = "rabbitMq"
)

// db
const (
	DefaultDbTimeout = 10 * time.Second
	DefaultDbId      = "db"
)

func baseUrl(addr string) string {
	u, e := url.Parse(addr)
	if e != nil {
		panic("http parse url failed, " + e.Error())
	}

	return u.Scheme + "://" + u.Host + u.Path
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
