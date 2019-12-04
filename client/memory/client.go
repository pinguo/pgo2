package memory

import (
    "fmt"
    "sync"
    "time"

    "github.com/pinguo/pgo2/core"
    "github.com/pinguo/pgo2/util"
    "github.com/pinguo/pgo2/value"
)

type item struct {
    value  interface{}
    expire time.Time
}

func (i item) isExpired() bool {
    return !i.expire.IsZero() && time.Since(i.expire) > 0
}

// Memory Client component, configuration:
// components:
//      memory:
//          gcInterval: "60s"
//          gcMaxItems: 1000
func New(config map[string]interface{}) (interface{}, error) {
    memory := &Client{
        items:      make(map[string]*item),
        gcInterval: defaultGcInterval,
        gcMaxItems: defaultGcMaxItems,
    }

    err := core.ClientConfigure(memory, config)
    if err != nil {
        return nil, err
    }
    memory.Init()

    return memory, nil
}

type Client struct {
    lock       sync.RWMutex
    items      map[string]*item
    gcInterval time.Duration
    gcMaxItems int
}

func (c *Client) Init() {
    go c.gcLoop()
}

func (c *Client) SetGcInterval(v string) error {
    if gcInterval, e := time.ParseDuration(v); e != nil {
        return fmt.Errorf(errSetProp, "gcInterval", e.Error())
    } else {
        c.gcInterval = gcInterval
    }

    return nil
}

func (c *Client) SetGcMaxItems(gcMaxItems int) {
    if gcMaxItems > 0 {
        c.gcMaxItems = gcMaxItems
    }
}

func (c *Client) Get(key string) *value.Value {
    c.lock.RLock()
    defer c.lock.RUnlock()

    if item := c.items[key]; item != nil && !item.isExpired() {
        return value.New(item.value)
    }

    return value.New(nil)
}

func (c *Client) MGet(keys []string) map[string]*value.Value {
    c.lock.RLock()
    defer c.lock.RUnlock()

    result := make(map[string]*value.Value)
    for _, key := range keys {
        if item := c.items[key]; item != nil && !item.isExpired() {
            result[key] = value.New(item.value)
        } else {
            result[key] = value.New(nil)
        }
    }

    return result
}

func (c *Client) Set(key string, value interface{}, expire ...time.Duration) bool {
    c.lock.Lock()
    defer c.lock.Unlock()

    expire, now := append(expire, defaultExpire), time.Now()
    c.items[key] = &item{
        value:  value,
        expire: now.Add(expire[0]),
    }

    return true
}

func (c *Client) MSet(items map[string]interface{}, expire ...time.Duration) bool {
    c.lock.Lock()
    defer c.lock.Unlock()

    expire, now := append(expire, defaultExpire), time.Now()
    for key, value := range items {
        c.items[key] = &item{
            value:  value,
            expire: now.Add(expire[0]),
        }
    }
    return true
}

func (c *Client) Add(key string, value interface{}, expire ...time.Duration) bool {
    c.lock.Lock()
    defer c.lock.Unlock()

    expire, now := append(expire, defaultExpire), time.Now()
    if old := c.items[key]; old == nil || old.isExpired() {
        c.items[key] = &item{
            value:  value,
            expire: now.Add(expire[0]),
        }
        return true
    }
    return false
}

func (c *Client) MAdd(items map[string]interface{}, expire ...time.Duration) bool {
    c.lock.Lock()
    defer c.lock.Unlock()

    expire, now, success := append(expire, defaultExpire), time.Now(), 0
    for key, value := range items {
        if old := c.items[key]; old == nil || old.isExpired() {
            c.items[key] = &item{
                value:  value,
                expire: now.Add(expire[0]),
            }
            success++
        }
    }

    return success == len(items)
}

func (c *Client) Del(key string) bool {
    c.lock.Lock()
    defer c.lock.Unlock()

    if _, ok := c.items[key]; ok {
        delete(c.items, key)
        return true
    }
    return false
}

func (c *Client) MDel(keys []string) bool {
    c.lock.Lock()
    defer c.lock.Unlock()

    success := 0
    for _, key := range keys {
        if _, ok := c.items[key]; ok {
            delete(c.items, key)
            success++
        }
    }

    return success == len(keys)
}

func (c *Client) Exists(key string) bool {
    c.lock.RLock()
    defer c.lock.RUnlock()

    _, ok := c.items[key]
    return ok
}

func (c *Client) Incr(key string, delta int) int {
    c.lock.Lock()
    defer c.lock.Unlock()

    cur := c.items[key]
    if cur == nil {
        cur = &item{value: 0}
        c.items[key] = cur
    }

    newVal := util.ToInt(cur.value) + delta
    cur.value = newVal
    return newVal
}

func (c *Client) gcLoop() {
    if c.gcInterval < minGcInterval || c.gcInterval > maxGcInterval {
        c.gcInterval = defaultGcInterval
    }

    for {
        <-time.After(c.gcInterval)
        if expiredKeys := c.getExpireKeys(); len(expiredKeys) > 0 {
            c.clearExpiredKeys(expiredKeys)
        }
    }
}

func (c *Client) getExpireKeys() []string {
    c.lock.RLock()
    defer c.lock.RUnlock()

    keys, now := make([]string, 0), time.Now()
    for key, item := range c.items {
        if !item.expire.IsZero() && item.expire.Sub(now) < 0 {
            keys = append(keys, key)
            if len(keys) >= c.gcMaxItems {
                break
            }
        }
    }
    return keys
}

func (c *Client) clearExpiredKeys(keys []string) {
    c.lock.Lock()
    defer c.lock.Unlock()

    for _, key := range keys {
        delete(c.items, key)
    }
}
