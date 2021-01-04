package pgo2

import (
	"net/http"
	"reflect"
	"time"

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
		} else {
			//runeArr :=  []rune( name[0:1])
			//if unicode.IsUpper(runeArr[0]) {
			//	actions[name] = i
			//}
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
	pErrorType := ""

	recoverErr := func(message string) {
		if err := recover(); err != nil {
			c.Json(EmptyObject, status, message)
			c.Context().Error("%s, trace[%s]", util.ToString(err), util.PanicTrace(TraceMaxDepth, false, debug))
		}
	}

	switch e := v.(type) {
	case *perror.Error:
		status = e.Status()
		pErrorType = e.ErrType()

		defer recoverErr(e.Message())

		App().Router().ErrorController(c.Context()).(iface.IErrorController).Error(status, e.Message())
	default:
		defer recoverErr("")

		App().Router().ErrorController(c.Context(), status).(iface.IErrorController).Error(status, "")
	}

	if status == http.StatusOK {
		return
	}

	switch pErrorType {
	case perror.ErrTypeWarn:
		c.Context().Warn("%s, trace[%s]", util.ToString(v), util.PanicTrace(TraceMaxDepth, false, debug))
	case perror.ErrTypeIgnore:
	default:
		c.Context().Error("%s, trace[%s]", util.ToString(v), util.PanicTrace(TraceMaxDepth, false, debug))
	}

}

// Response response values action returned
func (c *Controller) Response(v interface{}, err error) {
	if err != nil {
		if pErr, ok := err.(*perror.Error); ok {
			switch pErr.ErrType() {
			case perror.ErrTypeError:
				c.Context().Error(pErr.Error())
			case perror.ErrTypeWarn:
				c.Context().Warn(pErr.Error())
			}

			c.Json(nil, pErr.Status(), pErr.Message())
			return
		}

		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	if r, isRender := v.(render.Render); isRender {
		c.Render(r)
		return
	}

	c.Json(v, 200)
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
	out := map[string]interface{}{
		"data": data,
	}
	c.json(out, status, msg...)
}

// JsonV2 output json response
func (c *Controller) JsonV2(data interface{}, status int, msg ...string) {
	out := map[string]interface{}{
		"data":       data,
		"serverTime": float64(time.Now().UnixNano()) / 1e9,
	}
	c.json(out, status, msg...)
}

// Json output json response
func (c *Controller) json(out map[string]interface{}, status int, msg ...string) {
	ctx := c.Context()
	message := App().Status().Text(status, ctx.Header("Accept-Language", ""), msg...)
	out["status"] = status
	out["message"] = message
	r := render.NewJson(out)
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
func (c *Controller) Data(data []byte, dftContentType ...string) {
	ctx := c.Context()
	r := render.NewData(data)
	httpStatus := r.HttpCode()
	if ctx.Status() > 0 && ctx.Status() != httpStatus {
		httpStatus = ctx.Status()
	}

	ctx.PushLog("status", httpStatus)
	contentType := ""
	if len(dftContentType) > 0 {
		contentType = dftContentType[0]
	}

	if contentType != "" {
		ctx.SetHeader("Content-Type", contentType)
	}

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

// SetActionDesc
// Deprecated: Delete the next version directly
func (c *Controller) SetActionDesc(message string) {

}
