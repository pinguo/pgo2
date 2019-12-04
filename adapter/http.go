package adapter

import (
    "net/http"
    "sync"
    "time"

    "github.com/pinguo/pgo2"
    "github.com/pinguo/pgo2/client/phttp"
    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/util"
)

func init() {
    container := pgo2.App().Container()
    container.Bind(&Http{})
}

// Http of Http Client, add context support.
// usage: http := this.GetObj(adapter.NewHttp()).(adapter.IHttp)/(*adapter.Http)
func NewHttp(componentId ...string) *Http {
    id := DefaultHttpId
    if len(componentId) > 0 {
        id = componentId[0]
    }

    h := &Http{}

    client := pgo2.App().Component(id, phttp.New)
    h.client = client.(*phttp.Client)
    h.panicRecover = true

    return h
}

// Http of Http Client from pool, add context support.
// usage: http := this.GetObjPool(adapter.NewHttpPool).(adapter.IHttp)/(*adapter.Http)
func NewHttpPool(ctr *pgo2.Context, componentId ...interface{}) iface.IObject {
    id := DefaultHttpId
    if len(componentId) > 0 {
        id = componentId[0].(string)
    }

    h := pgo2.App().GetObjPool(HttpClass, ctr).(*Http)

    client := pgo2.App().Component(id, phttp.New)
    h.client = client.(*phttp.Client)
    h.panicRecover = true

    return h
}

type Http struct {
    pgo2.Object
    client       *phttp.Client
    panicRecover bool
}

func (h *Http) SetPanicRecover(v bool) {
    h.panicRecover = v
}

func (h *Http) GetClient() *phttp.Client {
    return h.client
}

func (h *Http) handlePanic() {
    if h.panicRecover {
        if v := recover(); v != nil {
            h.Context().Error(util.ToString(v))
        }
    }
}

func (h *Http) parseErr(err error) {
    if err != nil {
        panic(err)
    }
}

// Get perform a get request
func (h *Http) Get(addr string, data interface{}, option ...*phttp.Option) *http.Response {
    profile := baseUrl(addr)
    h.Context().ProfileStart(profile)
    defer h.Context().ProfileStop(profile)
    defer h.handlePanic()

    res, err := h.client.Get(addr, data, option...)
    h.parseErr(err)

    return res
}

// Post perform a post request
func (h *Http) Post(addr string, data interface{}, option ...*phttp.Option) *http.Response {
    profile := baseUrl(addr)
    h.Context().ProfileStart(profile)
    defer h.Context().ProfileStop(profile)
    defer h.handlePanic()

    res, err := h.client.Post(addr, data, option...)
    h.parseErr(err)

    return res
}

// Do perform a single request
func (h *Http) Do(req *http.Request, option ...*phttp.Option) *http.Response {
    profile := baseUrl(req.URL.String())
    h.Context().ProfileStart(profile)
    defer h.Context().ProfileStop(profile)
    defer h.handlePanic()

    res, err := h.client.Do(req, option...)
    h.parseErr(err)

    return res
}

// DoMulti perform multi requests concurrently
func (h *Http) DoMulti(requests []*http.Request, option ...*phttp.Option) []*http.Response {
    if num := len(option); num != 0 && num != len(requests) {
        panic("http multi request invalid num of options")
    }

    profile := "Http.DoMulti"
    h.Context().ProfileStart(profile)
    defer h.Context().ProfileStop(profile)
    defer h.handlePanic()

    lock, wg := new(sync.Mutex), new(sync.WaitGroup)
    responses := make([]*http.Response, len(requests))

    fn := func(k int) {
        start, profile := time.Now(), baseUrl(requests[k].URL.String())
        var res *http.Response
        var err error

        defer func() {
            if v := recover(); v != nil {
                h.Context().Error(util.ToString(v))
            }

            lock.Lock()
            h.Context().ProfileAdd(profile, time.Since(start)/1e6)
            responses[k] = res
            lock.Unlock()
            wg.Done()
        }()

        if len(option) > 0 && option[k] != nil {
            res, err = h.client.Do(requests[k], option[k])
        } else {
            res, err = h.client.Do(requests[k])
        }
        h.parseErr(err)
    }

    wg.Add(len(requests))
    for k := range requests {
        go fn(k)
    }

    wg.Wait()

    return responses
}
