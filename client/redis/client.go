package redis

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pinguo/pgo2/core"
	"github.com/pinguo/pgo2/util"
	"github.com/pinguo/pgo2/value"
)

// Redis Client component, require redis-server 2.6.12+
// configuration:
// components:
//      redis:
//          prefix: "pgo2_"
//          password: ""
//          db: 0
//          maxIdleConn: 10
//          maxIdleTime: "60s"
//          netTimeout: "1s"
//          probInterval: "0s"
//          mod:"cluster"
//          servers:
//              - "127.0.0.1:6379"
//              - "127.0.0.1:6380"
func New(config map[string]interface{}) (interface{}, error) {

	c := &Client{}
	c.hashRing = util.NewHashRing()
	c.connLists = make(map[string]*connList)
	c.servers = make(map[string]*serverInfo)

	c.prefix = defaultPrefix
	c.password = defaultPassword
	c.db = defaultDb
	c.maxIdleConn = defaultIdleConn
	c.maxIdleTime = defaultIdleTime
	c.netTimeout = defaultTimeout
	c.probeInterval = defaultProbe
	c.mod = ModCluster

	if err := core.ClientConfigure(c, config); err != nil {
		return nil, err
	}

	if err := c.Init(); err != nil {
		return nil, err
	}

	return c, nil
}

type Client struct {
	Pool
}

func (c *Client) Get(key string) (*value.Value, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("GET", newKey)
	if err != nil {
		return nil, err
	}
	defer conn.Close(false)
	v, err := conn.Do("GET", newKey)
	if err != nil {
		return nil, err
	}
	return value.New(v), nil
}

func (c *Client) MGet(keys []string) (map[string]*value.Value, error) {
	result := make(map[string]*value.Value)
	addrKeys, newKeys, err := c.AddrNewKeys("MGET", keys)
	if err != nil {
		return result, err
	}

	lock, wg := new(sync.Mutex), new(sync.WaitGroup)
	var retErr error
	wg.Add(len(addrKeys))
	for addr, keys := range addrKeys {
		go c.RunAddrFunc(addr, keys, wg, func(conn *Conn, keys []string) {
			ret, err := conn.Do("MGET", keys2Args(keys)...)
			if err != nil {
				retErr = err
				return
			}
			if items, ok := ret.([]interface{}); ok {
				lock.Lock()
				defer lock.Unlock()
				for i, item := range items {
					oldKey := newKeys[keys[i]]
					result[oldKey] = value.New(item)
				}
			}
		})
	}

	wg.Wait()

	return result, retErr
}

func (c *Client) Set(key string, value interface{}, expire ...time.Duration) (bool, error) {
	expire = append(expire, defaultExpire)
	return c.set(key, value, expire[0], "")
}

func (c *Client) MSet(items map[string]interface{}, expire ...time.Duration) (bool, error) {
	expire = append(expire, defaultExpire)
	return c.mset(items, expire[0], "")
}

func (c *Client) Add(key string, value interface{}, expire ...time.Duration) (bool, error) {
	expire = append(expire, defaultExpire)
	return c.set(key, value, expire[0], "NX")
}

func (c *Client) MAdd(items map[string]interface{}, expire ...time.Duration) (bool, error) {
	expire = append(expire, defaultExpire)
	return c.mset(items, expire[0], "NX")
}

func (c *Client) Del(key string) (bool, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("DEL", newKey)
	if err != nil {
		return false, err
	}
	defer conn.Close(false)

	ret, err := conn.Do("DEL", newKey)
	if err != nil {
		return false, err
	}
	num, ok := ret.(int64)
	return ok && num == 1, nil
}

func (c *Client) MDel(keys []string) (bool, error) {
	addrKeys, _, err := c.AddrNewKeys("DEL", keys)
	if err != nil {
		return false, err
	}

	var retErr error
	wg, success := new(sync.WaitGroup), uint32(0)
	wg.Add(len(addrKeys))
	for addr, keys := range addrKeys {
		go c.RunAddrFunc(addr, keys, wg, func(conn *Conn, keys []string) {
			ret, err := conn.Do("DEL", keys2Args(keys)...)
			if err != nil {
				retErr = err
			}

			if num, ok := ret.(int64); ok && num > 0 {
				atomic.AddUint32(&success, uint32(num))
			}
		})
	}

	wg.Wait()
	return success == uint32(len(keys)), retErr
}

func (c *Client) Exists(key string) (bool, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("EXISTS", newKey)
	if err != nil {
		return false, err
	}

	defer conn.Close(false)
	ret, err := conn.Do("EXISTS", newKey)
	num, ok := ret.(int64)
	return ok && num == 1, err
}

