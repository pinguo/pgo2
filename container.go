package pgo2

import (
	"reflect"
	"strings"
	"sync"

	"github.com/pinguo/pgo2/iface"
)

type bindItem struct {
	pool  sync.Pool     // object pool
	info  interface{}   // binding info
	zero  reflect.Value // zero value
	cmIdx int           // construct index
	pmIdx int           // prepare index
}

const (
	EnablePoolOn = "on"
	// EnablePoolOff = "off"
)

// Container the container component, configuration:
// container:
//     enablePool: on/off
func NewContainer(enable string) *Container {
	if enable == "" {
		enable = EnablePoolOn
	}

	return &Container{
		items:      make(map[string]*bindItem),
		enablePool: enable,
	}
}

type Container struct {
	enablePool string
	items      map[string]*bindItem
}

// Bind bind template object to class,
// param i must be a pointer of struct.
func (c *Container) Bind(i interface{}) string{
	iv := reflect.ValueOf(i)
	if iv.Kind() != reflect.Ptr {
		panic("Container: invalid type, need pointer")
	}

	// initialize binding
	rt := iv.Elem().Type()
	item := bindItem{zero: reflect.Zero(rt), cmIdx: -1, pmIdx: -1}
	item.pool.New = func() interface{} { return reflect.New(rt) }

	// get binding info
	if bind, ok := i.(iface.IBind); ok {
		item.info = bind.GetBindInfo(i)
	}

	// get method index
	it := iv.Type()
	nm := it.NumMethod()
	for i := 0; i < nm; i++ {
		switch it.Method(i).Name {
		case ConstructMethod:
			item.cmIdx = i
		case PrepareMethod:
			item.pmIdx = i
		}
	}

	// get class name
	pkgPath := rt.PkgPath()

	if index := strings.Index(pkgPath, "/"+ControllerWebPkg); index >= 0 {
		pkgPath = pkgPath[index+1:]
	}

	if index := strings.Index(pkgPath, "/"+ControllerCmdPkg); index >= 0 {
		pkgPath = pkgPath[index+1:]
	}

	name := pkgPath + "/" + rt.Name()

	if len(name) > VendorLength && name[:VendorLength] == VendorPrefix {
		name = name[VendorLength:]
	}

	c.items[name] = &item

	return name
}

// Has check if the class exists in container
func (c *Container) Has(name string) bool {
	_, ok := c.items[name]
	return ok
}

// GetInfo get class binding info
func (c *Container) GetInfo(name string) interface{} {
	if item, ok := c.items[name]; ok {
		return item.info
	}

	panic("Container: class not found, " + name)
}

// GetType get class reflect type
func (c *Container) GetType(name string) reflect.Type {
	if item, ok := c.items[name]; ok {
		return item.zero.Type()
	}

	panic("Container: class not found, " + name)
}

// Get get new class object. name is class name, config is properties map,
// params is optional construct parameters.
func (c *Container) Get(name string, ctx iface.IContext, params ...interface{}) reflect.Value {
	item, ok := c.items[name]
	if !ok {
		panic("Container: class not found, " + name)
	}

	// get new object from pool
	rv := item.pool.Get().(reflect.Value)
	if c.enablePool == EnablePoolOn {
		// reset properties
		rv.Elem().Set(item.zero)
		ctx.Cache(name, rv)
	}

	if obj, ok := rv.Interface().(iface.IObject); ok {
		// inject context
		obj.SetContext(ctx)
	}

	// call Init()
	if item.pmIdx != -1 {
		if im := rv.Method(item.pmIdx); im.IsValid() {
			in := make([]reflect.Value, len(params))
			for k, arg := range params {
				in[k] = reflect.ValueOf(arg)
			}
			im.Call(in)
		}
	}

	return rv
}

// Put put back reflect value to object pool
func (c *Container) Put(name string, rv reflect.Value) {
	if item, ok := c.items[name]; ok {
		item.pool.Put(rv)
		return
	}

	panic("Container: class not found, " + name)
}

// PathList Gets a list of paths with the specified path prefix
func (c *Container) PathList(prefix, suffix string) map[string]interface{} {
	list := make(map[string]interface{})
	for k, item := range c.items {
		//if strings.Index(k, prefix) == 0 && strings.Index(k, suffix) > 0 {
		//	list[k] = item.info
		//}
		if strings.Index(k, prefix) == 0 {
			list[k] = item.info
		}
	}
	return list
}
