package adapter

import (
    "time"

    "github.com/pinguo/pgo2"
    "github.com/pinguo/pgo2/client/memcache"
    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/util"
    "github.com/pinguo/pgo2/value"
)

func init() {
    container := pgo2.App().Container()
    container.Bind(&MemCache{})
}

// NewMemCache of MemCache Client, add context support.
// usage: mc := this.GetObj(adapter.New()).(adapter.IMemCache)/(*adapter.MemCache)
func NewMemCache(componentId ...string) *MemCache {
    id := DefaultMemCacheId
    if len(componentId) > 0 {
        id = componentId[0]
    }
    m := &MemCache{}

    m.client = pgo2.App().Component(id, memcache.New).(*memcache.Client)
    m.panicRecover = true

    return m
}

// NewMemCachePool of MemCache Client from pool, add context support.
// usage: memory := this.GetObjPool(adapter.NewMaxMindPool).(adapter.IMemCache)/(*adapter.MemCache)
func NewMemCachePool(ctr iface.IContext, componentId ...interface{}) iface.IObject {
    id := DefaultMemCacheId
    if len(componentId) > 0 {
        id = componentId[0].(string)
    }

    m := pgo2.App().GetObjPool(MemCacheClass, ctr).(*MemCache)

    m.client = pgo2.App().Component(id, memcache.New).(*memcache.Client)

    return m
}

type MemCache struct {
    pgo2.Object
    client       *memcache.Client
    panicRecover bool
}

func (m *MemCache) SetPanicRecover(v bool) {
    m.panicRecover = v
}


func (m *MemCache) GetClient() *memcache.Client {
    return m.client
}

func (m *MemCache) handlePanic() {
    if m.panicRecover {
        if v := recover(); v != nil {
            m.Context().Error(util.ToString(v))
        }
    }
}

func (m *MemCache) Get(key string) *value.Value {
    profile := "MemCache.Get"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    hit := 0
    res, err := m.client.Get(key)
    panicErr(err)
    if res != nil && res.Valid() {
        hit = 1
    }

    m.Context().Counting(profile, hit, 1)
    return res
}

func (m *MemCache) MGet(keys []string) map[string]*value.Value {
    profile := "MemCache.MGet"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    hit := 0
    res, err := m.client.MGet(keys)
    panicErr(err)
    for _, v := range res {
        if v != nil && v.Valid() {
            hit += 1
        }
    }

    m.Context().Counting(profile, hit, len(keys))
    return res
}

func (m *MemCache) Set(key string, value interface{}, expire ...time.Duration) bool {
    profile := "MemCache.Set"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    b, err := m.client.Set(key, value, expire...)
    panicErr(err)

    return b
}

func (m *MemCache) MSet(items map[string]interface{}, expire ...time.Duration) bool {
    profile := "MemCache.MSet"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    b,err := m.client.MSet(items, expire...)
    panicErr(err)

    return b
}

func (m *MemCache) Add(key string, value interface{}, expire ...time.Duration) bool {
    profile := "MemCache.Add"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    b, err := m.client.Add(key, value, expire...)
    panicErr(err)

    return b
}

func (m *MemCache) MAdd(items map[string]interface{}, expire ...time.Duration) bool {
    profile := "MemCache.MAdd"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStart(profile)
    defer m.handlePanic()

    b, err := m.client.MAdd(items, expire...)
    panicErr(err)

    return b
}

func (m *MemCache) Del(key string) bool {
    profile := "MemCache.Del"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    b, err := m.client.Del(key)
    panicErr(err)

    return b
}

func (m *MemCache) MDel(keys []string) bool {
    profile := "MemCache.MDel)"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    b, err := m.client.MDel(keys)
    panicErr(err)

    return b
}

func (m *MemCache) Exists(key string) bool {
    profile := "MemCache.Exists"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    b, err := m.client.Exists(key)
    panicErr(err)

    return b
}

func (m *MemCache) Incr(key string, delta int) int {
    profile := "MemCache.Incr"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    i, err := m.client.Incr(key, delta)
    panicErr(err)

    return i
}

func (m *MemCache) Retrieve(cmd, key string) *memcache.Item {
    profile := "MemCache.Retrieve"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    item, err := m.client.Retrieve(cmd, key)
    panicErr(err)

    return item
}

func (m *MemCache) MultiRetrieve(cmd string, keys []string) []*memcache.Item {
    profile := "MemCache.MultiRetrieve"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    items, err := m.client.MultiRetrieve(cmd, keys)
    panicErr(err)

    return items
}

func (m *MemCache) Store(cmd string, item *memcache.Item, expire ...time.Duration) bool {
    profile := "MemCache.Store"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    b, err := m.client.Store(cmd, item, expire...)
    panicErr(err)

    return b
}

func (m *MemCache) MultiStore(cmd string, items []*memcache.Item, expire ...time.Duration) bool {
    profile := "MemCache.MultiStore"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    b, err :=  m.client.MultiStore(cmd, items, expire...)
    panicErr(err)

    return b
}
