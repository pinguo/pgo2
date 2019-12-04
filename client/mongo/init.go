package mongo

import "time"

const (

    defaultDsn         = "mongodb://127.0.0.1:27017/"
    defaultOptions     = "connect=replicaSet&maxPoolSize=100&minPoolSize=1&maxIdleTimeMS=300000" +
        "&ssl=false&w=1&j=false&wtimeoutMS=10000&readPreference=secondaryPreferred"

    defaultConnectTimeout = 1 * time.Second
    defaultReadTimeout    = 10 * time.Second
    defaultWriteTimeout   = 10 * time.Second

    errSetProp    = "mongo: failed to set %s, %s"
    errInvalidDsn = "mongo: invalid dsn %s, %s"

    errDialFailed = "mongo: failed to dial %s, %s"
)