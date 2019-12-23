package adapter

import (
	"github.com/pinguo/pgo2"
	"github.com/pinguo/pgo2/client/maxmind"
	"github.com/pinguo/pgo2/iface"
)

var MaxMindClass string
func init() {
	MaxMindClass = pgo2.App().Container().Bind(&MaxMind{})
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
// usage: mmd := this.GetObjPool(adapter.MaxMindClass, adapter.NewMaxMindPool).(adapter.IMaxMind)/(*adapter.MaxMind)
func NewMaxMindPool(iObj iface.IObject, componentId ...interface{}) iface.IObject {
	id := DefaultMaxMindId
	if len(componentId) > 0 {
		id = componentId[0].(string)
	}

	m := iObj.(*MaxMind)

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
	profile := "GeoByIp:" + ip
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	geo, err := m.client.GeoByIp(ip, args...)

	if err != nil {
		panic("GeoByIp err:" + err.Error())
	}

	return geo
}
