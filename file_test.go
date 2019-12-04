package pgo2

import (
    "bytes"
    "net/http/httptest"
    "path/filepath"
    "strings"
    "testing"

    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/logs"
)

func TestNewFile(t *testing.T) {
    var obj interface{}
    obj = NewFile(nil)
    if _, ok := obj.(iface.IPlugin); ok == false {
        t.FailNow()
    }
}

func TestFile_SetExcludeExtensions(t *testing.T) {
    file := NewFile(nil)
    var Extensions []interface{}
    Extensions = []interface{}{".html", ".htm"}
    file.SetExcludeExtensions(Extensions)

    if len(file.excludeExtensions) != len(Extensions) {
        t.FailNow()
    }
}

func TestFile_HandleRequest(t *testing.T) {
    App(true).Log().SetTarget(logs.TargetConsole, &mockTarget{})
    App().publicPath, _ = filepath.Abs("./test/data")
    SetAlias("@public", App().publicPath)
    method := "GET"

    body := bytes.NewReader([]byte(""))

    t.Run("pathEmpty", func(t *testing.T) {

        r := httptest.NewRequest(method, "/view", body)
        w := httptest.NewRecorder()
        context := &Context{}
        context.HttpRW(true, r, w)

        file := NewFile(nil)
        file.HandleRequest(context)

        if w.Body.String() != "" {
            t.Fatal(`w.Body.String() !=""`)
        }

        if w.Result().Header.Get("Content-Type") != "" {
            t.Fatal(`w.Result().Header.Get("Content-Type")!=""`)
        }
    })

    t.Run("excludeExtensions", func(t *testing.T) {
        r := httptest.NewRequest(method, "/view.html", body)
        w := httptest.NewRecorder()
        context := &Context{}
        context.HttpRW(true, r, w)
        file := NewFile(nil)
        file.SetExcludeExtensions([]interface{}{".html"})
        file.HandleRequest(context)

        if w.Body.String() != "" {
            t.Fatal(`w.Body.String() !=""`)
        }

        if w.Result().Header.Get("Content-Type") != "" {
            t.Fatal(`w.Result().Header.Get("Content-Type")!=""`)
        }
    })

    t.Run("methodErr", func(t *testing.T) {
        r := httptest.NewRequest("POST", "/view.html", body)
        w := httptest.NewRecorder()
        context := &Context{}
        context.HttpRW(true, r, w)
        file := NewFile(nil)
        file.HandleRequest(context)

        if w.Code != 405 {
            t.Fatal(`w.Code != 405`)
        }
    })

    t.Run("pathNotExist", func(t *testing.T) {

        r := httptest.NewRequest(method, "/viewNotExists.html", body)
        w := httptest.NewRecorder()
        context := &Context{}
        context.HttpRW(true, r, w)

        file := NewFile(nil)
        file.HandleRequest(context)

        if w.Code != 404 {
            t.Fatal(` w.Code !=404`)
        }

    })

    t.Run("ServeContent", func(t *testing.T) {
        r := httptest.NewRequest(method, "/view.html", body)
        w := httptest.NewRecorder()
        context := &Context{}
        context.HttpRW(true, r, w)

        file := NewFile(nil)
        file.HandleRequest(context)

        if strings.Index(w.Body.String(), "test view") < 0 {
            t.Fatal("strings.Index(w.Body.String(), \"test view\") < 0")
        }

        if w.Result().Header.Get("Content-Type") != "text/html; charset=utf-8" {
            t.Fatal(`w.Result().Header.Get("Content-Type")!="text/html; charset=utf-8"`)
        }
    })
}
