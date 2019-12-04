package adapter

import (
    "github.com/pinguo/pgo2"
    "github.com/pinguo/pgo2/client/maxmind"
    "github.com/pinguo/pgo2/iface"
)

func init() {
    container := pgo2.App().Container()
    container.Bind(&MaxMind{})
}

// NewMaxMind of MaxMind Client, add context support.
// usage: mmd := this.GetObj(adapter.NewMaxMind()).(adapter.IMaxMind)/(*adapter.MaxMind)
func NewMaxMind(componentId ...string) *MaxMind {
    id := DefaultMaxMindId
    if len(componentId) > 0 {
        id = componentId[0]
    }

    m := &MaxMind{}

    m.client = pgo2.App().Component(id, maxmind.New).(*maxmind.Client)

    return m
}

// NewMaxMindPool of MaxMind Client from pool, add context support.
// usage: mmd := this.GetObjPool(adapter.NewMaxMindPool).(adapter.IMaxMind)/(*adapter.MaxMind)
func NewMaxMindPool(ctr iface.IContext, componentId ...interface{}) iface.IObject {
    id := DefaultMemoryId
    if len(componentId) > 0 {
        id = componentId[0].(string)
    }

    m := pgo2.App().GetObjPool(MaxMindClass, ctr).(*MaxMind)

    m.client = pgo2.App().Component(id, maxmind.New).(*maxmind.Client)

    return m
}

type MaxMind struct {
    pgo2.Object
    client *maxmind.Client
}

func (m *MaxMind) GetClient() *maxmind.Client {
    return m.client
}

// get geo info by ip, optional args:
// db int: preferred max mind db
// lang string: preferred i18n language
func (m *MaxMind) GeoByIp(ip string, args ...interface{}) *maxmind.Geo {
    profile :="GeoByIp:" + ip
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)
    geo, err := m.client.GeoByIp(ip, args...)

    if err != nil {
        panic("GeoByIp err:" + err.Error())
    }

    return geo
}
