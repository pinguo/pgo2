package pgo2

import (
    "fmt"
    "net/http"
    "net/url"
    "os"
    "reflect"
    "strings"
    "time"

    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/logs"
    "github.com/pinguo/pgo2/perror"
    "github.com/pinguo/pgo2/util"
    "github.com/pinguo/pgo2/validate"
)

type objectItem struct {
    name string
    rv   reflect.Value
}

// Context pgo request context, context is not goroutine
// safe, copy context to use in other goroutines
type Context struct {
    enableAccessLog bool
    response        Response
    input           *http.Request
    output          http.ResponseWriter

    startTime    time.Time
    controllerId string
    actionId     string
    userData     map[string]interface{}
    plugins      []iface.IPlugin
    index        int
    objects      []objectItem

    queryCache url.Values

    logs.Profiler
    logs.Logger
}

func (c *Context) HttpRW(enableAccessLog bool, r *http.Request, w http.ResponseWriter) {
    c.enableAccessLog = enableAccessLog
    c.input = r
    c.output = &c.response
    c.response.reset(w)
}

// start plugin chain process
func (c *Context) Process(plugins []iface.IPlugin) {
    // generate request id
    logId := c.Header("X-Log-Id", "")
    if logId == "" {
        logId = util.GenUniqueId()
    }

    // reset properties
    c.startTime = time.Now()
    c.controllerId = ""
    c.actionId = ""
    c.userData = nil
    c.plugins = plugins
    c.index = -1
    c.queryCache = nil
    c.Profiler.Reset()
    c.Logger.Init(App().Name(), logId, App().Log())

    // finish response
    defer c.finish(false)

    // process request
    c.Next()

}

// finish process request
func (c *Context) finish(goLog bool) {
    // process unhandled panic
    if v := recover(); v != nil {
       status := http.StatusInternalServerError
       switch e := v.(type) {
       case *perror.Error:
           status = e.Status()
           c.End(status, []byte(App().Status().Text(status, c.Header("Accept-Language", ""), e.Message())))
       default:
           c.End(status, []byte(http.StatusText(status)))
       }

       c.Error("%s, trace[%s]", util.ToString(v), util.PanicTrace(TraceMaxDepth, false))
    }

    if !goLog {
        // write header if not yet
        c.response.finish()

        // write access log
        if c.enableAccessLog {
            c.Notice("%s %s %d %d %dms pushlog[%s] profile[%s] counting[%s]",
                c.Method(), c.Path(), c.Status(), c.Size(), c.ElapseMs(),
                c.PushLogString(), c.ProfileString(), c.CountingString())
        }
    } else {
        if c.enableAccessLog {
            c.Notice("%dms pushlog[%s] profile[%s] counting[%s]",
                c.ElapseMs(), c.PushLogString(), c.ProfileString(), c.CountingString())
        }
    }
    // clean objects
    c.clean()
}

func (c *Context) Notice(format string, v ...interface{}) {
    c.Logger.Notice(format, v...)
}

func (c *Context) Debug(format string, v ...interface{}) {
    c.Logger.Debug(format, v...)
}

func (c *Context) Info(format string, v ...interface{}) {
    c.Logger.Info(format, v...)
}

func (c *Context) Warn(format string, v ...interface{}) {
    c.Logger.Warn(format, v...)
}

func (c *Context) Error(format string, v ...interface{}) {
    c.Logger.Error(format, v...)
}

func (c *Context) Fatal(format string, v ...interface{}) {
    c.Logger.Error(format, v...)
}

// finish coroutine log
func (c *Context) FinishGoLog() {
    c.finish(true)
}

// cache object in context
func (c *Context) Cache(name string, rv reflect.Value) {
    if App().Mode() == ModeWeb && len(c.objects) < MaxCacheObjects {
        c.objects = append(c.objects, objectItem{name, rv})
    }
}

// clean all cached objects
func (c *Context) clean() {
    container, num := App().Container(), len(c.objects)
    for i := 0; i < num; i++ {
        name, rv := c.objects[i].name, c.objects[i].rv
        container.Put(name, rv)
    }

    // reset object pool to empty
    if num > 0 {
        c.objects = c.objects[:0]
    }
}

// Next start running plugin chain
func (c *Context) Next() {
    c.index++
    for num := len(c.plugins); c.index < num; c.index++ {
        c.plugins[c.index].HandleRequest(c)
    }
}

// Abort stop running plugin chain
func (c *Context) Abort() {
    c.index = MaxPlugins
}

// Copy copy context
func (c *Context) Copy() iface.IContext {
    cp := *c
    cp.Profiler.Reset()
    cp.userData = nil
    cp.plugins = nil
    cp.index = MaxPlugins
    cp.objects = nil
    return &cp
}

