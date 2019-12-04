package pgo2

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "reflect"
    "runtime"
    "strings"
    "sync"

    "github.com/pinguo/pgo2/config"
    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/logs"
)

func init() {
    _ = flag.String("env", "", "set running env (requested), eg. --env=online")
    _ = flag.String("cmd", "", "set running cmd (optional), eg. --cmd=/foo/bar")
    _ = flag.String("base", "", "set base path (optional), eg. --base=/base/path")
    _ = flag.String("cmdList", "", "Displays a list of CMD controllers used (optional), eg. --cmdList")
}

func NewApp() *Application {
    exeBase := filepath.Base(os.Args[0])
    exeExt := filepath.Ext(os.Args[0])
    exeDir := filepath.Dir(os.Args[0])

    app := &Application{
        mode:       ModeWeb,
        env:        DefaultEnv,
        name:       strings.TrimSuffix(exeBase, exeExt),
        components: make(map[string]interface{}),
        objects:    make(map[string]iface.IObject),
    }

    app.basePath = app.genBasePath(exeDir)

    app.Init(os.Args)

    return app
}

type Application struct {
    mode        string // running mode, WEB or CMD
    env         string // running env, eg. develop/online/testing-dev/testing-qa
    name        string // App name
    basePath    string // base path of App
    runtimePath string // runtime path for log etc.
    publicPath  string // public path for web assets
    viewPath    string // path for view template

    config     config.IConfig
    container  *Container
    server     *Server
    router     *Router
    log        *logs.Log
    status     iface.IStatus
    i18n       iface.II18n
    view       iface.IView
    stopBefore *StopBefore // 服务停止前执行 [{"obj":"func"}]

    components map[string]interface{}
    objects    map[string]iface.IObject
    lock       sync.RWMutex

    args map[string]string
}

func (app *Application) initArgs(osArgs []string) map[string]string {
    app.args = make(map[string]string)

    nArg := len(osArgs)
    for k, cmd := range osArgs {
        if k == 0 {
            continue
        }

        if strings.Index(cmd, "--") != 0 {
            continue
        }

        value := ""
        name := strings.Replace(cmd, "--", "", 1)
        if strings.Index(name, "=") > 0 {
            tmpName := strings.Split(name, "=")
            name = tmpName[0]
            value = tmpName[1]
        } else {
            nextK := k + 1
            if nextK < nArg {
                nextV := osArgs[nextK]
                if strings.Index(nextV, "--") != 0 {
                    value = nextV
                }
            }
        }

        app.args[name] = value
    }
    return app.args
}

func (app *Application) HasArg(name string) bool {
    _, has := app.args[name]
    return has
}

func (app *Application) Arg(name string) string {
    if v, has := app.args[name]; has {
        return v
    }
    return ""
}

func (app *Application) cmdList() bool {
    return app.HasArg("cmdList")
}

func (app *Application) Init(osArgs []string) {
    args := app.initArgs(osArgs)

    // overwrite running env
    if env, has := args["env"]; has && len(env) > 0 {
        app.env = env
    }

    // overwrite running mode
    if cmd, has := args["cmd"]; has && len(cmd) > 0 {
        app.mode = ModeCmd
    }

    // overwrite base path
    if base, has := args["base"]; has && len(base) > 0 {
        app.basePath, _ = filepath.Abs(base)
    }

    // set basic path alias
    type dummy struct{}
    pkgPath := reflect.TypeOf(dummy{}).PkgPath()
    SetAlias("@app", app.basePath)
    SetAlias("@pgo2", strings.TrimPrefix(pkgPath, VendorPrefix))

    // initialize config object
    app.config = config.New(app.basePath, app.env)

    // initialize container object
    enablePool, _ := app.config.Get("app.container.enablePool").(string)
    app.container = NewContainer(enablePool)

    // initialize server object
    svrConf, _ := app.config.Get("app.server").(map[string]interface{})
    app.server = NewServer(svrConf)

    // overwrite appName
    if name := app.config.GetString("app.name", ""); len(name) > 0 {
        app.name = name
    }

    // overwrite GOMAXPROCS
    if n := app.config.GetInt("app.GOMAXPROCS", 0); n > 0 {
        runtime.GOMAXPROCS(n)
    }

    // set runtime path
    runtimePath := app.config.GetString("app.runtimePath", "@app/runtime")
    app.runtimePath, _ = filepath.Abs(GetAlias(runtimePath))
    SetAlias("@runtime", app.runtimePath)

    // set public path
    publicPath := app.config.GetString("app.publicPath", "@app/public")
    app.publicPath, _ = filepath.Abs(GetAlias(publicPath))
    SetAlias("@public", app.publicPath)

    // set view path
    viewPath := app.config.GetString("app.viewPath", "@app/view")
    app.viewPath, _ = filepath.Abs(GetAlias(viewPath))
    SetAlias("@view", app.viewPath)

    // create runtime directory if not exists
    if _, e := os.Stat(app.runtimePath); os.IsNotExist(e) {
        if e := os.MkdirAll(app.runtimePath, 0755); e != nil {
            panic(fmt.Sprintf("failed to create %s, %s", app.runtimePath, e))
        }
    }
}

