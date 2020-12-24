package mongodb

import "time"

// 基于mgo封装

const (
	defaultDsn     = "mongodb://127.0.0.1:27017/"
	defaultOptions = "directConnection=false&maxPoolSize=100&minPoolSize=1&maxIdleTimeMS=300000" +
		"&ssl=false&w=1&j=false&wtimeoutMS=10000&readPreference=secondaryPreferred"

	defaultConnectTimeout = 1 * time.Second
	defaultReadTimeout    = 10 * time.Second
	defaultWriteTimeout   = 10 * time.Second

	errSetProp    = "mongodb: failed to set %s, %s"
	errInvalidDsn = "mongodb: invalid dsn %s, %s"
	errInvalidConfig = "mongodb: invalid config %s, %s"

	errDialFailed = "mongodb: failed to dial %s, %s"
)