// ElapseMs get elapsed ms since request start
func (c *Context) ElapseMs() int {
    elapse := time.Now().Sub(c.startTime)
    return int(elapse.Nanoseconds() / 1e6)
}

// LogId get log id of current context
func (c *Context) LogId() string {
    return c.Logger.LogId()
}

// Status get response status
func (c *Context) Status() int {
    return c.response.status
}

// Size get response size
func (c *Context) Size() int {
    return c.response.size
}

// SetInput
func (c *Context) SetInput(r *http.Request) {
    c.input = r
}

// SetInput
func (c *Context) Input() *http.Request {
    return c.input
}

// SetOutput
func (c *Context) SetOutput(w http.ResponseWriter) {
    c.output = w
}

// Output
func (c *Context) Output() http.ResponseWriter {
    return c.output
}

// SetControllerId
func (c *Context) SetControllerId(id string) {
    c.controllerId = id
}

// ControllerId
func (c *Context) ControllerId() string {
    return c.controllerId
}

// SetActionId
func (c *Context) SetActionId(id string) {
    c.actionId = id
}

// ActionId
func (c *Context) ActionId() string {
    return c.actionId
}

// SetUserData set user data to current context
func (c *Context) SetUserData(key string, data interface{}) {
    if nil == c.userData {
        c.userData = make(map[string]interface{})
    }

    c.userData[key] = data
}

// UserData get user data from current context
func (c *Context) UserData(key string, dft interface{}) interface{} {
    if data, ok := c.userData[key]; ok {
        return data
    }

    return dft
}

// Method get request method
func (c *Context) Method() string {
    if c.input != nil {
        return c.input.Method
    }

    return "CMD"
}

// genQueryCache query cache
func (c *Context) genQueryCache() {
    if c.input != nil && c.queryCache == nil {
        c.queryCache = c.input.URL.Query()
    }
}

// Query get first url query value by name
func (c *Context) Query(name, dft string) string {
    if c.input != nil {
        c.genQueryCache()
        v := c.queryCache.Get(name)
        if len(v) > 0 {
            return v
        }
    }

    return dft
}

// QueryAll get first value of all url queries
func (c *Context) QueryAll() map[string]string {
    m := make(map[string]string)
    if c.input != nil {
        c.genQueryCache()
        if c.queryCache != nil {
            for k, v := range c.queryCache {
                if len(v) > 0 {
                    m[k] = v[0]
                } else {
                    m[k] = ""
                }
            }
        }

    }

    return m
}

// Post get first post value by name
func (c *Context) Post(name, dft string) string {
    if c.input != nil {
        v := c.input.PostFormValue(name)
        if len(v) > 0 {
            return v
        }
    }

    return dft
}

// PostAll get first value of all posts
func (c *Context) PostAll() map[string]string {
    m := make(map[string]string)
    if c.input != nil {
        // make sure c.input.ParseMultipartForm has been called
        c.input.PostFormValue("")
        for k, v := range c.input.PostForm {
            if len(v) > 0 {
                m[k] = v[0]
            } else {
                m[k] = ""
            }
        }
    }

    return m
}

// Param get first param value by name, post take precedence over get
func (c *Context) Param(name, dft string) string {
    if c.input != nil {
        v := c.input.FormValue(name)
        if len(v) > 0 {
            return v
        }
    }

    return dft
}

// ParamAll get first value of all params, post take precedence over get
func (c *Context) ParamAll() map[string]string {
    m := make(map[string]string)
    if c.input != nil {
        // make sure c.input.ParseMultipartForm has been called
        c.input.FormValue("")
        for k, v := range c.input.Form {
            if len(v) > 0 {
                m[k] = v[0]
            } else {
                m[k] = ""
            }
        }
    }

    return m
}

// ParamMap get map value from GET/POST
func (c *Context) ParamMap(name string) map[string]string {
    // name[k1]=v1&name[k2]=v2
    if c.input != nil {
        c.input.FormValue("")
        ret, exist := c.getMap(c.input.Form, name)
        if exist == true {
            return ret
        }
    }
    return nil
}

// QueryMap get map value from GET
func (c *Context) QueryMap(name string) map[string]string {
    // name[k1]=v1&name[k2]=v2
    if c.input != nil {
        c.genQueryCache()
        ret, exist := c.getMap(c.queryCache, name)
        if exist == true {
            return ret
        }
    }
    return nil
}

// PostMap get map value from POST
func (c *Context) PostMap(name string) map[string]string {
    // name[k1]=v1&name[k2]=v2
    if c.input != nil {
        // make sure c.input.ParseMultipartForm has been called
        c.input.PostFormValue("")
        ret, exist := c.getMap(c.input.PostForm, name)
        if exist == true {
            return ret
        }
    }

    return nil
}

