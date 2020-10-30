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

var HttpClass string

func init() {
	container := pgo2.App().Container()
	HttpClass = container.Bind(&Http{})
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
// usage: http := this.GetObjPool(adapter.HttpClass, adapter.NewHttpPool).(adapter.IHttp)/(*adapter.Http)
// It is recommended to use : http := this.GetObjBox(adapter.HttpClass).(adapter.IHttp)/(*adapter.Http)
func NewHttpPool(iObj iface.IObject, componentId ...interface{}) iface.IObject {

	return iObj
}

type Http struct {
	pgo2.Object
	client       *phttp.Client
	panicRecover bool
}

// GetObjPool, GetObjBox fetch is performed automatically
func (h *Http) Prepare(componentId ...interface{}) {
	id := DefaultHttpId
	if len(componentId) > 0 {
		id = componentId[0].(string)
	}

	h.client = pgo2.App().Component(id, phttp.New).(*phttp.Client)
	h.panicRecover = true
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

	if option == nil {
		option = make([]*phttp.Option, 1)
		option[0] = &phttp.Option{}
	}

	if _, has := option[0].Header["X-Log-Id"]; has == false {
		option[0].SetHeader("X-Log-Id", h.Context().LogId())
	}

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

	if option == nil {
		option = make([]*phttp.Option, 1)
		option[0] = &phttp.Option{}
	}

	if _, has := option[0].Header["X-Log-Id"]; has == false {
		option[0].SetHeader("X-Log-Id", h.Context().LogId())
	}

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

	if option == nil {
		option = make([]*phttp.Option, 1)
		option[0] = &phttp.Option{}
	}

	if _, has := option[0].Header["X-Log-Id"]; has == false {
		option[0].SetHeader("X-Log-Id", h.Context().LogId())
	}

	res, err := h.client.Do(req, option...)
	h.parseErr(err)

	return res
}

// DoMulti perform multi requests concurrently
func (h *Http) DoMulti(requests []*http.Request, option ...*phttp.Option) []*http.Response {
	if num := len(option); num != 0 && num != len(requests) {
		panic("http multi request invalid num of options")
	}

	baseOption := &phttp.Option{}
	if len(option) > 0 {
		baseOption = option[0]
	}

	if _, has := baseOption.Header["X-Log-Id"]; has == false {
		baseOption.SetHeader("X-Log-Id", h.Context().LogId())
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
			if _, has := option[k].Header["X-Log-Id"]; has == false {
				option[k].SetHeader("X-Log-Id", h.Context().LogId())
			}

			res, err = h.client.Do(requests[k], option[k])
		} else {

			res, err = h.client.Do(requests[k], baseOption)
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
