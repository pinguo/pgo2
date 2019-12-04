package pgo2

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "net/url"
    "reflect"
    "strings"
    "testing"

    "github.com/agiledragon/gomonkey"
    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/logs"
    "github.com/pinguo/pgo2/perror"
)

func TestContext_Process(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})

    context := &Context{}

    context.enableAccessLog = true
    plugins := []iface.IPlugin{&mockPlugin{}}
    context.Process(plugins)
    if context.index != 1 {
        t.FailNow()
    }

}

func TestContext_finish(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}

    context.enableAccessLog = true
    plugins := []iface.IPlugin{&mockPlugin{}}

    t.Run("NextPanicPError", func(t *testing.T) {
        outStatus := 502
        context.SetOutput(&context.response)
        context.response.reset(httptest.NewRecorder())

        patches := gomonkey.ApplyMethod(reflect.TypeOf(context), "Next", func(_ *Context) {
            t.Log("mock Context.Next")
            panic(perror.New(outStatus, "testNextErr"))
        })
        defer patches.Reset()
        context.Process(plugins)
        if context.Status() != outStatus {
            t.FailNow()
        }

    })
}

func TestContext_FinishGoLog(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}

    context.Logger.Init(App().Name(), "test_logId", App().Log())
    context.enableAccessLog = true
    context.FinishGoLog()
}

func TestContext_Abort(t *testing.T) {
    context := &Context{}
    context.Abort()
    if context.index != MaxPlugins {
        t.FailNow()
    }
}

func TestContext_Copy(t *testing.T) {
    context := &Context{}
    plugins := []iface.IPlugin{&mockPlugin{}}
    context.Process(plugins)
    context.SetUserData("name", "v1")
    newContext := context.Copy().(*Context)
    if newContext.UserData("name", "") != "" {
        t.Fatal("newContext.UserData(\"name\", \"\")!=\"\"")
    }

    if newContext.plugins != nil {
        t.Fatal("newContext.plugins != nil")
    }

    if newContext.index != MaxPlugins {
        t.Fatal("newContext.index != MaxPlugins")
    }

    if newContext.objects != nil {
        t.Fatal("newContext.objects != nil")
    }
}

func TestContext_Getter(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    str := "ss"
    logId := "ffff"
    method := "POST"
    path := "/test"
    actionId := "testActionId"
    controllerId := "testControllerId"

    rCodeName := "rCode"
    rCode := "303"
    rICode := 303

    userDataName := "uData"
    userDataValue := "uValue"

    body := bytes.NewReader([]byte(str))

    r := httptest.NewRequest(method, path, body)
    r.Header.Set("X-Log-Id", logId)
    w := httptest.NewRecorder()
    w.Header().Set(rCodeName, rCode)
    rb := bytes.NewBuffer([]byte(str))
    w.Body = rb

    context.HttpRW(true, r, w)

    context.Output().WriteHeader(rICode)
    context.Output().Write([]byte(str))

    plugins := []iface.IPlugin{&mockPlugin{}}
    context.Process(plugins)

    context.SetHeader(rCodeName, rCode)

    context.SetControllerId(controllerId)
    context.SetActionId(actionId)
    context.SetUserData(userDataName, userDataValue)

    if context.Status() != rICode {
        t.Fatal("context.Status() != ", rCode)
    }

    if context.Path() != path {
        t.Fatal("context.Path() != ", path)
    }

    if context.Method() != method {
        t.Fatal("context.Method() != ", method)
    }

    if context.ControllerId() != controllerId {
        t.Fatal("context.ControllerId() != ", controllerId)
    }

    if context.ActionId() != actionId {
        t.Fatal("context.ActionId() != " + actionId)
    }

    if context.Size() != len(str) {
        t.Fatal("context.Size() != ", len(str))
    }

    if context.LogId() != logId {
        t.Fatal("context.LogId() != ", logId)
    }

    if context.UserData(userDataName, "") != userDataValue {
        t.Fatal("context.UserData(userDataName,\"\") != ", userDataValue)
    }

    if context.Output().Header().Get(rCodeName) != rCode {
        t.Fatal("context.Header(rCodeName,\"\") != ", rCode)
    }
}

