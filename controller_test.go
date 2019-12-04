package pgo2

import (
    "bytes"
    "net/http/httptest"
    "path/filepath"
    "reflect"
    "strings"
    "testing"

    "github.com/agiledragon/gomonkey"
    "github.com/pinguo/pgo2/logs"
    "github.com/pinguo/pgo2/render"
)

func TestController_GetBindInfo(t *testing.T) {
    //className := "github.com/pinguo/pgo2/mockController"
    mockC := &mockController{}
    iActions := mockC.GetBindInfo(mockC)
    actions, _ := iActions.(map[string]int)
    if _, has := actions["Index"]; has == false {
        t.Fatal("_,has:=actions[\"Index\"];has==false")
    }

    if _, has := actions["Info"]; has == false {
        t.Fatal("_,has:=actions[\"Info\"];has==false")
    }
}

type mockContext struct {
}

func (m *mockContext) Error(_ *Context, _ string, _ ...interface{}) {

}

func TestController_HandlePanic(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    var c *Context
    context := &Context{}
    context.Init("", "", App().Log())
    patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Error", func(_ *Context, _ string, _ ...interface{}) { t.Log("mock mock Context.Error") })
    defer patches.Reset()

    mockC := &Controller{}
    mockC.SetContext(context)
    patchesController := gomonkey.ApplyMethod(reflect.TypeOf(mockC), "Json", func(_ *Controller, data interface{}, status int, msg ...string) {
        t.Log("mock Controller.Json")
    })
    defer patchesController.Reset()

    mockC.HandlePanic("testerr")

}

func TestController_Json(t *testing.T) {

    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    method := "POST"
    path := "/test"
    body := bytes.NewReader([]byte(""))

    r := httptest.NewRequest(method, path, body)
    w := httptest.NewRecorder()

    context.HttpRW(true, r, w)

    mockC := &Controller{}
    mockC.SetContext(context)

    data := map[string]string{"name": "name1"}
    status := 201
    sStatus := "201"

    mockC.Json(data, status)

    if strings.Index(w.Body.String(), `"name":"name1"`) < 0 {
        t.Fatal("strings.Index(w.Body.String(), \"name\":\"name1\")<", 0)
    }

    if strings.Index(w.Body.String(), `"status":`+sStatus) < 0 {
        t.Fatal("strings.Index(w.Body.String(), \"status\":\"201\")<", 0)
    }

    if w.Result().Header.Get("Content-Type") != "application/json; charset=utf-8" {
        t.Fatal(`w.Result().Header.Get("Content-Type")!="application/json; charset=utf-8"`)
    }

}

func TestController_Jsonp(t *testing.T) {

    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    method := "POST"
    path := "/test"
    body := bytes.NewReader([]byte(""))

    r := httptest.NewRequest(method, path, body)
    w := httptest.NewRecorder()

    context.HttpRW(true, r, w)

    mockC := &Controller{}
    mockC.SetContext(context)

    data := map[string]string{"name": "name1"}
    status := 201
    sStatus := "201"
    callback := "callback1"

    mockC.Jsonp(callback, data, status)

    if strings.Index(w.Body.String(), callback+"(") < 0 {
        t.Fatal("strings.Index(w.Body.String(), `callback(`)<", 0)
    }

    if strings.Index(w.Body.String(), `"name":"name1"`) < 0 {
        t.Fatal("strings.Index(w.Body.String(), \"name\":\"name1\")<", 0)
    }

    if strings.Index(w.Body.String(), `"status":`+sStatus) < 0 {
        t.Fatal("strings.Index(w.Body.String(), \"status\":\"201\")<", 0)
    }

    if w.Result().Header.Get("Content-Type") != "application/json; charset=utf-8" {
        t.Fatal(`w.Result().Header.Get("Content-Type")!="application/json; charset=utf-8"`)
    }

}

