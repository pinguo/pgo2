package pgo2

import (
    "net/http"
    "reflect"

    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/perror"
    "github.com/pinguo/pgo2/render"
    "github.com/pinguo/pgo2/util"
)

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
func (c *Controller) HandlePanic(v interface{}) {
    status := http.StatusInternalServerError
    switch e := v.(type) {
    case *perror.Error:
        status = e.Status()
        c.Json(EmptyObject, status, e.Message())
    default:
        c.Json(EmptyObject, status)
    }

    c.Context().Error("%s, trace[%s]", util.ToString(v), util.PanicTrace(TraceMaxDepth, false))
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
    ctx.End(r.HttpCode(), r.Content())
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
    ctx.End(r.HttpCode(), r.Content())
}

// Data output data response
func (c *Controller) Data(data []byte) {
    ctx := c.Context()
    r := render.NewData(data)

    ctx.PushLog("status", http.StatusOK)
    ctx.SetHeader("Content-Type", r.ContentType())
    ctx.End(r.HttpCode(), r.Content())
}

// Xml output xml response
func (c *Controller) Xml(data interface{}, status ...int) {
    status = append(status, http.StatusOK)
    ctx := c.Context()
    r := render.NewXml(data)

    ctx.PushLog("status", status[0])
    ctx.SetHeader("Content-Type", r.ContentType())
    ctx.End(r.HttpCode(), r.Content())
}

// ProtoBuf output proto buf response
func (c *Controller) ProtoBuf(data interface{}, status ...int) {
    status = append(status, http.StatusOK)
    ctx := c.Context()
    r := render.NewProtoBuf(data)

    ctx.PushLog("status", status[0])
    ctx.SetHeader("Content-Type", r.ContentType())
    ctx.End(r.HttpCode(), r.Content())
}

// Render Custom renderer
func (c *Controller) Render(r render.Render, status ...int) {
    status = append(status, http.StatusOK)
    ctx := c.Context()
    ctx.PushLog("status", status[0])
    ctx.SetHeader("Content-Type", r.ContentType())
    ctx.End(r.HttpCode(), r.Content())
}

// View output rendered view
func (c *Controller) View(view string, data interface{}, contentType ...string) {
    ctx := c.Context()
    contentType = append(contentType, "text/html; charset=utf-8")
    ctx.PushLog("status", http.StatusOK)
    ctx.SetHeader("Content-Type", contentType[0])
    ctx.End(http.StatusOK, App().View().Render(view, data))
}
