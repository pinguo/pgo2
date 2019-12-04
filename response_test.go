package pgo2

import (
    "bytes"
    "net/http/httptest"
    "testing"
)

func TestResponse_WriteHeader(t *testing.T) {
    code := 201
    w := httptest.NewRecorder()
    r := &Response{}
    r.reset(w)
    r.WriteHeader(code)
    if r.Status() != code {
        t.FailNow()
    }
}

func TestResponse_Write(t *testing.T) {
    str := "ss"
    w := httptest.NewRecorder()
    r := &Response{}
    r.reset(w)
    r.Write([]byte(str))

    if r.Size() != len(str) {
        t.FailNow()
    }

    if len(w.Body.String()) != len(str) {
        t.FailNow()
    }
}

func TestResponse_WriteString(t *testing.T) {
    str := "ss"
    w := httptest.NewRecorder()
    r := &Response{}
    r.reset(w)
    r.WriteString(str)

    if r.Size() != len(str) {
        t.FailNow()
    }

    if len(w.Body.String()) != len(str) {
        t.FailNow()
    }
}

func TestResponse_ReadFrom(t *testing.T) {
    str := "ss"
    var buf bytes.Buffer
    buf.WriteString(str)

    read := bytes.NewReader(buf.Bytes())

    w := httptest.NewRecorder()
    r := &Response{}
    r.reset(w)

    r.ReadFrom(read)

    if r.Size() != len(str) {
        t.FailNow()
    }

    if len(w.Body.String()) != len(str) {
        t.FailNow()
    }

}
