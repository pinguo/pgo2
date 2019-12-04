package pgo2

import (
    "reflect"
    "testing"
)

type containerTestCommand struct {
    Controller
}

func TestContainer(t *testing.T) {
    container := NewContainer("on")
    container.Bind(&containerTestCommand{})
    className := "github.com/pinguo/pgo2/containerTestCommand"
    t.Run("Bind&Has", func(t *testing.T) {
        if container.Has(className) == false {
            t.FailNow()
        }
    })

    t.Run("GetInfo", func(t *testing.T) {
        if info := container.GetInfo(className); info == nil {
            t.FailNow()
        }
    })

    t.Run("GetType", func(t *testing.T) {
        if container.GetType(className).String() != "pgo2.containerTestCommand" {
            t.FailNow()
        }
    })

    t.Run("Get", func(t *testing.T) {
        IC := container.Get(className, &Context{}).Interface()
        if _, ok := IC.(*containerTestCommand); ok == false {
            t.FailNow()
        }
    })

    t.Run("PathList", func(t *testing.T) {

        if len(container.PathList("github.com/pinguo/pgo2", "Command")) < 1 {
            t.FailNow()
        }
    })

    t.Run("Put", func(t *testing.T) {
        rt := reflect.TypeOf(&containerTestCommand{})

        container.Put(className, reflect.New(rt))
    })

}