func TestController_Data(t *testing.T) {

    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    method := "POST"
    path := "/test"
    body := bytes.NewReader([]byte(""))

    r := httptest.NewRequest(method, path, body)
    w := httptest.NewRecorder()

    context.HttpRW(true, r, w)

    mockC := &Controller{}
    mockC.SetContext(context)

    data := "data1"

    mockC.Data([]byte(data))

    if w.Body.String() != data {
        t.Fatal("w.Body.String() != ", data)
    }

    if w.Result().Header.Get("Content-Type") != "" {
        t.Fatal(`w.Result().Header.Get("Content-Type")!=""`)
    }

}

func TestController_Xml(t *testing.T) {

    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    method := "POST"
    path := "/test"
    body := bytes.NewReader([]byte(""))

    r := httptest.NewRequest(method, path, body)
    w := httptest.NewRecorder()

    context.HttpRW(true, r, w)

    mockC := &Controller{}
    mockC.SetContext(context)

    //data := map[string]string{"name": "name1"}
    type data struct {
        Id   string `xml:"id"`
        Name string `xml:"name"`
    }

    mockC.Xml(&data{"1", "name1"})

    if strings.Index(w.Body.String(), "<id>1</id>") < 0 {
        t.Fatal("strings.Index(w.Body.String(), \"<id>1</id>\")<", 0)
    }

    if strings.Index(w.Body.String(), "<name>name1</name>") < 0 {
        t.Fatal("strings.Index(w.Body.String(), \"<name>name1</name>\")<", 0)
    }

    if w.Result().Header.Get("Content-Type") != "application/xml; charset=utf-8" {
        t.Fatal(`w.Result().Header.Get("Content-Type")!="application/xml; charset=utf-8"`)
    }

}

//func TestController_ProtoBuf(t *testing.T) {
//    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
//    context := &Context{}
//    method := "POST"
//    path := "/test"
//    body := bytes.NewReader([]byte(""))
//
//    r := httptest.NewRequest(method, path, body)
//    w := httptest.NewRecorder()
//
//    context.HttpRW(true, r, w)
//
//    mockC := &Controller{}
//    mockC.SetContext(context)
//
//    data := map[string]string{"name": "name1"}
//
//    mockC.ProtoBuf(data)
//
//    fmt.Println(w.Body,w.Result())
//
//    if strings.Index(w.Body.String(), "<id>1</id>") < 0 {
//        t.Fatal("strings.Index(w.Body.String(), \"<id>1</id>\")<", 0)
//    }
//
//    if strings.Index(w.Body.String(), "<name>name1</name>") < 0 {
//        t.Fatal("strings.Index(w.Body.String(), \"<name>name1</name>\")<", 0)
//    }
//
//    if w.Result().Header.Get("Content-Type") != "application/xml; charset=utf-8" {
//        t.Fatal(`w.Result().Header.Get("Content-Type")!="application/xml; charset=utf-8"`)
//    }
//}

func TestController_Render(t *testing.T) {

    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    context := &Context{}
    method := "POST"
    path := "/test"
    body := bytes.NewReader([]byte(""))

    r := httptest.NewRequest(method, path, body)
    w := httptest.NewRecorder()

    context.HttpRW(true, r, w)

    mockC := &Controller{}
    mockC.SetContext(context)

    data := "data1"

    mockC.Render(render.NewData([]byte(data)))

    if w.Body.String() != data {
        t.Fatal("w.Body.String() != ", data)
    }

    if w.Result().Header.Get("Content-Type") != "" {
        t.Fatal(`w.Result().Header.Get("Content-Type")!=""`)
    }

}

func TestController_View(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    App().viewPath, _ = filepath.Abs("./test/data")
    SetAlias("@view", App().viewPath)

    context := &Context{}
    method := "POST"
    path := "/test"
    body := bytes.NewReader([]byte(""))

    r := httptest.NewRequest(method, path, body)
    w := httptest.NewRecorder()

    context.HttpRW(true, r, w)

    mockC := &Controller{}
    mockC.SetContext(context)

    mockC.View("view.html", nil)

    if strings.Index(w.Body.String(), "test view") < 0 {
        t.Fatal("strings.Index(w.Body.String(), \"test view\") < 0")
    }

    if w.Result().Header.Get("Content-Type") != "text/html; charset=utf-8" {
        t.Fatal(`w.Result().Header.Get("Content-Type")!="text/html; charset=utf-8"`)
    }

}
