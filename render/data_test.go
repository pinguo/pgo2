package render

import (
    "bytes"
    "testing"
)

func TestNewData(t *testing.T) {
    var obj interface{}
    obj = NewData([]byte("aa"))
    if _, ok := obj.(*Data); ok == false {
        t.FailNow()
    }
}

func TestData_SetHttpCode(t *testing.T) {
    d := NewData([]byte("aa"))
    d.SetHttpCode(100)
    if d.HttpCode() != 100 {
        t.FailNow()
    }
}

func TestData_HttpCode(t *testing.T) {
    d := NewData([]byte("aa"))
    if d.HttpCode() != 200 {
        t.FailNow()
    }
}

func TestData_Content(t *testing.T) {
    d := NewData([]byte("aa"))
    if bytes.Equal(d.Content(), []byte("aa")) == false {
        t.FailNow()
    }
}

func TestData_ContentType(t *testing.T) {
    d := NewData([]byte("aa"))
    if d.ContentType() != "" {
        t.FailNow()
    }
}
