package iface

import (
    "html/template"
    "io"
    "net/http"
    "reflect"
    "time"

    "github.com/pinguo/pgo2/validate"
)

//
//import "time"

type IBind interface {
    GetBindInfo(v interface{}) interface{}
}

type IObject interface {
    SetContext(ctx IContext)
    Context() IContext
    GetObj(obj IObject) IObject
    GetObjPool(funcName IObjPoolFunc, params ...interface{}) IObject
    GetObjSingle(name string, funcName IObjSingleFunc, params ...interface{}) IObject
    GetObjPoolCtr(ctr IContext, funcName IObjPoolFunc, params ...interface{}) IObject
    GetObjCtr(ctr IContext, obj IObject) IObject
    GetObjSingleCtr(ctr IContext, name string, funcName IObjSingleFunc, params ...interface{}) IObject
}

type IController interface {
    BeforeAction(action string)
    AfterAction(action string)
    HandlePanic(v interface{})
}

type IPlugin interface {
    HandleRequest(ctx IContext)
}

type IEvent interface {
    HandleEvent(event string, ctx IContext, args ...interface{})
}

//
//
//
//type ICache interface {
//   Get(key string) *Value
//   MGet(keys []string) map[string]*Value
//   Set(key string, value interface{}, expire ...time.Duration) bool
//   MSet(items map[string]interface{}, expire ...time.Duration) bool
//   Add(key string, value interface{}, expire ...time.Duration) bool
//   MAdd(items map[string]interface{}, expire ...time.Duration) bool
//   Del(key string) bool
//   MDel(keys []string) bool
//   Exists(key string) bool
//   Incr(key string, delta int) int
//}

type IStatus interface {
    Text(status int, lang string, dft ...string) string
}

type II18n interface {
    Translate(message, lang string, params ...interface{}) string
}

type IView interface {
    AddFuncMap(funcMap template.FuncMap)
    Render(view string, data interface{}) []byte
    Display(w io.Writer, view string, data interface{})
}

type IContext interface {
    HttpRW(enableAccessLog bool, r *http.Request, w http.ResponseWriter)
    Process(plugins []IPlugin)
    Notice(format string, v ...interface{})
    Debug(format string, v ...interface{})
    Info(format string, v ...interface{})
    Warn(format string, v ...interface{})
    Error(format string, v ...interface{})
    Fatal(format string, v ...interface{})
    FinishGoLog()
    Cache(name string, rv reflect.Value)
    Next()
    Abort()
    Copy() IContext
    ElapseMs() int
    LogId() string
    Status() int
    Size() int
    SetInput(r *http.Request)
    Input() *http.Request
    SetOutput(w http.ResponseWriter)
    Output() http.ResponseWriter
    SetControllerId(id string)
    ControllerId() string
    SetActionId(id string)
    ActionId() string
    SetUserData(key string, data interface{})
    UserData(key string, dft interface{}) interface{}
    Method() string
    Query(name, dft string) string
    QueryAll() map[string]string
    Post(name, dft string) string
    PostAll() map[string]string
    Param(name, dft string) string
    ParamAll() map[string]string
    ParamMap(name string) map[string]string
    QueryMap(name string) map[string]string
    PostMap(name string) map[string]string
    ParamArray(name string) []string
    QueryArray(name string) []string
    PostArray(name string) []string
    Cookie(name, dft string) string
    CookieAll() map[string]string
    Header(name, dft string) string
    HeaderAll() map[string]string
    Path() string
    ClientIp() string
    ValidateQuery(name string, dft ...interface{}) *validate.String
    ValidatePost(name string, dft ...interface{}) *validate.String
    ValidateParam(name string, dft ...interface{}) *validate.String
    SetHeader(name, value string)
    SetCookie(cookie *http.Cookie)
    End(status int, data []byte)
    PushLog(key string, v interface{})
    Counting(key string, hit, total int)
    ProfileStart(key string)
    ProfileStop(key string)
    ProfileAdd(key string, elapse time.Duration)
    PushLogString() string
    CountingString() string
    ProfileString() string
}

type IObjPoolFunc func(ctr IContext, params ...interface{}) IObject
type IObjSingleFunc func(params ...interface{}) IObject
type IComponentFunc func(config map[string]interface{}) (interface{}, error)
