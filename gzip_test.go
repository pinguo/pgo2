package pgo2

import (
    "bytes"
    "net/http/httptest"
    "testing"

    "github.com/pinguo/pgo2/iface"
)

func TestNewGzip(t *testing.T) {
    var obj interface{}
    obj = NewGzip()
    if _, ok := obj.(iface.IPlugin); ok == false {
        t.FailNow()
    }
}

func TestGzip_HandleRequest(t *testing.T) {
    body := bytes.NewReader([]byte(""))
    t.Run("Accept-Encoding!=gzip", func(t *testing.T) {
        r := httptest.NewRequest("GET", "/view", body)
        w := httptest.NewRecorder()
        context := &Context{}
        context.HttpRW(true, r, w)

        gzip := NewGzip()
        gzip.HandleRequest(context)

        if w.Result().Header.Get("Content-Encoding") != "" {
            t.FailNow()
        }
    })

    t.Run("img", func(t *testing.T) {
        r := httptest.NewRequest("GET", "/view.png", body)
        r.Header.Set("Accept-Encoding", "gzip")
        w := httptest.NewRecorder()
        context := &Context{}
        context.HttpRW(true, r, w)

        gzip := NewGzip()
        gzip.HandleRequest(context)

        if w.Result().Header.Get("Content-Encoding") != "" {
            t.FailNow()
        }
    })

    t.Run("normal", func(t *testing.T) {
        r := httptest.NewRequest("GET", "/view.html", body)
        r.Header.Set("Accept-Encoding", "gzip")
        w := httptest.NewRecorder()
        context := &Context{}
        context.HttpRW(true, r, w)

        gzip := NewGzip()
        gzip.HandleRequest(context)

        str := "ss"

        //context.response.WriteString(str)
        context.Output().Write([]byte(str))

        if w.Result().Header.Get("Content-Encoding") != "gzip" {
            t.Fatal(` w.Result().Header.Get("Content-Encoding")!="gzip"`)
        }

        if len(w.Body.String()) < 1 {
            t.Fatal(`len(w.Body.String()) <1`)
        }

    })

}
