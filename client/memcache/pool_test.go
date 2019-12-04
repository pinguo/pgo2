package memcache

import (
    "testing"

    "github.com/pinguo/pgo2/util"
)

func initPool() *Pool{
    p := &Pool{}
    p.hashRing = util.NewHashRing()
    p.connLists = make(map[string]*connList)
    p.servers = make(map[string]*serverInfo)
    p.initServers()
    p.prefix = "pgp2_"
    return p
}

func TestPool_initServers(t *testing.T) {
    p := initPool()
    if p.hashRing.GetNode("test") != defaultServer {
        t.FailNow()
    }

}

func TestPool_AddrNewKeys(t *testing.T) {
    p := initPool()

    testKey := "testKey"

    t.Run("[]string", func(t *testing.T) {
        keys,newKeys,err:=p.AddrNewKeys([]string{testKey})
        if keys[defaultServer][0]!=p.prefix + testKey{
            t.Fatal(`keys[defaultServer][0]!=p.prefix + testKey`)
        }

        if v,has:=newKeys[p.prefix + testKey];has==false || v!=testKey{
            t.Fatal( `v,has:=newKeys[p.prefix + testKey];has==false || v!=testKey`)
        }

        if err !=nil{
            t.Fatal(`err !=nil`)
        }
    })

    t.Run("map[string]interface{}", func(t *testing.T) {
        keys,newKeys,err:=p.AddrNewKeys(map[string]interface{}{testKey:""})
        if keys[defaultServer][0]!=p.prefix + testKey{
            t.Fatal(`keys[defaultServer][0]!=p.prefix + testKey`)
        }

        if v,has:=newKeys[p.prefix + testKey];has==false || v!=testKey{
            t.Fatal( `v,has:=newKeys[p.prefix + testKey];has==false || v!=testKey`)
        }

        if err !=nil{
            t.Fatal(`err !=nil`)
        }
    })

    t.Run("invalid", func(t *testing.T) {
        _,_,err:=p.AddrNewKeys(map[string]string{testKey:""})
        if err == nil{
            t.FailNow()
        }
    })

}

func TestPool_BuildKey(t *testing.T) {
    p := initPool()
    testKey := "testKey"
    if p.BuildKey(testKey) != p.prefix + testKey{
        t.FailNow()
    }
}

func TestPool_GetAddrByKey(t *testing.T) {
    p := initPool()

    testKey := "testKey"

    if p.GetAddrByKey(testKey) != defaultServer {
        t.FailNow()
    }
}

func TestPool_GetConnByAddr(t *testing.T) {

}

func TestPool_GetConnByKey(t *testing.T) {

}

func TestPool_GetServers(t *testing.T) {

}

func TestPool_RunAddrFunc(t *testing.T) {

}

func TestPool_SetMaxIdleConn(t *testing.T) {

}

func TestPool_SetMaxIdleTime(t *testing.T) {

}

func TestPool_SetNetTimeout(t *testing.T) {

}

func TestPool_SetPrefix(t *testing.T) {

}

func TestPool_SetProbeInterval(t *testing.T) {

}

func TestPool_SetServers(t *testing.T) {

}