func TestContext_ClientIp(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    str := "ss"
    method := "POST"
    path := "/test"

    body := bytes.NewReader([]byte(str))
    w := httptest.NewRecorder()

    t.Run("X-Forwarded-For_more", func(t *testing.T) {
        r := httptest.NewRequest(method, path, body)
        ip := "129.1.1.1"
        r.Header.Set("X-Forwarded-For", ip+",")
        context.HttpRW(true, r, w)
        if context.ClientIp() != ip {
            t.FailNow()
        }
    })

    t.Run("X-Forwarded-For", func(t *testing.T) {
        r := httptest.NewRequest(method, path, body)
        ip := "129.1.1.1"
        r.Header.Set("X-Forwarded-For", ip)
        context.HttpRW(true, r, w)
        if context.ClientIp() != ip {
            t.FailNow()
        }
    })

    t.Run("X-Client-Ip", func(t *testing.T) {
        r := httptest.NewRequest(method, path, body)
        ip := "129.1.1.1"
        r.Header.Set("X-Client-Ip", ip)
        context.HttpRW(true, r, w)
        if context.ClientIp() != ip {
            t.FailNow()
        }
    })

    t.Run("X-Real-Ip", func(t *testing.T) {
        r := httptest.NewRequest(method, path, body)
        ip := "129.1.1.1"
        r.Header.Set("X-Real-Ip", ip)
        context.HttpRW(true, r, w)
        if context.ClientIp() != ip {
            t.FailNow()
        }
    })

    t.Run("RemoteAddr", func(t *testing.T) {
        r := httptest.NewRequest(method, path, body)
        ip := "129.1.1.1"
        r.RemoteAddr = ip
        context.HttpRW(true, r, w)
        if context.ClientIp() != ip {
            t.FailNow()
        }
    })

}

func TestContext_Cookie(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    str := "ss"
    method := "POST"
    path := "/test"

    cName := "cname"
    cNameV := "cname"

    body := bytes.NewReader([]byte(str))
    w := httptest.NewRecorder()
    r := httptest.NewRequest(method, path, body)
    cookieUser := &http.Cookie{Name: cName, Value: cNameV, Path: "/"}
    r.AddCookie(cookieUser)
    cookieUser1 := &http.Cookie{Name: cName + "1", Value: cNameV + "1", Path: "/"}
    r.AddCookie(cookieUser1)

    context.HttpRW(true, r, w)
    t.Run("GetCookie", func(t *testing.T) {
        if context.Cookie(cName, "") != cNameV {
            t.FailNow()
        }
    })

    t.Run("CookieAll", func(t *testing.T) {
        if len(context.CookieAll()) != len(r.Cookies()) {
            t.FailNow()
        }
    })

    t.Run("SetCookie", func(t *testing.T) {
        cName2 := cName + "2"
        cNameV2 := cNameV + "2"
        cookieUser := &http.Cookie{Name: cName2, Value: cNameV2, Path: "/"}
        context.SetCookie(cookieUser)
        strCookie := context.Output().Header().Get("Set-Cookie")
        if strings.Index(strCookie, cName2) == -1 || strings.Index(strCookie, cNameV2) == -1 {
            t.FailNow()
        }
    })
}

func TestContext_Header(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    str := "ss"
    method := "POST"
    path := "/test"

    cName := "cname"
    cNameV := "cname"

    body := bytes.NewReader([]byte(str))
    w := httptest.NewRecorder()
    r := httptest.NewRequest(method, path, body)
    r.Header.Set(cName, cNameV)
    r.Header.Set(cName+"1", cNameV+"1")

    context.HttpRW(true, r, w)
    t.Run("GetHeader", func(t *testing.T) {
        if context.Header(cName, "") != cNameV {
            t.FailNow()
        }
    })

    t.Run("GetAll", func(t *testing.T) {
        if len(context.HeaderAll()) != 2 {
            t.FailNow()
        }
    })

    t.Run("SetHeader", func(t *testing.T) {
        cName2 := cName + "2"
        cNameV2 := cNameV + "2"
        context.SetHeader(cName2, cNameV2)
        if context.Output().Header().Get(cName2) != cNameV2 {
            t.FailNow()
        }
    })
}

func TestContext_Query(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    str := "ss"
    method := "GET"
    param1 := "name1"
    param2 := "name2"
    paramV1 := "nameV1"
    paramV2 := "nameV2"

    var uri = &url.Values{}
    uri.Add(param1, paramV1)
    uri.Add(param2, paramV2)
    path := "/test?" + uri.Encode()

    body := bytes.NewReader([]byte(str))
    w := httptest.NewRecorder()
    r := httptest.NewRequest(method, path, body)
    context.HttpRW(true, r, w)
    t.Run("Query", func(t *testing.T) {
        if context.Query(param1, "") != paramV1 {
            t.FailNow()
        }
    })

    t.Run("", func(t *testing.T) {
        if len(context.QueryAll()) != 2 {
            t.FailNow()
        }
    })

}

