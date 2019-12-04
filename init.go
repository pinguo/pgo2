package pgo2

import (
    "regexp"
    "strings"
    "time"

    "github.com/pinguo/pgo2/logs"
    "github.com/pinguo/pgo2/util"
)

const (
    ModeWeb               = "web"
    ModeCmd               = "cmd"
    DefaultEnv            = "develop"
    DefaultControllerPath = "index"
    DefaultActionPath     = "index"
    DefaultHttpAddr       = "0.0.0.0:8000"
    DefaultTimeout        = 30 * time.Second
    DefaultHeaderBytes    = 1 << 20
    ControllerWebPkg      = "controller"
    ControllerCmdPkg      = "command"
    ControllerWebType     = "Controller"
    ControllerCmdType     = "Command"
    ConstructMethod       = "Construct"
    InitMethod            = "Init"
    VendorPrefix          = "vendor/"
    VendorLength          = 7
    ActionPrefix          = "Action"
    ActionLength          = 6
    TraceMaxDepth         = 10
    MaxPlugins            = 32
    MaxCacheObjects       = 100
)

var (
    application *Application
    appTime     = time.Now()
    aliasMap    = make(map[string]string)
    aliasRe     = regexp.MustCompile(`^@[^\\/]+`)
    logger      *logs.Logger
    EmptyObject struct{}
)

func App(newApp ...bool) *Application {
    if application != nil && newApp == nil {
        return application
    }

    application = NewApp()

    return application
}

// Run run app
func Run() {
    // Initialization route
    App().Router().InitHandlers()
    // Check config path
    App().Config().CheckPath()
    // Listen for server or start CMD
    App().Server().Serve()
}

// GLogger get global logger
func GLogger() *logs.Logger {
    if logger == nil {
        // defer creation to first call, give opportunity to customize log target
        logger = App().Log().Logger(App().name, util.GenUniqueId())
    }

    return logger
}

// SetAlias set path alias, eg. @app => /path/to/base
func SetAlias(alias, path string) {
    if len(alias) > 0 && alias[0] != '@' {
        alias = "@" + alias
    }

    if strings.IndexAny(alias, `\/`) != -1 {
        panic("SetAlias: invalid alias, " + alias)
    }

    if len(alias) <= 1 || len(path) == 0 {
        panic("SetAlias: empty alias or path, " + alias)
    }

    aliasMap[alias] = path
}

// GetAlias resolve path alias, eg. @runtime/app.log => /path/to/runtime/app.log
func GetAlias(alias string) string {
    if prefix := aliasRe.FindString(alias); len(prefix) == 0 {
        return alias // not an alias
    } else if path, ok := aliasMap[prefix]; ok {
        return strings.Replace(alias, prefix, path, 1)
    }

    return ""
}