func (app *Application) genBasePath(exeDir string) string {
    basePath, _ := filepath.Abs(filepath.Join(exeDir, ".."))

    return basePath
}

// Mode  running mode, web:1, cmd:2
func (app *Application) Mode() string {
    return app.mode
}

// Env  running env
func (app *Application) Env() string {
    return app.env
}

// Name  appName, default is executable name
func (app *Application) Name() string {
    return app.name
}

// BasePath  base path, default is parent of executable
func (app *Application) BasePath() string {
    return app.basePath
}

// RuntimePath  runtime path, default is @app/runtime
func (app *Application) RuntimePath() string {
    return app.runtimePath
}

// PublicPath  public path, default is @app/public
func (app *Application) PublicPath() string {
    return app.publicPath
}

// ViewPath  view path, default is @app/view
func (app *Application) ViewPath() string {
    return app.viewPath
}

// Config  config component
func (app *Application) Config() config.IConfig {
    return app.config
}

// Container  container component
func (app *Application) Container() *Container {
    return app.container
}

// Server  server component
func (app *Application) Server() *Server {
    return app.server
}

// Router  router component
func (app *Application) Router() *Router {
    if app.router == nil {
        app.router = NewRouter(app.componentConf("router"))
    }

    return app.router
}

// Log  log component
func (app *Application) Log() *logs.Log {
    if app.log == nil {
        app.log = logs.NewLog(app.RuntimePath(), app.componentConf("log"))
    }

    return app.log
}

// Status  status component
func (app *Application) Status() iface.IStatus {
    if app.status == nil {
        app.status = NewStatus(app.componentConf("status"))
    }

    return app.status
}

// I18n  i18n component
func (app *Application) I18n() iface.II18n {
    if app.i18n == nil {
        app.i18n = NewI18n(app.componentConf("i18n"))
    }

    return app.i18n
}

// View  view component
func (app *Application) View() iface.IView {
    if app.view == nil {
        app.view = NewView(app.componentConf("view"))
    }

    return app.view
}

// StopBefore  stopBefore component
func (app *Application) StopBefore() *StopBefore {
    if app.stopBefore == nil {
        app.stopBefore = NewStopBefore()
    }
    return app.stopBefore
}

// SetStatus  set status component
func (app *Application) SetStatus(status iface.IStatus) {
    app.status = status
}

// SetI18n  set i18n component
func (app *Application) SetI18n(i18n iface.II18n) {
    app.i18n = i18n
}

// SetView  set view component
func (app *Application) SetView(view iface.IView) {
    app.view = view
}

// SetView  set view component
func (app *Application) SetConfig(config config.IConfig) {
    app.config = config
}

// component conf by id
func (app *Application) componentConf(id string) map[string]interface{} {
    conf := app.config.Get("app.components." + id)
    if conf == nil {
        return nil
    }

    if retConf, ok := conf.(map[string]interface{}); ok {
        for k, v := range retConf {
            if vv, ok := v.(string); ok == true {
                retConf[k] = GetAlias(vv)
            }
        }

        return retConf
    }

    return nil
}

// Get get component by id
func (app *Application) Component(id string, funcName iface.IComponentFunc, params ...map[string]interface{}) interface{} {
    if _, ok := app.components[id]; !ok {
        if len(params) > 0 {
            obj, err := funcName(params[0])
            if err != nil {
                panic("Component " + id + " err:" + err.Error())
            }

            app.setComponent(id, obj)
        } else {
            obj, err := funcName(app.componentConf(id))
            if err != nil {
                panic("Component " + id + " err:" + err.Error())
            }

            app.setComponent(id, obj)
        }

    }

    app.lock.RLock()
    defer app.lock.RUnlock()

    return app.components[id]
}

// Get get pool class object. name is class name, ctx is context,
func (app *Application) GetObjPool(name string, ctx iface.IContext) iface.IObject {
    if name := GetAlias(name); len(name) > 0 {
        return app.container.Get(name, ctx).Interface().(iface.IObject)
    }

    panic("unknown class: " + name)

}

// Get get single class object. name is class name, ctx is context,
func (app *Application) GetObjSingle(name string, funcName iface.IObjSingleFunc, param ...interface{}) iface.IObject {
    if _, ok := app.objects[name]; !ok {
        app.setObject(name, funcName(param...))
    }

    app.lock.RLock()
    defer app.lock.RUnlock()

    return app.objects[name]
}

func (app *Application) setObject(name string, object iface.IObject, force ...bool) {
    force = append(force, false)

    app.lock.Lock()
    defer app.lock.Unlock()

    // avoid repeated loading
    if _, ok := app.objects[name]; ok && force[0] == false {
        return
    }

    app.objects[name] = object
}

func (app *Application) setComponent(id string, component interface{}, force ...bool) {
    force = append(force, false)

    app.lock.Lock()
    defer app.lock.Unlock()

    // avoid repeated loading
    if _, ok := app.components[id]; ok && force[0] == false {
        return
    }
    app.components[id] = component
}
