package pgo2

import (
	"testing"

	"github.com/pinguo/pgo2/iface"
)

var initData = "initData"
func newMockObjectTestPool(iObj iface.IObject, params ...interface{}) iface.IObject {
	obj := iObj.(*mockObjectTest)
	obj.Data = initData
	return obj
}

var MockObjectTestClass = "github.com/pinguo/pgo2/mockObjectTest"

func newMockObjectTest(params ...interface{}) iface.IObject {
	return &mockObjectTest{}
}

type mockObjectTest struct {
	Object
	Data string
}

func TestObject_Context(t *testing.T) {
	obj := &Object{}
	ctr := &Context{}
	obj.SetContext(ctr)
	if _, ok := obj.Context().(iface.IContext); ok == false {
		t.FailNow()
	}
}

func TestObject_GetObj(t *testing.T) {
	obj := &Object{}
	ctr := &Context{}
	obj.SetContext(ctr)
	m := obj.GetObj(newMockObjectTest())
	if _, ok := m.(iface.IObject); ok == false {
		t.FailNow()
	}
}

func TestObject_GetObjSingle(t *testing.T) {
	obj := &Object{}
	ctr := &Context{}
	obj.SetContext(ctr)
	m := obj.GetObjSingle("mocktestobj", newMockObjectTest)
	if _, ok := m.(iface.IObject); ok == false {
		t.Fatal(`m.(iface.IObject) == false `)
	}

	data := "data111"
	mm := m.(*mockObjectTest)
	mm.Data = data
	m = obj.GetObjSingle("mocktestobj", newMockObjectTest)
	mm = m.(*mockObjectTest)
	if mm.Data != data {
		t.Fatal(`mm.Data != `, data)
	}

}

func TestObject_GetObjPool(t *testing.T) {
	App().Container().Bind(&mockObjectTest{})
	obj := &Object{}
	ctr := &Context{}
	obj.SetContext(ctr)

	m := obj.GetObjPool(MockObjectTestClass, newMockObjectTestPool)
	if _, ok := m.(iface.IObject); ok == false {
		t.FailNow()
	}

	data := "data111"
	mm := m.(*mockObjectTest)
	mm.Data = data
	m = obj.GetObjPool(MockObjectTestClass, newMockObjectTestPool)
	mm = m.(*mockObjectTest)
	if mm.Data != initData {
		t.Fatal(`mm.Data == `, initData)
	}
}

func TestObject_GetObjCtr(t *testing.T) {
	obj := &Object{}
	ctr := &Context{}
	ctr.actionId = "test"
	obj.SetContext(ctr)
	ctr1 := &Context{}
	ctr1.actionId = "test1"
	m := obj.GetObjCtx(ctr1, newMockObjectTest())
	mm := m.(*mockObjectTest)
	if _, ok := m.(iface.IObject); ok == false {
		t.Fatal(`_,ok:=m.(iface.IObject);ok==false`)
	}

	if mm.Context().ActionId() != "test1" {
		t.Fatal(`mm.Context().ActionId()!= "test1"`)
	}
}

func TestObject_GetObjSingleCtx(t *testing.T) {
	obj := &Object{}
	ctr := &Context{}
	ctr.actionId = "test"
	obj.SetContext(ctr)

	ctr1 := &Context{}
	ctr1.actionId = "test1"
	m := obj.GetObjSingleCtx(ctr1, "mocktestobj", newMockObjectTest)
	if _, ok := m.(iface.IObject); ok == false {
		t.Fatal(`m.(iface.IObject) == false `)
	}

	data := "data111"
	mm := m.(*mockObjectTest)
	mm.Data = data
	m = obj.GetObjSingleCtx(ctr1, "mocktestobj", newMockObjectTest)
	mm = m.(*mockObjectTest)
	if mm.Data != data {
		t.Fatal(`mm.Data != `, data)
	}

	if mm.Context().ActionId() != "test1" {
		t.Fatal(`mm.Context().ActionId()!= "test1"`)
	}

}

func TestObject_GetObjPoolCtx(t *testing.T) {
	App().Container().Bind(&mockObjectTest{})
	obj := &Object{}
	ctr := &Context{}
	ctr.actionId = "test"
	obj.SetContext(ctr)

	ctr1 := &Context{}
	ctr1.actionId = "test1"

	m := obj.GetObjPoolCtx(ctr1,MockObjectTestClass, newMockObjectTestPool)
	if _, ok := m.(iface.IObject); ok == false {
		t.FailNow()
	}

	ctr2 := &Context{}
	ctr2.actionId = "test1"
	data := "data111"
	mm := m.(*mockObjectTest)
	mm.Data = data
	m = obj.GetObjPoolCtx(ctr2, MockObjectTestClass, newMockObjectTestPool)
	mm = m.(*mockObjectTest)
	if mm.Data != initData {
		t.Fatal(`mm.Data == `, initData)
	}

	if mm.Context().ActionId() != "test1" {
		t.Fatal(`mm.Context().ActionId()!= "test1"`)
	}
}
