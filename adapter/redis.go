package adapter

import (
    "time"

    "github.com/pinguo/pgo2"
    "github.com/pinguo/pgo2/client/redis"
    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/util"
    "github.com/pinguo/pgo2/value"
)

func init() {
    container := pgo2.App().Container()
    container.Bind(&Redis{})
}

// NewRedis of Redis Client, add context support.
// usage: redis := this.GetObject(adapter.NewRedis()).(*adapter.Redis)
func NewRedis(componentId ...string) *Redis {
    id := DefaultRedisId
    if len(componentId) > 0 {
        id = componentId[0]
    }

    c := &Redis{}
    c.client = pgo2.App().Component(id, redis.New).(*redis.Client)
    c.panicRecover = true

    return c
}

// NewRedisPool of Redis Client from pool, add context support.
// usage: redis := this.GetObjPool(adapter.NewRedisPool).(*adapter.Redis)
func NewRedisPool(ctr iface.IContext, componentId ...interface{}) iface.IObject {
    id := DefaultRedisId
    if len(componentId) > 0 {
        id = componentId[0].(string)
    }

    c := pgo2.App().GetObjPool(RedisClass, ctr).(*Redis)

    c.client = pgo2.App().Component(id, redis.New).(*redis.Client)
    c.panicRecover = true

    return c
}

type Redis struct {
    pgo2.Object
    client       *redis.Client
    panicRecover bool
}

func (r *Redis) SetPanicRecover(v bool) {
    r.panicRecover = v
}

func (r *Redis) GetClient() *redis.Client {
    return r.client
}

func (r *Redis) handlePanic() {
    if r.panicRecover {
        if v := recover(); v != nil {
            r.Context().Error(util.ToString(v))
        }
    }
}

func (r *Redis) Get(key string) *value.Value {
    profile := "Redis.Get"
    r.Context().ProfileStart(profile)
    defer r.Context().ProfileStop(profile)
    defer r.handlePanic()

    res, err := r.client.Get(key)
    r.parseErr(err)
    hit := 0
    if res != nil && res.Valid() {
        hit = 1
    }

    r.Context().Counting(profile, hit, 1)
    return res
}

func (r *Redis) MGet(keys []string) map[string]*value.Value {
    profile := "Redis.MGet"
    r.Context().ProfileStart(profile)
    defer r.Context().ProfileStop(profile)
    defer r.handlePanic()

    res, err := r.client.MGet(keys)
    r.parseErr(err)
    hit := 0
    for _, v := range res {
        if v != nil && v.Valid() {
            hit += 1
        }
    }

    r.Context().Counting(profile, hit, len(keys))
    return res
}

func (r *Redis) Set(key string, value interface{}, expire ...time.Duration) bool {
    profile := "Redis.Set"
    r.Context().ProfileStart(profile)
    defer r.Context().ProfileStop(profile)
    defer r.handlePanic()

    res, err := r.client.Set(key, value, expire...)
    r.parseErr(err)

    return res
}

func (r *Redis) MSet(items map[string]interface{}, expire ...time.Duration) bool {
    profile := "Redis.MSet"
    r.Context().ProfileStart(profile)
    defer r.Context().ProfileStop(profile)
    defer r.handlePanic()

    res, err := r.client.MSet(items, expire...)
    r.parseErr(err)

    return res
}

func (r *Redis) Add(key string, value interface{}, expire ...time.Duration) bool {
    profile := "Redis.Add"
    r.Context().ProfileStart(profile)
    defer r.Context().ProfileStop(profile)
    defer r.handlePanic()

    res, err := r.client.Add(key, value, expire...)
    r.parseErr(err)

    return res
}

func (r *Redis) MAdd(items map[string]interface{}, expire ...time.Duration) bool {
    profile := "Redis.MAdd"
    r.Context().ProfileStart(profile)
    defer r.Context().ProfileStart(profile)
    defer r.handlePanic()

    res, err := r.client.MAdd(items, expire...)
    r.parseErr(err)

    return res
}

func (r *Redis) Del(key string) bool {
    profile := "Redis.Del"
    r.Context().ProfileStart(profile)
    defer r.Context().ProfileStop(profile)
    defer r.handlePanic()

    res, err := r.client.Del(key)
    r.parseErr(err)

    return res
}

func (r *Redis) MDel(keys []string) bool {
    profile := "Redis.MDel"
    r.Context().ProfileStart(profile)
    defer r.Context().ProfileStop(profile)
    defer r.handlePanic()

    res, err := r.client.MDel(keys)
    r.parseErr(err)

    return res
}

func (r *Redis) Exists(key string) bool {
    profile := "Redis.Exists"
    r.Context().ProfileStart(profile)
    defer r.Context().ProfileStop(profile)
    defer r.handlePanic()

    res, err := r.client.Exists(key)
    r.parseErr(err)

    return res
}

func (r *Redis) Incr(key string, delta int) int {
    profile := "Redis.Incr"
    r.Context().ProfileStart(profile)
    defer r.Context().ProfileStop(profile)
    defer r.handlePanic()

    res, err := r.client.Incr(key, delta)
    r.parseErr(err)

    return res
}

// 支持的命令请查阅：Redis.allRedisCmd
// args = [0:"key"]
// Example:
// redis := t.GetObject(Redis.AdapterClass).(*Redis.Adapter)
// retI := redis.Do("SADD","myTest", "test1"
// ret := retI.(int)
// fmt.Println(ret) = 1
// retList :=redis.Do("SMEMBERS","myTest"
// retListI,_:=ret.([]interface{})
// for _,v:=range retListI{
//    vv :=pgo.NewValue(v) // 写入的时候有pgo.Encode(),如果存入的是结构体或slice map 需要decode,其他类型直接断言类型
//    fmt.Println(vv.String()) // test1
// }
func (r *Redis) Do(cmd string, args ...interface{}) interface{} {
    profile := "Redis.Do." + cmd
    r.Context().ProfileStart(profile)
    defer r.Context().ProfileStop(profile)
    defer r.handlePanic()

    res, err := r.client.Do(cmd, args ...)
    r.parseErr(err)

    return res
}

func (r *Redis) parseErr(err error) {
    if err != nil {
        panic(err)
    }
}
