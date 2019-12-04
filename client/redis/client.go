package redis

import (
    "bytes"
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
//          prefix: "pgo_"
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
    num, ok := ret.(int)
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

            if num, ok := ret.(int); ok && num > 0 {
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
    num, ok := ret.(int)
    return ok && num == 1, err
}

func (c *Client) Incr(key string, delta int) (int, error) {
    newKey := c.BuildKey(key)
    conn, err := c.GetConnByKey("INCRBY", newKey)
    if err != nil {
        return 0, err
    }

    defer conn.Close(false)
    ret, err := conn.Do("INCRBY", newKey, delta)
    num, _ := ret.(int)
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

// args = [0:"key"]
func (c *Client) Do(cmd string, args ...interface{}) (interface{}, error) {
    if len(args) == 0 {
        panic("The length of args has to be greater than 1")
    }

    key, ok := args[0].(string)
    if ok == false {
        panic("Invalid key string:" + util.ToString(args[0]))
    }

    cmd = strings.ToUpper(cmd)
    if util.SliceSearchString(allRedisCmd, cmd) == -1 {
        panic("Undefined command:" + cmd)
    }

    newKey := c.BuildKey(key)
    conn, err := c.GetConnByKey(cmd, newKey)
    if err != nil {
        return nil, err
    }
    defer conn.Close(false)

    args[0] = newKey

    return conn.Do(cmd, args ...)
}
