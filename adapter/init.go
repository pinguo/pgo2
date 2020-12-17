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

// orm
const (
	DefaultOrmId      = "orm"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
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