func TestContext_QueryArray(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    str := "ss"
    method := "GET"
    param1 := "name1[]"
    param2 := "name1[]"
    paramV1 := "nameV1"
    paramV2 := "nameV2"

    var uri = &url.Values{}
    uri.Add(param1, paramV1)
    uri.Add(param2, paramV2)
    path := "/test?" + uri.Encode()

    body := bytes.NewReader([]byte(str))
    w := httptest.NewRecorder()
    r := httptest.NewRequest(method, path, body)
    context.HttpRW(true, r, w)
    if len(context.QueryArray(param1)) != 2 {
        t.FailNow()
    }

}

func TestContext_Post(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    method := "POST"
    param1 := "name1"
    param2 := "name2"
    paramV1 := "nameV1"
    paramV2 := "nameV2"

    var uri = &url.Values{}
    uri.Add(param1, paramV1)
    uri.Add(param2, paramV2)
    path := "/test"

    body := bytes.NewReader([]byte("&" + uri.Encode()))
    w := httptest.NewRecorder()

    r := httptest.NewRequest(method, path, body)
    r.Header.Set("Content-type", "application/x-www-form-urlencoded")

    context.HttpRW(true, r, w)

    t.Run("Post", func(t *testing.T) {
        if context.Post(param1, "") != paramV1 {
            t.FailNow()
        }
    })

    t.Run("PostAll", func(t *testing.T) {
        if len(context.PostAll()) != 2 {
            t.FailNow()
        }
    })

}

func TestContext_PostArray(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    method := "POST"
    param1 := "name1"
    param2 := "name1"
    paramV1 := "nameV1"
    paramV2 := "nameV2"

    var uri = &url.Values{}
    uri.Add(param1, paramV1)
    uri.Add(param2, paramV2)

    path := "/test"

    body := bytes.NewReader([]byte(uri.Encode()))
    w := httptest.NewRecorder()

    r := httptest.NewRequest(method, path, body)
    r.Header.Set("Content-type", "application/x-www-form-urlencoded")

    context.HttpRW(true, r, w)

    if len(context.PostArray(param1)) != 2 {
        t.FailNow()
    }

}

func TestContext_Param(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    method := "POST"
    param1 := "name1"
    param2 := "name2"
    paramV1 := "nameV1"
    paramV2 := "nameV2"

    var uri = &url.Values{}
    uri.Add(param1, paramV1)
    uri.Add(param2, paramV2)

    param1_1 := "name11"
    param2_1 := "name12"
    paramV1_1 := "nameV11"
    paramV2_1 := "nameV21"

    var post = &url.Values{}
    post.Add(param1_1, paramV1_1)
    post.Add(param2_1, paramV2_1)

    path := "/test?" + uri.Encode()

    body := bytes.NewReader([]byte(post.Encode()))
    w := httptest.NewRecorder()

    r := httptest.NewRequest(method, path, body)
    r.Header.Set("Content-type", "application/x-www-form-urlencoded")

    context.HttpRW(true, r, w)

    if context.Param(param1, "") != paramV1 {
        t.Fatal("context.Param(param1, \"\") != ", paramV1)
    }

    if context.Param(param1_1, "") != paramV1_1 {
        t.Fatal("context.Param(param1_1, \"\") != ", paramV1_1)
    }

    if len(context.ParamAll()) != 4 {
        t.Fatal("len(context.ParamAll()) !=", 4)
    }

}

func TestContext_ParamArray(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    method := "POST"
    param1 := "name1"
    param2 := "name1"
    paramV1 := "nameV1"
    paramV2 := "nameV2"

    var uri = &url.Values{}
    uri.Add(param1, paramV1)
    uri.Add(param2, paramV2)

    param1_1 := "name11"
    param2_1 := "name11"
    paramV1_1 := "nameV11"
    paramV2_1 := "nameV21"

    var post = &url.Values{}
    post.Add(param1_1, paramV1_1)
    post.Add(param2_1, paramV2_1)

    path := "/test?" + uri.Encode()

    body := bytes.NewReader([]byte(post.Encode()))
    w := httptest.NewRecorder()

    r := httptest.NewRequest(method, path, body)
    r.Header.Set("Content-type", "application/x-www-form-urlencoded")

    context.HttpRW(true, r, w)

    if len(context.ParamArray(param1)) != 2 {
        t.Fatal("len(context.ParamArray(param1)) !=", 2)
    }

}
