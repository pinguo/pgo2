package pgo2

import (
    "testing"
)

func TestNewStopBefore(t *testing.T) {
    var obj interface{}

    obj = NewStopBefore()

    if _, ok := obj.(*StopBefore); ok == false {
        t.FailNow()
    }
}

type mockTestStopBefore struct {
}

func (m *mockTestStopBefore) Index(name chan string, v string) {
    name <- v
}

func TestStopBefore_Add(t *testing.T) {

    t.Run("queue.len>10", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        obj := NewStopBefore()
        for i := 0; i < 10000; i++ {
            name := make(chan string)
            obj.Add(&mockTestStopBefore{}, "Index", name, "t")
        }
    })

    t.Run("action valid", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        name := make(chan string)
        obj := NewStopBefore()
        obj.Add(&mockTestStopBefore{}, "Index1", name, "t")
    })

    t.Run("normal", func(t *testing.T) {
        name := make(chan string)
        obj := NewStopBefore()
        obj.Add(&mockTestStopBefore{}, "Index", name, "t")
        if len(obj.performQueue) != 1 {
            t.FailNow()
        }
    })

}

func TestStopBefore_Exec(t *testing.T) {
    name := make(chan string)
    v := "test"
    obj := NewStopBefore()
    obj.Add(&mockTestStopBefore{}, "Index", name, v)
    go obj.Exec()
    if n, ok := <-name; ok == false || n != v {
        t.FailNow()
    }

}
