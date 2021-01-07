package adapter

import (
	"time"

	"github.com/pinguo/pgo2"
	"github.com/pinguo/pgo2/client/redis"
	"github.com/pinguo/pgo2/iface"
	"github.com/pinguo/pgo2/util"
	"github.com/pinguo/pgo2/value"
)

var RedisClass string
func init() {
	container := pgo2.App().Container()
	RedisClass = container.Bind(&Redis{})
}

// NewRedis of Redis Client, add context support.
// usage: redis := this.GetObject(adapter.NewRedis()).(*adapter.Redis)
func NewRedis(componentId ...string) *Redis {
	id := DefaultRedisId
	if len(componentId) > 0 {
		id = componentId[0]
	}

	c := &Redis{}
	c.client = pgo2.App().Component(id, redis.New, map[string]interface{}{"logger":pgo2.GLogger()}).(*redis.Client)
	c.panicRecover = true

	return c
}

// NewRedisPool of Redis Client from pool, add context support.
// usage: redis := this.GetObjPool(adapter.RedisClass, adapter.NewRedisPool).(*adapter.Redis)
// It is recommended to use : redis := this.GetObjBox(adapter.RedisClass).(*adapter.Redis)
func NewRedisPool(iObj iface.IObject, componentId ...interface{}) iface.IObject {

	return iObj
}

type Redis struct {
	pgo2.Object
	client       *redis.Client
	panicRecover bool
}

// GetObjPool, GetObjBox fetch is performed automatically
func (r *Redis) Prepare(componentId ...interface{}) {
	id := DefaultRedisId
	if len(componentId) > 0 {
		id = componentId[0].(string)
	}

	r.client = pgo2.App().Component(id, redis.New, map[string]interface{}{"logger":pgo2.GLogger()}).(*redis.Client)
	r.panicRecover = true

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

	res, err := r.client.Incr(key, int64(delta))
	r.parseErr(err)

	return int(res)
}

func (r *Redis) IncrBy(key string, delta int64) (int64, error) {
	profile := "Redis.IncrBy"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.Incr(key, delta)
	r.parseErr(err)

	return res,err
}

// 支持的命令请查阅：Redis.allRedisCmd
// args = [0:"key"]
// Example:
// redis := t.GetObjBox(Redis.AdapterClass).(*Redis.Adapter)
// retI := redis.Do("SADD","myTest", "test1"
// ret := retI.(int64)
// fmt.Println(ret) = 1
// retList :=redis.Do("SMEMBERS","myTest"
// retListI,_:=ret.([]interface{})
// for _,v:=range retListI{
//    vv :=value.New(v) // 写入的时候有value.Encode(),如果存入的是结构体或slice map 需要decode,其他类型直接断言类型
//    fmt.Println(vv.String()) // test1
// }
func (r *Redis) Do(cmd string, args ...interface{}) interface{} {
	profile := "Redis.Do." + cmd
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.Do(cmd, args...)
	r.parseErr(err)

	return res
}

func (r *Redis) ExpireAt(key string, timestamp int64) bool {
	profile := "Redis.ExpireAt"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.ExpireAt(key, timestamp)
	r.parseErr(err)

	return res
}

func (r *Redis) Expire(key string, expire time.Duration) bool {
	profile := "Redis.Expire"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.Expire(key, expire)
	r.parseErr(err)

	return res
}

func (r *Redis) RPush(key string, values ...interface{}) bool {
	profile := "Redis.RPush"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.RPush(key, values...)
	r.parseErr(err)

	return res
}

func (r *Redis) LPush(key string, values ...interface{}) bool {
	profile := "Redis.LPush"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.LPush(key, values...)
	r.parseErr(err)

	return res
}

func (r *Redis) RPop(key string) *value.Value {
	profile := "Redis.RPop"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.RPop(key)
	r.parseErr(err)

	hit := 0
	if res != nil && res.Valid() {
		hit += 1
	}

	r.Context().Counting(profile, hit, 1)

	return res
}

