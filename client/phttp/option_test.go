package phttp

import (
    "testing"
    "time"
)

func TestOption_SetCookie(t *testing.T) {
    o := &Option{}
    o.SetCookie("name", "v")
    if len(o.Cookies) != 1 {
        t.FailNow()
    }
}

func TestOption_SetHeader(t *testing.T) {
    o := &Option{}
    o.SetHeader("name", "v")
    if o.Header.Get("name") != "v" {
        t.FailNow()
    }
}

func TestOption_SetTimeout(t *testing.T) {
    o := &Option{}
    o.SetTimeout(1 * time.Second)
    if o.Timeout != 1*time.Second {
        t.FailNow()
    }
}
