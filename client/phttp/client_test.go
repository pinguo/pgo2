package phttp

import (
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "net/url"
    "reflect"
    "strings"
    "testing"
    "time"

    "github.com/agiledragon/gomonkey"
)

func TestNew(t *testing.T) {
    c, err := New(nil)
    if _, ok := c.(*Client); ok == false {
        t.FailNow()
    }

    if err != nil {
        t.FailNow()
    }
}

func TestClient_SetTimeout(t *testing.T) {
    ic, _ := New(nil)
    c, _ := ic.(*Client)
    t.Run("err", func(t *testing.T) {
        if c.SetTimeout("1ss") == nil {
            t.FailNow()
        }
    })

    t.Run("normal", func(t *testing.T) {
        if c.SetTimeout("1s") != nil {
            t.FailNow()
        }
    })

}

func TestClient_SetVerifyPeer(t *testing.T) {
    ic, _ := New(nil)
    c, _ := ic.(*Client)
    c.SetVerifyPeer(true)
    if c.verifyPeer != true {
        t.FailNow()
    }
}

func TestClient_SetUserAgent(t *testing.T) {
    ic, _ := New(nil)
    c, _ := ic.(*Client)
    c.SetUserAgent("testUA")
    if c.userAgent != "testUA" {
        t.FailNow()
    }
}

func TestClient_Init(t *testing.T) {

}

func TestClient_Get(t *testing.T) {

    ic, _ := New(nil)
    c, _ := ic.(*Client)

    t.Run("data=nil", func(t *testing.T) {
        patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *Client, _ *http.Request, _ ...*Option) (*http.Response, error) {
            t.Log("mock Client.Do")
            return httptest.NewRecorder().Result(), nil
        })
        defer patches.Reset()
        _, err := c.Get("http://127.0.0.1", nil)
        if err != nil {
            t.FailNow()
        }

    })

    t.Run("data=url.values", func(t *testing.T) {
        data := url.Values{}
        data.Set("k", "v")

        patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *Client, _ *http.Request, _ ...*Option) (*http.Response, error) {
            t.Log("mock Client.Do")
            return httptest.NewRecorder().Result(), nil
        })
        defer patches.Reset()

        _, err := c.Get("http://127.0.0.1?id=1", data)
        if err != nil {
            t.FailNow()
        }
    })

    t.Run("data=map", func(t *testing.T) {
        data := make(map[string]interface{})
        data["k"] = "v"

        patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *Client, _ *http.Request, _ ...*Option) (*http.Response, error) {
            t.Log("mock Client.Do")
            return httptest.NewRecorder().Result(), nil
        })
        defer patches.Reset()

        _, err := c.Get("http://127.0.0.1", data)
        if err != nil {
            t.FailNow()
        }
    })

    t.Run("other", func(t *testing.T) {
        _, err := c.Get("http://127.0.0.1", 1)
        if err == nil {
            t.FailNow()
        }
    })

    t.Run("return=nil", func(t *testing.T) {
        patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "HttpNewRequest", func(_ *Client, _, _ string, _ io.Reader) (*http.Request, error) {
            t.Log("mock Client.HttpNewRequest")
            return nil, fmt.Errorf("errmock")
        })
        defer patches.Reset()

        _, err := c.Get("http://127.0.0.1", nil)
        if err == nil {
            t.FailNow()
        }

    })
}

func TestClient_Post(t *testing.T) {
    ic, _ := New(nil)
    c, _ := ic.(*Client)

    t.Run("data=nil", func(t *testing.T) {
        patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *Client, _ *http.Request, _ ...*Option) (*http.Response, error) {
            t.Log("mock Client.Do")
            return httptest.NewRecorder().Result(), nil
        })
        defer patches.Reset()
        _, err := c.Post("http://127.0.0.1", nil)
        if err != nil {
            t.FailNow()
        }

    })

    t.Run("data=url.values", func(t *testing.T) {
        data := url.Values{}
        data.Set("k", "v")

        patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *Client, _ *http.Request, _ ...*Option) (*http.Response, error) {
            t.Log("mock Client.Do")
            return httptest.NewRecorder().Result(), nil
        })
        defer patches.Reset()

        _, err := c.Post("http://127.0.0.1?id=1", data)
        if err != nil {
            t.FailNow()
        }
    })

    t.Run("data=map", func(t *testing.T) {
        data := make(map[string]interface{})
        data["k"] = "v"

        patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *Client, _ *http.Request, _ ...*Option) (*http.Response, error) {
            t.Log("mock Client.Do")
            return httptest.NewRecorder().Result(), nil
        })
        defer patches.Reset()

        _, err := c.Post("http://127.0.0.1", data)
        if err != nil {
            t.FailNow()
        }
    })

    t.Run("data=string", func(t *testing.T) {
        patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *Client, _ *http.Request, _ ...*Option) (*http.Response, error) {
            t.Log("mock Client.Do")
            return httptest.NewRecorder().Result(), nil
        })
        defer patches.Reset()

        _, err := c.Post("http://127.0.0.1", "aaa")
        if err != nil {
            t.FailNow()
        }
    })

    t.Run("data=[]byte", func(t *testing.T) {
        patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *Client, _ *http.Request, _ ...*Option) (*http.Response, error) {
            t.Log("mock Client.Do")
            return httptest.NewRecorder().Result(), nil
        })
        defer patches.Reset()

        _, err := c.Post("http://127.0.0.1", []byte("aaa"))
        if err != nil {
            t.FailNow()
        }
    })

    t.Run("data=io.Reader", func(t *testing.T) {
        patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *Client, _ *http.Request, _ ...*Option) (*http.Response, error) {
            t.Log("mock Client.Do")
            return httptest.NewRecorder().Result(), nil
        })
        defer patches.Reset()

        _, err := c.Post("http://127.0.0.1", strings.NewReader("aaa"))
        if err != nil {
            t.FailNow()
        }
    })

    t.Run("other", func(t *testing.T) {
        _, err := c.Post("http://127.0.0.1", 1)
        if err == nil {
            t.FailNow()
        }
    })

    t.Run("return=nil", func(t *testing.T) {
        patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "HttpNewRequest", func(_ *Client, _, _ string, _ io.Reader) (*http.Request, error) {
            t.Log("mock Client.HttpNewRequest")
            return nil, fmt.Errorf("errmock")
        })
        defer patches.Reset()

        _, err := c.Post("http://127.0.0.1", nil)
        if err == nil {
            t.FailNow()
        }

    })
}

func TestClient_Do(t *testing.T) {
    config := map[string]interface{}{"userAgent": "test userAgent"}
    ic, _ := New(config)
    c, _ := ic.(*Client)

    patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *Client, _ *http.Request, _ ...*Option) (*http.Response, error) {
        t.Log("mock Client.Do")
        return httptest.NewRecorder().Result(), nil
    })
    defer patches.Reset()

    option := &Option{}
    option.SetTimeout(1 * time.Second)
    option.SetHeader("Content-Type", "test")
    option.SetCookie("name", "v1")
    req := httptest.NewRequest("GET", "http://127.0.0.1", nil)

    _, err := c.Do(req, option)
    if err != nil {
        t.FailNow()
    }

}
