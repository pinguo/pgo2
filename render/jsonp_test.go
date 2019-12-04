package render

import "testing"

func TestNewJsonp(t *testing.T) {
    var obj interface{}
    obj = NewJsonp("call", "aa")
    if _, ok := obj.(*Jsonp); ok == false {
        t.FailNow()
    }
}

func TestJsonp_Content(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()

        j := NewJsonp("call", TestJsonp_Content)
        j.Content()
    })

    t.Run("normal", func(t *testing.T) {
        j := NewJsonp("call", "aa")
        if j.Content() == nil {
            t.FailNow()
        }
    })

}

func TestJsonp_ContentType(t *testing.T) {
    j := NewJsonp("call", "aa")
    if j.ContentType() != "application/json; charset=utf-8" {
        t.FailNow()
    }
}

func TestJsonp_HttpCode(t *testing.T) {
    j := NewJsonp("call", "aa")
    if j.HttpCode() != 200 {
        t.FailNow()
    }
}

func TestJsonp_SetHttpCode(t *testing.T) {
    j := NewJsonp("call", "aa")
    j.SetHttpCode(100)
    if j.HttpCode() != 100 {
        t.FailNow()
    }
}
