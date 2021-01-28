package adapter

import (
    "time"

    "github.com/pinguo/pgo2"
    "github.com/pinguo/pgo2/client/es"
    "github.com/pinguo/pgo2/util"
)

var EsClass string
func init() {
    container := pgo2.App().Container()

    EsClass = container.Bind(&Es{})
}

// NewEs of ElasticSearch Client, add context support.
// usage: es := this.GetObj(adapter.NewEs()).(adapter.IEs)/(*adapter.Es)
func NewEs(componentId ...string) *Es {
    id := defaultEsId
    if len(componentId) > 0 {
        id = componentId[0]
    }
    e := &Es{}
    e.client = pgo2.App().Component(id, es.New).(*es.Client)

    return e
}

// Adapter of ElasticSearch Client, add context support.
// Adapter of Http Client, add context support.
// usage: http := this.GetObjBox(adapter.EsClass).(*adapter.Es)
type Es struct {
    pgo2.Object
    client       *es.Client
    panicRecover bool
}

func (e *Es) Prepare(componentId ...string) {
    id := defaultEsId
    if len(componentId) > 0 {
        id = componentId[0]
    }
    e.client = pgo2.App().Component(id, es.New).(*es.Client)
}

func (e *Es) GetClient() *es.Client {
    return e.client
}

func (e *Es) handlePanic() {
    if e.panicRecover {
        if v := recover(); v != nil {
            e.Context().Error(util.ToString(v))
        }
    }
}

// method: POST GET PUT DELETE
// uri
// body
// timeout 超时时间
func (e *Es) Single(method, uri string, body []byte, timeout time.Duration) ([]byte,error) {
    profile := uri
    e.Context().ProfileStart(profile)
    defer e.Context().ProfileStop(profile)

    return e.client.Single(method,uri,body,timeout)
}

// 批量增加命令  异步执行
// action :  index，create，delete, update
// head : {“_ index”：“test”，“_ type”：“_ doc”，“_ id”：“1”}
// body : {"filed1":"value1"}
func (e *Es) Batch(action, head, body string) error{
    profile := action
    e.Context().ProfileStart(profile)
    defer e.Context().ProfileStop(profile)

    return e.client.Batch(action, head, body)
}


