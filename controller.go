package pgo2

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/pinguo/pgo2/iface"
	"github.com/pinguo/pgo2/perror"
	"github.com/pinguo/pgo2/render"
	"github.com/pinguo/pgo2/util"
)

func init() {
	container := App().Container()
	App().Router().SetErrorController(container.Bind(&Controller{}))
}

// Controller the base class of web and cmd controller
type Controller struct {
	Object
}

// GetBindInfo get action map as extra binding info
func (c *Controller) GetBindInfo(v interface{}) interface{} {
	if _, ok := v.(iface.IController); !ok {
		panic("param require a controller")
	}

	rt := reflect.ValueOf(v).Type()
	num := rt.NumMethod()
	actions := make(map[string]int)

	for i := 0; i < num; i++ {
		name := rt.Method(i).Name
		if _, ok := restFulActions[name]; ok {
			actions[name] = i
			continue
		}

		if len(name) > ActionLength && name[:ActionLength] == ActionPrefix {
			actions[name[ActionLength:]] = i
		}
	}

	return actions
}

// BeforeAction before action hook
func (c *Controller) BeforeAction(action string) {
}

// AfterAction after action hook
func (c *Controller) AfterAction(action string) {
}

// HandlePanic process unhandled action panic
func (c *Controller) HandlePanic(v interface{}, debug bool) {
	status := http.StatusInternalServerError

	switch e := v.(type) {
	case *perror.Error:
		status := e.Status()
		defer func() {
			if err := recover(); err != nil {
				c.Json(EmptyObject, status, e.Message())
				c.Context().Error("%s, trace[%s]", util.ToString(err), util.PanicTrace(TraceMaxDepth, false, debug))
			}
		}()

		App().Router().ErrorController(c.Context(), status).(iface.IErrorController).Error(status, e.Message())
	default:
		defer func() {
			if err := recover(); err != nil {
				c.Json(EmptyObject, status)
				c.Context().Error("%s, trace[%s]", util.ToString(err), util.PanicTrace(TraceMaxDepth, false, debug))
			}
		}()

		App().Router().ErrorController(c.Context(), status).(iface.IErrorController).Error(status, "")
	}

	c.Context().Error("%s, trace[%s]", util.ToString(v), util.PanicTrace(TraceMaxDepth, false, debug))
}

// Redirect output redirect response
func (c *Controller) Redirect(location string, permanent bool) {
	ctx := c.Context()
	ctx.SetHeader("Location", location)
	if permanent {
		ctx.End(http.StatusMovedPermanently, nil)
	} else {
		ctx.End(http.StatusFound, nil)
	}
}

// Json output json response
func (c *Controller) Json(data interface{}, status int, msg ...string) {
	ctx := c.Context()
	message := App().Status().Text(status, ctx.Header("Accept-Language", ""), msg...)
	r := render.NewJson(map[string]interface{}{
		"status":  status,
		"message": message,
		"data":    data,
	})

	ctx.PushLog("status", status)
	ctx.SetHeader("Content-Type", r.ContentType())
	httpStatus := r.HttpCode()
	if ctx.Status() > 0 && ctx.Status() != httpStatus {
		httpStatus = ctx.Status()
	}
	ctx.End(httpStatus, r.Content())
}

// Jsonp output jsonp response
func (c *Controller) Jsonp(callback string, data interface{}, status int, msg ...string) {
	ctx := c.Context()
	message := App().Status().Text(status, ctx.Header("Accept-Language", ""), msg...)
	r := render.NewJsonp(callback, map[string]interface{}{
		"status":  status,
		"message": message,
		"data":    data,
	})

	ctx.PushLog("status", status)
	ctx.SetHeader("Content-Type", r.ContentType())
	httpStatus := r.HttpCode()
	if ctx.Status() > 0 && ctx.Status() != httpStatus {
		httpStatus = ctx.Status()
	}

	ctx.End(httpStatus, r.Content())
}

// Data output data response
func (c *Controller) Data(data []byte) {
	ctx := c.Context()
	r := render.NewData(data)
	httpStatus := r.HttpCode()
	if ctx.Status() > 0 && ctx.Status() != httpStatus {
		httpStatus = ctx.Status()
	}

	ctx.PushLog("status", httpStatus)
	ctx.SetHeader("Content-Type", r.ContentType())
	ctx.End(httpStatus, r.Content())
}

// Xml output xml response
func (c *Controller) Xml(data interface{}, statuses ...int) {
	status := http.StatusOK
	if len(statuses) > 0 {
		status = statuses[0]
	}

	ctx := c.Context()
	r := render.NewXml(data)

	ctx.PushLog("status", status)
	ctx.SetHeader("Content-Type", r.ContentType())
	httpStatus := r.HttpCode()
	if ctx.Status() > 0 && ctx.Status() != httpStatus {
		httpStatus = ctx.Status()
	}

	ctx.End(httpStatus, r.Content())
}

// ProtoBuf output proto buf response
func (c *Controller) ProtoBuf(data interface{}, statuses ...int) {
	status := http.StatusOK
	if len(statuses) > 0 {
		status = statuses[0]
	}

	ctx := c.Context()
	r := render.NewProtoBuf(data)

	ctx.PushLog("status", status)
	ctx.SetHeader("Content-Type", r.ContentType())
	httpStatus := r.HttpCode()
	if ctx.Status() > 0 && ctx.Status() != httpStatus {
		httpStatus = ctx.Status()
	}

	ctx.End(httpStatus, r.Content())
}

// Render Custom renderer
func (c *Controller) Render(r render.Render, statuses ...int) {
	status := http.StatusOK
	if len(statuses) > 0 {
		status = statuses[0]
	}

	ctx := c.Context()
	ctx.PushLog("status", status)
	ctx.SetHeader("Content-Type", r.ContentType())
	httpStatus := r.HttpCode()
	if ctx.Status() > 0 && ctx.Status() != httpStatus {
		httpStatus = ctx.Status()
	}

	ctx.End(httpStatus, r.Content())
}

// View output rendered view
func (c *Controller) View(view string, data interface{}, contentTypes ...string) {
	ctx := c.Context()
	contentType := "text/html; charset=utf-8"
	if len(contentTypes) > 0 {
		contentType = contentTypes[0]
	}
	httpStatus := http.StatusOK
	ctx.PushLog("status", httpStatus)
	ctx.SetHeader("Content-Type", contentType)

	if ctx.Status() > 0 && ctx.Status() != httpStatus {
		httpStatus = ctx.Status()
	}
	ctx.End(httpStatus, App().View().Render(view, data))
}

// Error
func (c *Controller) Error(status int, message string) {
	c.Json(EmptyObject, status, message)
}
