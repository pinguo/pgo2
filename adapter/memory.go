package adapter

import (
    "fmt"
    "time"

    "github.com/pinguo/pgo2"
    "github.com/pinguo/pgo2/client/memory"
    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/util"
    "github.com/pinguo/pgo2/value"
)

func init() {
    container := pgo2.App().Container()
    container.Bind(&Memory{})
}

// NewMemory of Memory Client, add context support.
// usage: memory := this.GetObj(adapter.NewMemory()).(adapter.IMemory)/(*adapter.Memory)
func NewMemory(componentId ...string) *Memory {

    id := DefaultMemoryId
    if len(componentId) > 0 {
        id = componentId[0]
    }

    m := &Memory{}

    client := pgo2.App().Component(id, memory.New)
    m.client = client.(*memory.Client)
    m.panicRecover = true

    return m
}

// NewMemoryPool of Memory Client from object pool, add context support.
// usage: memory := this.GetObjPool(adapter.NewMemoryPool).(adapter.IMemory)/(*adapter.Memory)
func NewMemoryPool(ctr iface.IContext, componentId ...interface{}) iface.IObject {
    id := DefaultMemoryId
    if len(componentId) > 0 {
        fmt.Println(componentId[0])
        id = componentId[0].(string)
    }

    m := pgo2.App().GetObjPool(MemoryClass, ctr).(*Memory)

    client := pgo2.App().Component(id, memory.New)
    m.client = client.(*memory.Client)
    m.panicRecover = true

    return m
}

type Memory struct {
    pgo2.Object
    client       *memory.Client
    panicRecover bool
}

func (m *Memory) SetPanicRecover(v bool) {
    m.panicRecover = v
}

func (m *Memory) Client() *memory.Client {
    return m.client
}

func (m *Memory) handlePanic() {
    if m.panicRecover {
        if v := recover(); v != nil {
            m.Context().Error(util.ToString(v))
        }
    }
}

func (m *Memory) Get(key string) *value.Value {
    profile := "Memory.Get"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    res, hit := m.client.Get(key), 0
    if res != nil && res.Valid() {
        hit = 1
    }

    m.Context().Counting(profile, hit, 1)
    return res
}

func (m *Memory) MGet(keys []string) map[string]*value.Value {
    profile := "Memory.MGet"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    res, hit := m.client.MGet(keys), 0
    for _, v := range res {
        if v != nil && v.Valid() {
            hit += 1
        }
    }

    m.Context().Counting(profile, hit, len(keys))
    return res
}

func (m *Memory) Set(key string, value interface{}, expire ...time.Duration) bool {
    profile := "Memory.Set"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    return m.client.Set(key, value, expire...)
}

func (m *Memory) MSet(items map[string]interface{}, expire ...time.Duration) bool {
    profile := "Memory.MSet"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    return m.client.MSet(items, expire...)
}

func (m *Memory) Add(key string, value interface{}, expire ...time.Duration) bool {
    profile := "Memory.Add"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    return m.client.Add(key, value, expire...)
}

func (m *Memory) MAdd(items map[string]interface{}, expire ...time.Duration) bool {
    profile := "Memory.MAdd"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStart(profile)
    defer m.handlePanic()

    return m.client.MAdd(items, expire...)
}

func (m *Memory) Del(key string) bool {
    profile := "Memory.Del"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    return m.client.Del(key)
}

func (m *Memory) MDel(keys []string) bool {
    profile := "Memory.MDel"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    return m.client.MDel(keys)
}

func (m *Memory) Exists(key string) bool {
    profile := "Memory.Exists"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    return m.client.Exists(key)
}

func (m *Memory) Incr(key string, delta int) int {
    profile := "Memory.Incr"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    defer m.handlePanic()

    return m.client.Incr(key, delta)
}
