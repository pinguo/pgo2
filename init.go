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
		"help": {Name: "help", Usage: "Displays a list of CMD controllers used (optional), eg. --help=1"},
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
	startTime := float64(time.Now().UnixNano()/1e6)
	defer GLogger().Info(fmt.Sprintf("exec cmdList cost:%f ms",float64(time.Now().UnixNano()/1e6)-startTime))
	list := App().Router().CmdHandlers()

	showText:= func(name,nameType,DefValue,usage string ) string{
		s := fmt.Sprintf("    \t  --%s", name)
		if len(nameType) > 0 {
			s += " " + nameType
			s += "    \t"
		}
		s += strings.ReplaceAll(usage, "\n", "\n    \t")
		if DefValue != "" {
			s += fmt.Sprintf(" (default %v)", DefValue)
		}

		return s
	}

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

			nameType, usage := flag.UnquoteUsage(vFlag)

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
			defValue := ""
			if !isZeroValue(vFlag, vFlag.DefValue) {
				defValue = vFlag.DefValue
			}
			s := showText(vFlag.Name,nameType,defValue, usage)
			retStr += s + "\n"

		})

		return retStr
	}

	fmt.Println("Global parameters:\n " + flagParams(true))

	fmt.Println("The --cmd path list:")
	actionInfo := func(handler *Handler) (actionDesc, paramsMsg string) {
		defer func() {
			if err := recover(); err != nil {
				paramsMsg = path + ":不能解析参数，err:" + util.ToString(err)
			}
		}()
		ctx := &Context{}
		ctx.setPath(path)
		rv := App().Container().getNew(GetAlias(handler.cPath))

		if !rv.IsValid() {
			return
		}

		methodT := rv.Type().Method(handler.aId)

		actionInfo := NewParser().GetActionInfo(rv.Type().Elem().PkgPath(),rv.Type().Elem().Name(),methodT.Name)

		if actionInfo == nil {
			return
		}
		actionDesc = actionInfo.Desc
		if actionInfo.ParamsDesc == nil {
			return
		}


		for _,v:= range actionInfo.ParamsDesc{
			paramsMsg = paramsMsg + showText(v.Name,v.NameType,v.DftValue,v.Usage) + "\n"
		}

		return

	}

	if path != "" {
		path = strings.ToLower(util.CleanPath(path))
	}

	for uri, handler := range list {
		if path != "" && path != strings.ToLower(uri) {
			continue
		}

		actionDesc,paramsStr := actionInfo(handler)
		fmt.Println("  --cmd=" + uri + " \t" + actionDesc)
		if App().CmdMode() && paramsStr != "" {
			fmt.Println(paramsStr)
		}

	}
	fmt.Println("")
	fmt.Println("")
}