func (c *Client) Incr(key string, delta int64) (int64, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("INCRBY", newKey)
	if err != nil {
		return 0, err
	}

	defer conn.Close(false)
	ret, err := conn.Do("INCRBY", newKey, delta)
	num, _ := ret.(int64)
	return num, err
}

func (c *Client) set(key string, value interface{}, expire time.Duration, flag string) (bool, error) {
	newKey := c.BuildKey(key)
	conn, errConn := c.GetConnByKey("SET", newKey)
	if errConn != nil {
		return false, errConn
	}

	defer conn.Close(false)

	var res interface{}
	var err error
	if len(flag) == 0 {
		res, err = conn.Do("SET", newKey, value, "EX", expire/time.Second)
	} else {
		res, err = conn.Do("SET", newKey, value, "EX", expire/time.Second, flag)
	}

	if err != nil {
		return false, err
	}

	payload, ok := res.([]byte)
	return ok && bytes.Equal(payload, replyOK), nil
}

func (c *Client) mset(items map[string]interface{}, expire time.Duration, flag string) (bool, error) {
	addrKeys, newKeys, errAddKey := c.AddrNewKeys("SET", items)
	if errAddKey != nil {
		return false, errAddKey
	}

	var err error
	wg, success := new(sync.WaitGroup), uint32(0)
	wg.Add(len(addrKeys))
	for addr, keys := range addrKeys {
		go c.RunAddrFunc(addr, keys, wg, func(conn *Conn, keys []string) {
			for _, key := range keys {
				if oldKey := newKeys[key]; len(flag) == 0 {
					err = conn.WriteCmd("SET", key, items[oldKey], "EX", expire/time.Second)
				} else {
					err = conn.WriteCmd("SET", key, items[oldKey], "EX", expire/time.Second, flag)
				}
				if err != nil {
					return
				}
			}

			for range keys {
				ret, err := conn.ReadReply()
				if err != nil {
					return
				}

				payload, ok := ret.([]byte)
				if ok && bytes.Equal(payload, replyOK) {
					atomic.AddUint32(&success, 1)
				}
			}
		})
	}

	wg.Wait()

	return success == uint32(len(items)), err
}

func (c *Client) ExpireAt(key string, timestamp int64) (bool, error) {
	newKey := c.BuildKey(key)
	conn, errConn := c.GetConnByKey("SET", newKey)
	if errConn != nil {
		return false, errConn
	}

	defer conn.Close(false)

	var res interface{}
	var err error
	res, err = conn.Do("EXPIREAT", newKey, timestamp)

	if err != nil {
		return false, err
	}

	num, ok := res.(int64)
	return ok && num == 1, nil
}

// args = [0:"key"]
func (c *Client) Do(cmd string, args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("The length of args has to be greater than 1")
	}

	key, ok := args[0].(string)
	if ok == false {
		return nil, fmt.Errorf("Invalid key string:%s", args[0])
	}

	cmd = strings.ToUpper(cmd)
	if util.SliceSearchString(allRedisCmd, cmd) == -1 {
		return nil, fmt.Errorf("Undefined command:%s", cmd)
	}

	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey(cmd, newKey)
	if err != nil {
		return nil, err
	}
	defer conn.Close(false)

	args[0] = newKey

	return conn.Do(cmd, args...)
}

func (c *Client) Expire(key string, expire time.Duration) (bool, error) {
	newKey := c.BuildKey(key)
	conn, errConn := c.GetConnByKey("EXPIRE", newKey)
	if errConn != nil {
		return false, errConn
	}

	defer conn.Close(false)

	var res interface{}
	var err error
	res, err = conn.Do("EXPIRE", newKey, expire/time.Second)

	if err != nil {
		return false, err
	}

	num, ok := res.(int64)
	return ok && num == 1, nil
}

func (c *Client) RPush(key string, values ...interface{}) (bool, error) {
	return c.listPush("RPUSH", key, values...)
}

func (c *Client) LPush(key string, values ...interface{}) (bool, error) {
	return c.listPush("LPUSH", key, values...)
}

func (c *Client) listPush(cmd, key string, values ...interface{}) (bool, error) {
	newKey := c.BuildKey(key)
	conn, errConn := c.GetConnByKey(cmd, newKey)
	if errConn != nil {
		return false, errConn
	}

	defer conn.Close(false)

	var res interface{}
	var err error
	args := c.mergeKey(newKey, values...)
	res, err = conn.Do(cmd, args...)

	if err != nil {
		return false, err
	}

	num, ok := res.(int64)

	return ok && num > 0, nil
}

