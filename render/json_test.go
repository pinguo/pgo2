package render

import (
    "testing"
)

func TestNewJson(t *testing.T) {
    var obj interface{}
    obj = NewJson("aa")
    if _, ok := obj.(*Json); ok == false {
        t.FailNow()
    }
}

func TestJson_Content(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()

        j := NewJson(TestJson_Content)
        j.Content()
    })

    t.Run("normal", func(t *testing.T) {
        j := NewJson("aa")
        if j.Content() == nil {
            t.FailNow()
        }
    })

}

func TestJson_ContentType(t *testing.T) {
    j := NewJson("aa")
    if j.ContentType() != "application/json; charset=utf-8" {
        t.FailNow()
    }
}

func TestJson_HttpCode(t *testing.T) {
    j := NewJson("aa")
    if j.HttpCode() != 200 {
        t.FailNow()
    }
}

func TestJson_SetHttpCode(t *testing.T) {
    j := NewJson("aa")
    j.SetHttpCode(100)
    if j.HttpCode() != 100 {
        t.FailNow()
    }
}
