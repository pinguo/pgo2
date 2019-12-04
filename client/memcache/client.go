package memcache

import (
    "sync"
    "sync/atomic"
    "time"

    "github.com/pinguo/pgo2/core"
    "github.com/pinguo/pgo2/util"
    "github.com/pinguo/pgo2/value"
)

// Memcache Client component, configuration:
// components:
//      memcache:
//          class: "@pgo/Client/Memcache/Client"
//          prefix: "pgo_"
//          maxIdleConn: 10
//          maxIdleTime: "60s"
//          netTimeout: "1s"
//          probInterval: "0s"
//          servers:
//              - "127.0.0.1:11211"
//              - "127.0.0.1:11212"
func New(config map[string]interface{}) (interface{}, error) {
    c := &Client{}
    c.hashRing = util.NewHashRing()
    c.connLists = make(map[string]*connList)
    c.servers = make(map[string]*serverInfo)

    c.prefix = defaultPrefix
    c.maxIdleConn = defaultIdleConn
    c.maxIdleTime = defaultIdleTime
    c.netTimeout = defaultTimeout
    c.probeInterval = defaultProbe

    err := core.ClientConfigure(c, config)
    if err != nil {
        return nil, err
    }

    c.Init()

    return c, nil
}

type Client struct {
    Pool
}

func (c *Client) Get(key string) (*value.Value, error) {
    if item, err := c.Retrieve(CmdGet, key); item != nil {
        if err != nil {
            return nil, err
        }

        return value.New(item.Data), nil
    }
    return value.New(nil), nil
}

func (c *Client) MGet(keys []string) (map[string]*value.Value, error) {
    result := make(map[string]*value.Value)
    for _, key := range keys {
        result[key] = value.New(nil)
    }

    if items, err := c.MultiRetrieve(CmdGet, keys); len(items) > 0 {
        if err != nil {
            return nil, err
        }
        for _, item := range items {
            result[item.Key] = value.New(item.Data)
        }
    }
    return result, nil
}

func (c *Client) Set(key string, v interface{}, expire ...time.Duration) (bool, error) {

    return c.Store(CmdSet, &Item{Key: key, Data: value.Encode(v)}, expire...)
}

func (c *Client) MSet(items map[string]interface{}, expire ...time.Duration) (bool, error) {
    newItems := make([]*Item, 0, len(items))
    for key, v := range items {
        newItems = append(newItems, &Item{Key: key, Data: value.Encode(v)})
    }
    return c.MultiStore(CmdSet, newItems, expire...)
}

func (c *Client) Add(key string, v interface{}, expire ...time.Duration) (bool, error) {
    return c.Store(CmdAdd, &Item{Key: key, Data: value.Encode(v)}, expire...)
}

func (c *Client) MAdd(items map[string]interface{}, expire ...time.Duration) (bool, error) {
    newItems := make([]*Item, 0, len(items))
    for key, v := range items {
        newItems = append(newItems, &Item{Key: key, Data: value.Encode(v)})
    }
    return c.MultiStore(CmdAdd, newItems, expire...)
}

func (c *Client) Del(key string) (bool, error) {
    newKey := c.BuildKey(key)
    conn, err := c.GetConnByKey(newKey)
    if err != nil {
        return false, err
    }
    defer conn.Close(false)

    return conn.Delete(newKey)
}

func (c *Client) MDel(keys []string) (bool, error) {
    addrKeys, _, errKey := c.AddrNewKeys(keys)
    if errKey != nil {
        return false, errKey
    }
    wg, success := new(sync.WaitGroup), uint32(0)

    var err error
    wg.Add(len(addrKeys))
    for addr, keys := range addrKeys {
        go c.RunAddrFunc(addr, keys, wg, func(conn *Conn, keys []string) {
            for _, key := range keys {
                // extend deadline for every operation
                conn.ExtendDeadLine()
                delOk, err := conn.Delete(key)
                if err != nil{
                    return
                }

                if delOk {
                    atomic.AddUint32(&success, 1)
                }
            }
        })
    }

    wg.Wait()
    return success == uint32(len(keys)), err
}

func (c *Client) Exists(key string) (bool, error) {
    v, err := c.Get(key)
    if err != nil {
        return false, err
    }

    if v == nil {
        return false, nil
    }

    return v.Valid(), nil
}

func (c *Client) Incr(key string, delta int) (int, error) {
    newKey := c.BuildKey(key)
    conn, err := c.GetConnByKey(newKey)
    if err != nil {
        return 0, err
    }
    defer conn.Close(false)

    return conn.Increment(newKey, delta)
}

func (c *Client) Retrieve(cmd, key string) (*Item, error) {
    newKey := c.BuildKey(key)
    conn, err := c.GetConnByKey(newKey)
    if err != nil {
        return nil, err
    }
    defer conn.Close(false)


    if items, err := conn.Retrieve(cmd, newKey);err!=nil {
        return nil,err
    }else if len(items) == 1{
        return items[0], nil
    }

    return nil, nil
}

func (c *Client) MultiRetrieve(cmd string, keys []string) ([]*Item, error) {
    result := make([]*Item, 0, len(keys))
    addrKeys, newKeys, errKey := c.AddrNewKeys(keys)
    if errKey != nil {
        return nil, errKey
    }
    lock, wg := new(sync.Mutex), new(sync.WaitGroup)

    var err error
    wg.Add(len(addrKeys))
    for addr, keys := range addrKeys {
        go c.RunAddrFunc(addr, keys, wg, func(conn *Conn, keys []string) {
            items, err := conn.Retrieve(cmd, keys...)
            if err !=nil{
                return
            }

            if len(items) > 0 {
                lock.Lock()
                defer lock.Unlock()
                for _, item := range items {
                    item.Key = newKeys[item.Key]
                    result = append(result, item)
                }
            }
        })
    }

    wg.Wait()
    return result, err
}

func (c *Client) Store(cmd string, item *Item, expire ...time.Duration) (bool, error) {
    item.Key = c.BuildKey(item.Key)
    conn, err := c.GetConnByKey(item.Key)
    if err != nil {
        return false, err
    }
    defer conn.Close(false)

    expire = append(expire, defaultExpire)
    return conn.Store(cmd, item, int(expire[0]/time.Second))
}

func (c *Client) MultiStore(cmd string, items []*Item, expire ...time.Duration) (bool, error) {
    expire = append(expire, defaultExpire)
    addrItems := make(map[string][]*Item)
    wg, success := new(sync.WaitGroup), uint32(0)

    for _, item := range items {
        item.Key = c.BuildKey(item.Key)
        addr := c.GetAddrByKey(item.Key)
        addrItems[addr] = append(addrItems[addr], item)
    }

    var err error
    wg.Add(len(addrItems))
    for addr := range addrItems {
        go c.RunAddrFunc(addr, nil, wg, func(conn *Conn, keys []string) {
            for _, item := range addrItems[addr] {
                conn.ExtendDeadLine() // extend deadline for every store
                ok, err := conn.Store(cmd, item, int(expire[0]/time.Second))
                if err != nil{
                    return
                }
                if ok {
                    atomic.AddUint32(&success, 1)
                }
            }
        })
    }

    wg.Wait()
    return success == uint32(len(items)), err
}
