package pgo2

import (
    "reflect"
    "regexp"
    "strings"

    "github.com/pinguo/pgo2/core"
    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/util"
)

// format route string to CamelCase, eg.
// /api/foo-bar/say-hello => /Api/FooBar/SayHello
func routeFormatFunc(s string) string {
    s = strings.ToUpper(s)
    if s[0] == '-' {
        s = s[1:]
    }
    return s
}

// format route string to CamelCase, eg.
// /path/FooBar/SayHello => /api/foo-bar/say-hello
func pathFormatFunc(s string) string {
    s = strings.ToLower(s)
    return "-" + s
}

type routeRule struct {
    rePat   *regexp.Regexp
    pattern string
    route   string
}

// Router the router component, configuration:
// router:
//     rules:
//         - "^/foo/all$ => /foo/index"
//         - "^/api/user/(\d+)$ => /api/user"
func NewRouter(config map[string]interface{}) *Router {
    router := &Router{}
    router.reFmt = regexp.MustCompile(`([/-][a-z])`)
    router.rePathFmt = regexp.MustCompile(`([A-Z])`)
    router.rules = make([]routeRule, 0, 10)

    core.Configure(router, config)

    return router
}

type Handler struct {
    uri   string
    cPath string
    cId   string
    aName string
    aId   int
}

type Router struct {
    reFmt     *regexp.Regexp
    rePathFmt *regexp.Regexp
    rules     []routeRule

    webHandlers map[string]*Handler
    cmdHandlers map[string]*Handler
    modules     []string
}

var rePath = strings.NewReplacer("/"+ControllerCmdPkg+"/", "/", "/"+ControllerWebPkg+"/", "/", ControllerCmdType, "", ControllerWebType, "")

// SetRules set rule list, format: `^/api/user/(\d+)$ => /api/user`
func (r *Router) SetRules(rules []interface{}) {
    for _, v := range rules {
        parts := strings.Split(v.(string), "=>")
        if len(parts) != 2 {
            panic("Router: invalid rule: " + util.ToString(v))
        }

        pattern := strings.TrimSpace(parts[0])
        route := strings.TrimSpace(parts[1])
        r.AddRoute(pattern, route)
    }
}

// InitHandlers Initialization route
func (r *Router) InitHandlers() {
    r.webHandlers = make(map[string]*Handler)
    r.cmdHandlers = make(map[string]*Handler)
    webList := App().Container().PathList(ControllerWebPkg+"/", ControllerWebType)
    cmdList := App().Container().PathList(ControllerCmdPkg+"/", ControllerCmdType)
    r.SetHandlers(ControllerWebPkg, webList)
    r.SetHandlers(ControllerCmdPkg, cmdList)

}

// SetHandlers Set route
func (r *Router) SetHandlers(cmdType string, list map[string]interface{}) {
    if list == nil {
        return
    }

    for controllerOPath, info := range list {
        actions, _ := info.(map[string]int)
        controllerPath := rePath.Replace("/" + controllerOPath)

        paths := strings.Split(controllerPath, "/")
        oCname := paths[len(paths)-1:][0]
        cName := r.firstToLower(oCname)

        baseUrl := strings.Join(paths[0:len(paths)-1], "/") + "/"
        baseUrl = strings.Replace(baseUrl, "//", "/", -1)

        cNames := make([]string, 0, 2)
        cNames = append(cNames, cName)
        if r.web(cmdType) {
            fmtCName := r.rePathFmt.ReplaceAllStringFunc(cName, pathFormatFunc)
            if cName != fmtCName {
                cNames = append(cNames, fmtCName)
            }
        }

        for oAName, aNum := range actions {
            aNames := make([]string, 0, 2)
            aName := r.firstToLower(oAName)
            aNames = append(aNames, aName)
            if r.web(cmdType) {
                fmtAName := r.rePathFmt.ReplaceAllStringFunc(aName, pathFormatFunc)
                if aName != fmtAName {
                    aNames = append(aNames, fmtAName)
                }
            }

            for _, cPath := range cNames {
                for _, aPath := range aNames {
                    uri := baseUrl + cPath + "/" + aPath
                    r.setHandler(cmdType, uri, controllerOPath, baseUrl+cPath, oAName, aNum)
                }
            }
        }
    }
}

func (r *Router) web(cmdType string) bool {
    return cmdType == ControllerWebPkg
}

func (r *Router) firstToLower(s string) string {
    return strings.ToLower(s[0:1]) + s[1:]
}

func (r *Router) setHandler(cmdType, uri, cPath, cId, aName string, aNum int) {
    switch cmdType {
    case ControllerWebPkg:
        r.webHandlers[uri] = &Handler{uri: uri, cPath: cPath, cId: cId, aName: aName, aId: aNum}
    case ControllerCmdPkg:
        r.cmdHandlers[uri] = &Handler{uri: uri, cPath: cPath, cId: cId, aName: aName, aId: aNum}
    default:
        panic("Is the defined cmdType")
    }
}

// AddRoute add one route, the captured group will be passed to
// action method as function params
func (r *Router) AddRoute(pattern, route string) {
    rePat := regexp.MustCompile(pattern)
    rule := routeRule{rePat, pattern, route}
    r.rules = append(r.rules, rule)
}

// Resolve path to route and action params, then format route to CamelCase
func (r *Router) Resolve(path string) (handler *Handler, params []string) {
    // The first mapping
    handler = r.Handler(path)
    if handler != nil {
        return
    }

    // format path
    if path == "/" {
        path += DefaultControllerPath + "/" + DefaultActionPath
    } else {
        if r.modules != nil && util.SliceSearchString(r.modules, path) > 0 {
            path += "/" + DefaultControllerPath + "/" + DefaultActionPath
        }
    }

    path = util.CleanPath(path)
    // The second mapping
    handler = r.Handler(path)
    if handler != nil {
        return
    }

    // Custom route
    if len(r.rules) != 0 {
        for _, rule := range r.rules {
            matches := rule.rePat.FindStringSubmatch(path)
            if len(matches) != 0 {
                path = rule.route
                params = matches[1:]
                break
            }
        }
    }
    handler = r.Handler(path)

    return
}

func (r *Router) Handler(path string) *Handler {
    if ModeWeb == App().mode {
        if handler, ok := r.webHandlers[path]; ok {
            return handler
        }
    } else {
        if handler, ok := r.cmdHandlers[path]; ok {
            return handler
        }
    }
    return nil
}

func (r *Router) CmdHandlers() map[string]*Handler {
    return r.cmdHandlers
}

// CreateController Create the controller and parameters
func (r *Router) CreateController(path string, ctx iface.IContext) (reflect.Value, reflect.Value, []string) {

    container := App().Container()

    handler, params := r.Resolve(path)
    if handler == nil {
        return reflect.Value{}, reflect.Value{}, nil
    }

    controllerName := handler.cPath

    ctx.SetControllerId(handler.cId)

    ctx.SetActionId(handler.aName)

    controllerName = GetAlias(controllerName)

    controller := container.Get(controllerName, ctx)
    action := controller.Method(handler.aId)
    return controller, action, params
}