func (r *Redis) LPop(key string) *value.Value {
	profile := "Redis.LPop"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.LPop(key)
	r.parseErr(err)

	hit := 0
	if res != nil && res.Valid() {
		hit += 1
	}

	r.Context().Counting(profile, hit, 1)

	return res
}

func (r *Redis) LLen(key string) int64 {
	profile := "Redis.LLen"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.LLen(key)
	r.parseErr(err)

	return res
}

func (r *Redis) HDel(key string, fields ...interface{}) int64 {
	profile := "Redis.HDel"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.HDel(key,fields...)
	r.parseErr(err)

	return res
}

func (r *Redis) HExists(key, field string) bool {
	profile := "Redis.HExists"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.HExists(key, field)
	r.parseErr(err)

	return res
}

func (r *Redis) HSet(key string, fv ...interface{}) bool {
	profile := "Redis.HSet"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.HSet(key,fv...)
	r.parseErr(err)

	return res
}

func (r *Redis) HGet(key,field string) *value.Value {
	profile := "Redis.HGet"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.HGet(key,field)
	r.parseErr(err)
	hit := 0
	if res != nil && res.Valid() {
		hit = 1
	}

	r.Context().Counting(profile, hit, 1)
	return res
}

func (r *Redis) HGetAll(key string)map[string]*value.Value {
	profile := "Redis.HGetAll"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.HGetAll(key)
	r.parseErr(err)

	return res
}

func (r *Redis) HMSet(key string, fv ...interface{}) bool {
	profile := "Redis.HMSet"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.HMSet(key,fv...)
	r.parseErr(err)

	return res
}

func (r *Redis) HMGet(key string,fields ...interface{}) map[string]*value.Value {
	profile := "Redis.HMGet"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.HMGet(key,fields...)
	r.parseErr(err)
	hit := 0
	for _, v := range res {
		if v != nil && v.Valid() {
			hit += 1
		}
	}

	r.Context().Counting(profile, hit, len(fields))
	return res
}

func (r *Redis) HIncrBy(key, field string, delta int64) (int64,error) {
	profile := "Redis.HIncrBy"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.HIncrBy(key,field, delta)
	r.parseErr(err)

	return res,err
}

func (r *Redis) ZRange(key string, start, end int) []*value.Value {
	profile := "Redis.ZRange"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.ZRange(key,start, end)
	r.parseErr(err)

	return res
}


func (r *Redis) ZRevRange(key string, start, end int) []*value.Value {
	profile := "Redis.ZRevRange"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.ZRevRange(key,start, end)
	r.parseErr(err)

	return res
}

func (r *Redis) ZRangeWithScores(key string, start, end int) []*redis.ZV {
	profile := "Redis.ZRangeWithScores"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.ZRangeWithScores(key,start, end)
	r.parseErr(err)

	return res
}


func (r *Redis) ZRevRangeWithScores(key string, start, end int) []*redis.ZV {
	profile := "Redis.ZRevRangeWithScores"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.ZRevRangeWithScores(key,start, end)
	r.parseErr(err)

	return res
}

func (r *Redis) ZAdd(key string, members ...*redis.Z) int64 {
	profile := "Redis.ZAdd"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.ZAdd(key,members...)
	r.parseErr(err)

	return res
}

func (r *Redis) ZAddOpt(key string,opts []string, members ...*redis.Z) (int64,error) {
	profile := "Redis.ZAddOpt"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.ZAddOpt(key,opts,members...)
	r.parseErr(err)

	return res,err
}

func (r *Redis) ZCard(key string) int64 {
	profile := "Redis.ZCard"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.ZCard(key)
	r.parseErr(err)

	return res
}


func (r *Redis) ZRem(key string,members ...interface{}) int64 {
	profile := "Redis.ZRem"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.ZRem(key,members...)
	r.parseErr(err)

	return res
}


func (r *Redis) parseErr(err error) {
	if err != nil {
		panic(err)
	}
}