func (c *Client) mergeKey(key string, values ...interface{}) []interface{} {
	args := make([]interface{}, 0, len(values)+1)
	args = append(args, key)
	args = append(args, values...)

	return args
}

func (c *Client) listPop(cmd, key string) (*value.Value, error) {
	newKey := c.BuildKey(key)
	conn, errConn := c.GetConnByKey(cmd, newKey)
	if errConn != nil {
		return nil, errConn
	}

	defer conn.Close(false)

	var res interface{}
	var err error

	res, err = conn.Do(cmd, newKey)

	if err != nil {
		return nil, err
	}

	return value.New(res), nil
}

func (c *Client) RPop(key string) (*value.Value, error) {
	return c.listPop("RPOP", key)
}

func (c *Client) LPop(key string) (*value.Value, error) {
	return c.listPop("LPOP", key)
}

func (c *Client) LLen(key string) (int64, error) {
	newKey := c.BuildKey(key)
	conn, errConn := c.GetConnByKey("LLEN", newKey)
	if errConn != nil {
		return 0, errConn
	}

	defer conn.Close(false)

	res, err := conn.Do("LLEN", newKey)
	if err != nil {
		return 0, err
	}

	num, _ := res.(int64)

	return num, nil
}

func (c *Client) HDel(key string, fields ...interface{}) (int64, error) {
	newKey := c.BuildKey(key)
	conn, errConn := c.GetConnByKey("HDEL", newKey)
	if errConn != nil {
		return 0, errConn
	}

	defer conn.Close(false)

	args := c.mergeKey(newKey, fields...)

	res, err := conn.Do("HDEL", args...)
	if err != nil {
		return 0, err
	}

	num, _ := res.(int64)

	return num, nil
}

func (c *Client) HExists(key string, field string) (bool, error) {
	newKey := c.BuildKey(key)
	conn, errConn := c.GetConnByKey("HEXISTS", newKey)
	if errConn != nil {
		return false, errConn
	}

	defer conn.Close(false)

	res, err := conn.Do("HEXISTS", newKey, field)
	if err != nil {
		return false, err
	}

	num, ok := res.(int64)

	return ok && num > 0, nil
}

func (c *Client) HSet(key string, fv ...interface{}) (bool, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("HSET", newKey)
	if err != nil {
		return false, err
	}
	defer conn.Close(false)

	pl := len(fv)
	if pl < 2 || pl%2 != 0 {
		return false, errors.New(errParamsNum)
	}

	args := c.mergeKey(newKey, fv...)

	v, err := conn.Do("HSET", args...)
	if err != nil {
		return false, err
	}

	num, ok := v.(int64)

	return ok && num >= 0, nil
}

func (c *Client) HGet(key, field string) (*value.Value, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("HGET", newKey)
	if err != nil {
		return nil, err
	}
	defer conn.Close(false)
	v, err := conn.Do("HGET", newKey, field)
	if err != nil {
		return nil, err
	}
	return value.New(v), nil
}

func (c *Client) HGetAll(key string) (map[string]*value.Value, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("HGETALL", newKey)
	if err != nil {
		return nil, err
	}
	defer conn.Close(false)
	iV, err := conn.Do("HGETALL", newKey)
	if err != nil {
		return nil, err
	}

	v, ok := iV.([]interface{})
	if !ok {
		return nil, nil
	}

	ret := make(map[string]*value.Value)
	k := 0
	vL := len(v) + 1
	for _, vv := range v {
		k++
		if k%2 == 0 {
			continue
		}

		var nextV interface{}
		if vL >= k {
			nextV = v[k]
		}

		ret[value.New(vv).String()] = value.New(nextV)
	}

	return ret, nil
}

func (c *Client) HMSet(key string, fv ...interface{}) (bool, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("HMSET", newKey)
	if err != nil {
		return false, err
	}
	defer conn.Close(false)

	pl := len(fv)
	if pl < 2 || pl%2 != 0 {
		return false, errors.New(errParamsNum)
	}

	args := c.mergeKey(newKey, fv...)

	v, err := conn.Do("HMSET", args...)
	if err != nil {
		return false, err
	}
	payload, ok := v.([]byte)
	if ok && bytes.Equal(payload, replyOK) {
		return true, nil
	}

	return false, nil
}

func (c *Client) HMGet(key string, fields ...interface{}) (map[string]*value.Value, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("HMGET", newKey)
	if err != nil {
		return nil, err
	}
	defer conn.Close(false)

	args := c.mergeKey(newKey, fields...)
	v, err := conn.Do("HMGET", args...)
	if err != nil {
		return nil, err
	}

	res := make(map[string]*value.Value)

	if tmpV, ok := v.([]interface{}); ok {
		for k, vv := range tmpV {
			res[fields[k].(string)] = value.New(vv)
		}
	}

	return res, nil
}

