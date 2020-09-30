package pgo2

import (
	"flag"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/pinguo/pgo2/logs"
	"github.com/pinguo/pgo2/util"
)

const (
	ModeWeb                = "web"
	ModeCmd                = "cmd"
	DefaultEnv             = "develop"
	DefaultControllerPath  = "index"
	DefaultActionPath      = "index"
	DefaultHttpAddr        = "0.0.0.0:8000"
	DefaultTimeout         = 30 * time.Second
	DefaultHeaderBytes     = 1 << 20
	ControllerWebPkg       = "controller"
	ControllerCmdPkg       = "command"
	ControllerWebType      = "Controller"
	ControllerCmdType      = "Command"
	ConstructMethod        = "Construct"
	PrepareMethod          = "Prepare"
	VendorPrefix           = "vendor/"
	VendorLength           = 7
	ActionPrefix           = "Action"
	ActionLength           = 6
	TraceMaxDepth          = 10
	MaxPlugins             = 32
	MaxCacheObjects        = 100
	ParamsFlagMethodPrefix = "ParamsFlag"
)

var (
	application    *Application
	appTime        = time.Now()
	aliasMap       = make(map[string]string)
	aliasRe        = regexp.MustCompile(`^@[^\\/]+`)
	logger         *logs.Logger
	EmptyObject    struct{}
	restFulActions = map[string]int{"GET": 1, "POST": 1, "PUT": 1, "DELETE": 1, "PATCH": 1, "OPTIONS": 1, "HEAD": 1}
	globalParams   = map[string]*flag.Flag{
		"env":  {Name: "env", Usage: "set running env (optional), eg. --env=online"},
		"cmd":  {Name: "cmd", Usage: "set running cmd (optional), eg. --cmd=/foo/bar"},
		"base": {Name: "base", Usage: "set base path (optional), eg. --base=/base/path"},
		"help": {Name: "help",DefValue:"1", Usage: "Displays a list of CMD controllers used (optional), eg. --help=1"},
	}
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
	if err := App().Config().CheckPath(); err != nil {
		cmdList("")
		panic(err)
	}
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

func cmdList(path string) {
	list := App().Router().CmdHandlers()

	flagParams := func(global bool) string {
		retStr := ""
		flag.VisitAll(func(vFlag *flag.Flag) {
			_, hasGP := globalParams[vFlag.Name]
			if global && !hasGP {
				return
			}

			if !global && hasGP {
				return
			}
			s := fmt.Sprintf("    \t  --%s", vFlag.Name) // Two spaces before -; see next two comments.
			name, usage := flag.UnquoteUsage(vFlag)
			if len(name) > 0 {
				s += " " + name
			}

			s += "    \t"
			s += strings.ReplaceAll(usage, "\n", "\n    \t")
			isZeroValue := func(vFlag *flag.Flag, value string) bool {
				// Build a zero value of the flag's Value type, and see if the
				// result of calling its String method equals the value passed in.
				// This works unless the Value type is itself an interface type.
				typ := reflect.TypeOf(vFlag.Value)
				var z reflect.Value
				if typ.Kind() == reflect.Ptr {
					z = reflect.New(typ.Elem())
				} else {
					z = reflect.Zero(typ)
				}
				return value == z.Interface().(flag.Value).String()
			}
			if !isZeroValue(vFlag, vFlag.DefValue) {
				s += fmt.Sprintf(" (default %v)", vFlag.DefValue)
			}
			retStr += s + "\n"

		})

		return retStr
	}

	fmt.Println("Global parameters:\n " + flagParams(true))

	fmt.Println("The path list:")
	showParams := func(path string) string {
		ctx := &Context{}
		ctx.setPath(path)

		rv, _, _ := App().Router().CreateController(path, ctx)
		if !rv.IsValid() {
			return ""
		}

		name := ParamsFlagMethodPrefix + ctx.ActionId()
		methodV := rv.MethodByName(name)
		if !methodV.IsValid() {
			return ""
		}

		methodV.Call(nil)
		return flagParams(false)
	}

	if path != "" {
		path = strings.ToLower(util.CleanPath(path))
	}

	for uri, v := range list {
		// fmt.Println("path", path,"uri",uri)

		if path != "" && path != strings.ToLower(uri) {
			continue
		}

		paramsStr:= showParams(uri)
		fmt.Println("  --cmd=" + uri + " \t" + v.desc)
		if paramsStr != "" {
			fmt.Println(paramsStr)
		}


	}
	fmt.Println("")
	fmt.Println("")
}
