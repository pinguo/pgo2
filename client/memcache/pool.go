package memcache

import (
    "errors"
    "fmt"
    "net"
    "strings"
    "sync"
    "time"

    "github.com/pinguo/pgo2/util"
)

type serverInfo struct {
    weight   int
    disabled bool
}

type connList struct {
    count int
    head  *Conn
    tail  *Conn
}

type Pool struct {
    lock      sync.RWMutex
    hashRing  *util.HashRing
    connLists map[string]*connList
    servers   map[string]*serverInfo

    prefix        string
    maxIdleConn   int
    maxIdleTime   time.Duration
    netTimeout    time.Duration
    probeInterval time.Duration
}

func (p *Pool) Init() {
    p.initServers()
    if p.probeInterval > 0 {
        go p.probeLoop()
    }
}

func (p *Pool) initServers() {
    if len(p.servers) == 0 {
        p.servers[defaultServer] = &serverInfo{weight: 1}
    }

    for addr, item := range p.servers {
        p.hashRing.AddNode(addr, item.weight)
    }

}

func (p *Pool) SetPrefix(prefix string) {
    p.prefix = prefix
}

func (p *Pool) SetServers(v []interface{}) {
    for _, vv := range v {
        addr := vv.(string)
        if pos := strings.Index(addr, "://"); pos != -1 {
            addr = addr[pos+3:]
        }

        item := p.servers[addr]
        if item == nil {
            item = &serverInfo{}
            p.servers[addr] = item
        }

        item.weight += 1
    }
}

func (p *Pool) GetServers() (servers []string) {
    for server := range p.servers {
        servers = append(servers, server)
    }
    return
}

func (p *Pool) SetMaxIdleConn(v int) {
    p.maxIdleConn = v
}

func (p *Pool) SetMaxIdleTime(v string) error {
    if maxIdleTime, e := time.ParseDuration(v); e != nil {
        return fmt.Errorf(errSetProp, "maxIdleTime", e)
    } else {
        p.maxIdleTime = maxIdleTime
    }

    return nil
}

func (p *Pool) SetNetTimeout(v string) error {
    if netTimeout, e := time.ParseDuration(v); e != nil {
        return fmt.Errorf(errSetProp, "netTimeout", e)
    } else {
        p.netTimeout = netTimeout
    }

    return nil
}

func (p *Pool) SetProbeInterval(v string) error {
    if probeInterval, e := time.ParseDuration(v); e != nil {
        return fmt.Errorf(errSetProp, "probeInterval", e)
    } else {
        p.probeInterval = probeInterval
        if p.probeInterval > maxProbeInterval {
            p.probeInterval = maxProbeInterval
        } else if p.probeInterval < minProbeInterval {
            p.probeInterval = minProbeInterval
        }
    }

    return nil
}

func (p *Pool) BuildKey(key string) string {
    return p.prefix + key
}

func (p *Pool) AddrNewKeys(v interface{}) (map[string][]string, map[string]string, error) {
    addrKeys, newKeys := make(map[string][]string), make(map[string]string)
    switch vv := v.(type) {
    case []string:
        for _, key := range vv {
            newKey := p.BuildKey(key)
            addr := p.GetAddrByKey(newKey)
            newKeys[newKey] = key
            addrKeys[addr] = append(addrKeys[addr], newKey)
        }
    case map[string]interface{}:
        for key := range vv {
            newKey := p.BuildKey(key)
            addr := p.GetAddrByKey(newKey)
            newKeys[newKey] = key
            addrKeys[addr] = append(addrKeys[addr], newKey)
        }
    default:
        return nil, nil, errors.New(errBase + "addr new keys invalid")
    }
    return addrKeys, newKeys, nil
}

func (p *Pool) RunAddrFunc(addr string, keys []string, wg *sync.WaitGroup, f func(*Conn, []string)) error {
    defer func() {
        recover() // ignore panic
        wg.Done() // notify done
    }()

    conn, err := p.GetConnByAddr(addr)
    if err != nil {
        return err
    }

    defer conn.Close(false)

    f(conn, keys)

    return nil
}

func (p *Pool) GetConnByKey(key string) (*Conn, error) {
    if addr := p.GetAddrByKey(key); len(addr) == 0 {
        return nil, errors.New(errNoServer)
    } else {
        return p.GetConnByAddr(addr)
    }
}

func (p *Pool) GetConnByAddr(addr string) (conn *Conn, err error) {
    conn = p.getFreeConn(addr)
    if conn == nil || !p.checkConn(conn) {
        conn, err = p.dial(addr)
        if err != nil {
            return nil, err
        }

    }

    conn.ExtendDeadLine()
    return conn, err
}

func (p *Pool) GetAddrByKey(key string) string {
    p.lock.RLock()
    defer p.lock.RUnlock()
    return p.hashRing.GetNode(key)
}

func (p *Pool) getFreeConn(addr string) *Conn {
    p.lock.Lock()
    defer p.lock.Unlock()

    list := p.connLists[addr]
    if list == nil || list.count == 0 {
        return nil
    }

    conn := list.head
    if list.count--; list.count == 0 {
        list.head, list.tail = nil, nil
    } else {
        list.head, conn.next.prev = conn.next, nil
    }

    conn.next = nil
    return conn
}

func (p *Pool) putFreeConn(conn *Conn) bool {
    p.lock.Lock()
    defer p.lock.Unlock()

    list := p.connLists[conn.addr]
    if list == nil {
        list = new(connList)
        p.connLists[conn.addr] = list
    }

    if list.count >= p.maxIdleConn {
        return false
    }

    if list.count == 0 {
        list.head, list.tail = conn, conn
        conn.prev, conn.next = nil, nil
    } else {
        conn.prev, conn.next = list.tail, nil
        conn.prev.next, list.tail = conn, conn
    }

    list.count++
    return true
}

func (p *Pool) checkConn(conn *Conn) bool {
    defer func() {
        // if panic, return value is default(false)
        if v := recover(); v != nil {
            conn.Close(true)
        }
    }()

    if active,err:=conn.CheckActive();!active || err!=nil {
        conn.Close(true)
        return false
    }
    return true
}

func (p *Pool) dial(addr string) (*Conn, error) {
    if nc, e := net.DialTimeout(p.parseNetwork(addr), addr, p.netTimeout); e != nil {
        return nil, e
    } else {
        return newConn(addr, nc, p), nil
    }
}

func (p *Pool) parseNetwork(addr string) string {
    if pos := strings.IndexByte(addr, '/'); pos != -1 {
        return "unix"
    } else {
        return "tcp"
    }
}

func (p *Pool) probeServer(addr string) {
    nc, e := net.DialTimeout(p.parseNetwork(addr), addr, p.netTimeout)
    if e != nil && !p.servers[addr].disabled {
        func(){
            p.lock.Lock()
            defer p.lock.Unlock()
            p.servers[addr].disabled = true
            p.hashRing.DelNode(addr)
        }()

    } else if e == nil && p.servers[addr].disabled {
        func(){
            p.lock.Lock()
            defer p.lock.Unlock()
            p.servers[addr].disabled = false
            p.hashRing.AddNode(addr, p.servers[addr].weight)
        }()
    }

    if e == nil {
        nc.Close()
    }
}

func (p *Pool) probeLoop() {
    defer func() {
        if err := recover();err != nil{
            fmt.Println("memCache pool.probeLoop err:" + util.ToString(err))
        }

        p.probeLoop()
    }()

    for {
        <-time.After(p.probeInterval)
        for addr := range p.servers {
            p.probeServer(addr)
        }
    }
}
