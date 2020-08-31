package pgo2

import (
	"fmt"
	"reflect"
	"testing"
)

type containerTestCommand struct {
	Controller
}

func (c *Container) Prepare(a, b string){
	fmt.Println("a",a,"b",b)
}

func TestContainer(t *testing.T) {
	container := NewContainer("on")
	className := container.Bind(&containerTestCommand{})
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
		container.Get(className, &Context{}).Interface()
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
		rt := reflect.TypeOf(containerTestCommand{})

		container.Put(className, reflect.New(rt))
	})

	t.Run("GetPrepareNoParams", func(t *testing.T) {
		IC := container.Get(className, &Context{}).Interface()
		if _, ok := IC.(*containerTestCommand); ok == false {
			t.FailNow()
		}
	})

	//
	t.Run("GetPrepareNoParams", func(t *testing.T) {
		IC1 := container.Get(className, &Context{}, "aa","bb").Interface()
		if _, ok := IC1.(*containerTestCommand); ok == false {
			t.FailNow()
		}
	})

}