func (c *Client) HIncrBy(key, field string, delta int64) (int64, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("HINCRBY", newKey)
	if err != nil {
		return 0, err
	}
	defer conn.Close(false)

	v, err := conn.Do("HINCRBY", newKey, field, delta)
	if err != nil {
		return 0, err
	}

	num, _ := v.(int64)
	return num, nil
}

func (c *Client) zRange(cmd, key string, start, end int, opt ...string) ([]interface{}, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey(cmd, newKey)
	if err != nil {
		return nil, err
	}
	defer conn.Close(false)

	var iV interface{}
	if len(opt) > 0 {
		iV, err = conn.Do(cmd, newKey, start, end, opt[0])
	} else {
		iV, err = conn.Do(cmd, newKey, start, end)
	}

	if err != nil {
		return nil, err
	}

	v, ok := iV.([]interface{})
	if !ok {
		return nil, nil
	}

	return v, nil

}

func (c *Client) zRangeFormat(v []interface{}, err error) ([]*value.Value, error) {
	if v == nil || err != nil {
		return nil, err
	}
	ret := make([]*value.Value, 0, len(v))
	for _, vv := range v {
		ret = append(ret, value.New(vv))
	}

	return ret, err
}

func (c *Client) zRangeFormatWithScores(v []interface{}, err error) ([]*ZV, error) {
	if v == nil || err != nil {
		return nil, err
	}
	ret := make([]*ZV, 0, len(v))
	k := 0
	vL := len(v) + 1
	for _, vv := range v {
		k++
		if k%2 == 0 {
			continue
		}

		var score interface{}
		if vL >= k {
			score = v[k]
		}
		ret = append(ret, &ZV{value.New(score), value.New(vv)})
	}

	return ret, err
}

func (c *Client) ZRevRange(key string, start, end int) ([]*value.Value, error) {
	return c.zRangeFormat(c.zRange("ZREVRANGE", key, start, end))
}

func (c *Client) ZRange(key string, start, end int) ([]*value.Value, error) {
	return c.zRangeFormat(c.zRange("ZRANGE", key, start, end))
}

func (c *Client) ZRevRangeWithScores(key string, start, end int) ([]*ZV, error) {
	return c.zRangeFormatWithScores(c.zRange("ZREVRANGE", key, start, end, "WITHSCORES"))
}

func (c *Client) ZRangeWithScores(key string, start, end int) ([]*ZV, error) {
	return c.zRangeFormatWithScores(c.zRange("ZRANGE", key, start, end, "WITHSCORES"))
}

func (c *Client) zAdd(a []interface{}, n int, members ...*Z) (int64, error) {
	newKey := c.BuildKey(a[0].(string))
	a[0] = newKey
	conn, err := c.GetConnByKey("ZADD", newKey)
	if err != nil {
		return 0, err
	}
	defer conn.Close(false)

	for i, m := range members {
		a[n+2*i] = m.Score
		a[n+2*i+1] = m.Member
	}

	v, err := conn.Do("ZADD", a...)
	if err != nil {
		return 0, err
	}

	num, _ := v.(int64)
	return num, nil
}

func (c *Client) ZAdd(key string, members ...*Z) (int64, error) {
	const n = 1
	a := make([]interface{}, n+2*len(members))
	a[0] = key
	return c.zAdd(a, n, members...)
}

func (c *Client) ZAddOpt(key string, opts []string, members ...*Z) (int64, error) {
	n := 1 + len(opts)
	a := make([]interface{}, n+2*len(members))
	a[0] = key
	for k, v := range opts {
		a[k+1] = v
	}
	return c.zAdd(a, n, members...)
}

func (c *Client) ZCard(key string) (int64, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("ZCARD", newKey)
	if err != nil {
		return 0, err
	}
	defer conn.Close(false)

	v, err := conn.Do("ZCARD", newKey)
	if err != nil {
		return 0, err
	}

	num, _ := v.(int64)
	return num, nil
}

func (c *Client) ZRem(key string, members ...interface{}) (int64, error) {
	newKey := c.BuildKey(key)
	conn, err := c.GetConnByKey("ZREM", newKey)
	if err != nil {
		return 0, err
	}
	defer conn.Close(false)

	args := c.mergeKey(newKey, members...)

	v, err := conn.Do("ZREM", args...)
	if err != nil {
		return 0, err
	}

	num, _ := v.(int64)
	return num, nil
}
