package core

import (
    "errors"
    "testing"
)

type mockTestConfigure struct {
    Name string
    id   string
    Desc string
}

func (m *mockTestConfigure) SetId(v string) {
    m.id = v
}

func (m *mockTestConfigure) SetDesc(v string) error {
    m.id = v
    return errors.New("test err")
}

func TestConfigure(t *testing.T) {
    t.Run("config==nil", func(t *testing.T) {
        mm := &mockTestConfigure{}
        Configure(mm, nil)
        if mm.Name != "" {
            t.FailNow()
        }
    })

    t.Run("notPointer", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        var dd string
        Configure(dd, map[string]interface{}{"name": "vvv"})

    })

    t.Run("normal", func(t *testing.T) {
        mm := &mockTestConfigure{}
        Configure(mm, map[string]interface{}{"name": "vvv", "id": "123"})
        if mm.Name != "vvv" {
            t.Fatal(`mm.Name != "vvv"`)
        }

        if mm.id != "123" {
            t.Fatal(`mm.id != "123"`)
        }
    })
}

func TestClientConfigure(t *testing.T) {
    t.Run("config==nil", func(t *testing.T) {
        mm := &mockTestConfigure{}
        ClientConfigure(mm, nil)
        if mm.Name != "" {
            t.FailNow()
        }
    })

    t.Run("notPointer", func(t *testing.T) {

        var dd string
        err := ClientConfigure(dd, map[string]interface{}{"name": "vvv"})
        if err == nil {
            t.FailNow()
        }

    })

    t.Run("method err", func(t *testing.T) {

        var dd string
        err := ClientConfigure(dd, map[string]interface{}{"desc": "vvv"})
        if err == nil {
            t.FailNow()
        }

    })

    t.Run("normal", func(t *testing.T) {
        mm := &mockTestConfigure{}
        ClientConfigure(mm, map[string]interface{}{"name": "vvv", "id": "123"})
        if mm.Name != "vvv" {
            t.Fatal(`mm.Name != "vvv"`)
        }

        if mm.id != "123" {
            t.Fatal(`mm.id != "123"`)
        }
    })
}
