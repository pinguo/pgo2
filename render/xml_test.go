package render

import "testing"

func TestNewXml(t *testing.T) {
    var obj interface{}
    obj = NewXml("aa")
    if _, ok := obj.(*Xml); ok == false {
        t.FailNow()
    }
}

func TestXml_Content(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()

        j := NewXml(TestXml_Content)
        j.Content()
    })

    t.Run("normal", func(t *testing.T) {
        j := NewXml("aa")
        if j.Content() == nil {
            t.FailNow()
        }
    })

}

func TestXml_ContentType(t *testing.T) {
    j := NewXml("aa")
    if j.ContentType() != "application/xml; charset=utf-8" {
        t.FailNow()
    }
}

func TestXml_HttpCode(t *testing.T) {
    j := NewXml("aa")
    if j.HttpCode() != 200 {
        t.FailNow()
    }
}

func TestXml_SetHttpCode(t *testing.T) {
    j := NewXml("aa")
    j.SetHttpCode(100)
    if j.HttpCode() != 100 {
        t.FailNow()
    }
}