// getMap get map value from map
func (c Context) getMap(m url.Values, key string) (map[string]string, bool) {
    ret := make(map[string]string)
    exist := false
    if m == nil {
        return ret, exist
    }

    for k, v := range m {
        if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
            if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
                exist = true
                ret[k[i+1:][:j]] = v[0]
            }
        }
    }
    return ret, exist
}

// ParamArray get array value from GET/POST
func (c *Context) ParamArray(name string) []string {
    // name[]=v1&name[]=v2

    if c.input != nil {
        // make sure c.input.ParseMultipartForm has been called
        c.input.FormValue("")
        if vv, has := c.input.Form[name]; has == true {
            return vv
        }
    }
    return nil
}

// QueryArray get array value from GET
func (c *Context) QueryArray(name string) []string {
    // name[]=v1&name[]=v2
    if c.input != nil {
        c.genQueryCache()
        if vv, has := c.queryCache[name]; has == true {
            return vv
        }
    }
    return nil
}

// PostArray get array value from POST
func (c *Context) PostArray(name string) []string {
    // name[]=v1&name[]=v2
    if c.input != nil {
        // make sure c.input.ParseMultipartForm has been called
        c.input.PostFormValue("")
        if vv, has := c.input.PostForm[name]; has == true {
            return vv
        }
    }
    return nil
}

// Cookie get first cookie value by name
func (c *Context) Cookie(name, dft string) string {
    if c.input != nil {
        v, e := c.input.Cookie(name)
        if e == nil && len(v.Value) > 0 {
            return v.Value
        }
    }

    return dft
}

// CookieAll get first value of all cookies
func (c *Context) CookieAll() map[string]string {
    m := make(map[string]string)
    if c.input != nil {
        cookies := c.input.Cookies()
        for _, cookie := range cookies {
            if _, ok := m[cookie.Name]; !ok {
                m[cookie.Name] = cookie.Value
            }
        }
    }

    return m
}

// Header get first header value by name，name is case-insensitive
func (c *Context) Header(name, dft string) string {
    if c.input != nil {
        v := c.input.Header.Get(name)
        if len(v) > 0 {
            return v
        }
    }

    return dft
}

// HeaderAll get first value of all headers
func (c *Context) HeaderAll() map[string]string {
    m := make(map[string]string)
    if c.input != nil {
        for k, v := range c.input.Header {
            if len(v) > 0 {
                m[k] = v[0]
            } else {
                m[k] = ""
            }
        }
    }

    return m
}

// Path get request path
func (c *Context) Path() string {
    // for web
    if c.input != nil {
        return c.input.URL.Path
    }

    // for cmd
    if App().HasArg("cmd") && len(App().Arg("cmd")) > 0 {
        return App().Arg("cmd")
    }

    return "/"
}

// ClientIp get client ip
func (c *Context) ClientIp() string {
    if xff := c.Header("X-Forwarded-For", ""); len(xff) > 0 {
        if pos := strings.IndexByte(xff, ','); pos > 0 {
            return strings.TrimSpace(xff[:pos])
        } else {
            return xff
        }
    }

    if ip := c.Header("X-Client-Ip", ""); len(ip) > 0 {
        return ip
    }

    if ip := c.Header("X-Real-Ip", ""); len(ip) > 0 {
        return ip
    }

    if c.input != nil && len(c.input.RemoteAddr) > 0 {
        pos := strings.LastIndexByte(c.input.RemoteAddr, ':')
        if pos > 0 {
            return c.input.RemoteAddr[:pos]
        } else {
            return c.input.RemoteAddr
        }
    }

    return ""
}

// validate query param, return string validator
func (c *Context) ValidateQuery(name string, dft ...interface{}) *validate.String {
    return validate.StringData(c.Query(name, ""), name, dft...)
}

// validate post param, return string validator
func (c *Context) ValidatePost(name string, dft ...interface{}) *validate.String {
    return validate.StringData(c.Post(name, ""), name, dft...)
}

// validate get/post param, return string validator
func (c *Context) ValidateParam(name string, dft ...interface{}) *validate.String {
    return validate.StringData(c.Param(name, ""), name, dft...)
}

// set response header, no effect if any header has sent
func (c *Context) SetHeader(name, value string) {
    if c.output != nil {
        c.output.Header().Set(name, value)
    }
}

// convenient way to set response cookie
func (c *Context) SetCookie(cookie *http.Cookie) {
    if c.output != nil {
        http.SetCookie(c.output, cookie)
    }
}

// send response
func (c *Context) End(status int, data []byte) {
    if c.output != nil {
        c.SetHeader("X-Log-Id", c.Logger.LogId())
        c.SetHeader("X-Cost-Time", fmt.Sprintf("%dms", c.ElapseMs()))
        c.output.WriteHeader(status)
        c.output.Write(data)
    } else if len(data) > 0 {
        os.Stdout.Write(data)
        os.Stdout.WriteString("\n")
    }
}
