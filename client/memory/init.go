package memory

import (
    "time"
)

const (
    defaultGcMaxItems = 1000
    defaultGcInterval = 60 * time.Second
    defaultExpire     = 60 * time.Second

    minGcInterval = 10 * time.Second
    maxGcInterval = 600 * time.Second

    errSetProp = "memory: failed to set %s, %s"
)
