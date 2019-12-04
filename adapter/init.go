package adapter

import (
    "net/url"
    "time"
)

// memory
const (
    DefaultMemoryId = "memory"
    MemoryClass     = "@pgo2/adapter/Memory"
)

// http
const (
    DefaultHttpId = "http"
    HttpClass     = "@pgo2/adapter/Http"
)

// maxMind
const (
    DefaultMaxMindId = "maxMind"
    MaxMindClass     = "@pgo2/adapter/MaxMind"
)

// memCache
const (
    DefaultMemCacheId = "memCache"
    MemCacheClass = "@pgo2/adapter/MemCache"

)

// mongo
const (
    MongoClass = "@pgo2/adapter/Mongo"

    DefaultMongoId = "mongo"
    ErrInvalidOpt = "mongo: invalid option "
)

// redis
const (
    RedisClass = "@pgo2/adapter/Redis"
    DefaultRedisId = "redis"
)

// rabbitMq
const (
    RabbitMqClass = "@pgo2/adapter/RabbitMq"
    DefaultRabbitId = "rabbitMq"
)

// db
const (
    DbClass = "@pgo2/adapter/Db"
    DefaultDbTimeout     = 10 * time.Second
    DefaultDbId = "db"
)

func baseUrl(addr string) string {
    u, e := url.Parse(addr)
    if e != nil {
        panic("http parse url failed, " + e.Error())
    }

    return u.Scheme + "://" + u.Host + u.Path
}

func panicErr(err error){
    if err != nil {
        panic(err)
    }
}