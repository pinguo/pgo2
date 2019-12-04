package pgo2

import (
    "reflect"
    "testing"

    "github.com/agiledragon/gomonkey"
    "github.com/pinguo/pgo2/iface"
)

func TestNewRouter(t *testing.T) {
    var router interface{}
    router = NewRouter(nil)
    if _, ok := router.(*Router); ok == false {
        t.FailNow()
    }
}

func TestRouter_SetRules(t *testing.T) {
    router := NewRouter(nil)
    t.Run("err_rules_style", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        rules := []interface{}{`^/api/user/(\d+)$  /api/user`, `^/api/user1/(\d+)$ => /api/user1`}
        router.SetRules(rules)
    })

    t.Run("normal", func(t *testing.T) {
        rules := []interface{}{`^/api/user/(\d+)$ => /api/user`, `^/api/user1/(\d+)$ => /api/user1`}
        router.SetRules(rules)
        if len(router.rules) != len(rules) {
            t.Fatal(`len(router.rules)!=len(rules)`)
        }

    })

}

func TestRouter_AddRoute(t *testing.T) {
    router := NewRouter(nil)
    router.AddRoute(`^/api/user/(\d+)$`, "/api/user")
    if router.rules[0].pattern != `^/api/user/(\d+)$` || router.rules[0].route != "/api/user" {
        t.FailNow()
    }
}

func TestRouter_SetHandlers(t *testing.T) {
    router := NewRouter(nil)
    router.webHandlers = make(map[string]*Handler)
    router.cmdHandlers = make(map[string]*Handler)
    list := make(map[string]interface{})
    list["controller/testController"] = map[string]int{"Index": 0}
    list["controller/testAaController"] = map[string]int{"IndexAaa": 0, "Index": 1}
    router.SetHandlers(ControllerWebPkg, list)

    if len(router.webHandlers) != 7 {
        t.FailNow()
    }

    cmdList := make(map[string]interface{})
    cmdList["controller/testController"] = map[string]int{"Index": 0}
    cmdList["controller/testAaController"] = map[string]int{"IndexAaa": 0, "Index": 1}

    router.SetHandlers(ControllerCmdPkg, list)
    if len(router.CmdHandlers()) != 3 {
        t.FailNow()
    }
}

func TestRouter_Handler(t *testing.T) {
    router := NewRouter(nil)
    router.webHandlers = make(map[string]*Handler)
    router.cmdHandlers = make(map[string]*Handler)

    list := make(map[string]interface{})
    list["controller/testController"] = map[string]int{"Index": 0}
    router.SetHandlers(ControllerWebPkg, list)

    cmdList := make(map[string]interface{})
    cmdList["controller/testController"] = map[string]int{"Index": 0}
    router.SetHandlers(ControllerCmdPkg, list)

    t.Run("web", func(t *testing.T) {
        App().mode = ModeWeb
        if router.Handler("/test/index") == nil {
            t.FailNow()
        }
    })

    t.Run("cmd", func(t *testing.T) {
        App().mode = ModeCmd
        if router.Handler("/test/index") == nil {
            t.FailNow()
        }
    })

    t.Run("nil", func(t *testing.T) {
        App().mode = ModeCmd
        if router.Handler("/test1") != nil {
            t.FailNow()
        }
    })

}

func TestRouter_Resolve(t *testing.T) {
    App(true)
    router := NewRouter(nil)
    router.webHandlers = make(map[string]*Handler)
    router.cmdHandlers = make(map[string]*Handler)

    list := make(map[string]interface{})
    list["controller/IndexController"] = map[string]int{"Index": 0}
    router.SetHandlers(ControllerWebPkg, list)

    cmdList := make(map[string]interface{})
    cmdList["controller/testController"] = map[string]int{"Index": 0}
    router.SetHandlers(ControllerCmdPkg, list)

    t.Run("inHandlers", func(t *testing.T) {
        if h, _ := router.Resolve("/index/index"); h == nil {
            t.FailNow()
        }
    })

    t.Run("path=/", func(t *testing.T) {
        if h, _ := router.Resolve("/"); h == nil {
            t.FailNow()
        }
    })

    t.Run("path=/module/", func(t *testing.T) {
        t.Skip()
    })

    t.Run("path=/api/user/(\\d)", func(t *testing.T) {
        router.AddRoute(`^/api/user/(\d+)$`, "/index/index")
        h, p := router.Resolve("/api/user/123")
        if h == nil {
            t.Fatal("handler is valid")
        }

        if p[0] != "123" {
            t.Fatal("param is valid")
        }

    })
}

func TestRouter_CreateController(t *testing.T) {
    App(true)
    var c *Container
    patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Get", func(_ *Container, _ string, ctx iface.IContext) reflect.Value {
        t.Log("mock Container.Get")
        m := &mockController{}
        m.SetContext(ctx)
        return reflect.ValueOf(m)
    })
    defer patches.Reset()

    router := NewRouter(nil)
    router.webHandlers = make(map[string]*Handler)
    router.cmdHandlers = make(map[string]*Handler)

    list := make(map[string]interface{})
    list["controller/MockController"] = map[string]int{"Index": 0}
    router.SetHandlers(ControllerWebPkg, list)

    context := &Context{}

    controller, action, _ := router.CreateController("/mock/index", context)
    if !controller.IsValid() {
        t.Fatal(`!controller.IsValid() `)
    }

    if context.ControllerId() != "/mock" {
        t.Fatal(`context.ControllerId() !="/mock"`)
    }

    if context.ActionId() != "Index" {
        t.Fatal(`context.ActionId() != "Index"`)
    }

    if action.Type().NumIn() != 0 {
        t.Fatal(`action.Type().NumIn() != 0`)
    }

}
