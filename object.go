package pgo2

import (
	"github.com/pinguo/pgo2/iface"
)

// Object base class of context based object
type Object struct {
	context iface.IContext
}

// GetContext get context of this object
func (o *Object) Context() iface.IContext {
	return o.context
}

// SetContext set context of this object
func (o *Object) SetContext(ctx iface.IContext) {
	o.context = ctx
}

// GetObject create new object
func (o *Object) GetObj(obj iface.IObject) iface.IObject {
	return o.GetObjCtx(o.Context(), obj)
}

// GetObjPool Get Object from pool
// Recommended:  Use GetObjBox instead.
func (o *Object) GetObjPool(className string, funcName iface.IObjPoolFunc, params ...interface{}) iface.IObject {
	return o.GetObjPoolCtx(o.Context(), className, funcName, params...)
}


// GetObjPool Get Object from box,Have the function of the pool
// params: Parameter passed into Prepare
func (o *Object) GetObjBox(className string, params ...interface{}) iface.IObject {
	return o.GetObjBoxCtx(o.Context(), className, params...)
}

// GetObject Get single object
func (o *Object) GetObjSingle(name string, funcName iface.IObjSingleFunc, params ...interface{}) iface.IObject {
	return o.GetObjSingleCtx(o.Context(), name, funcName, params...)
}

// GetObject create new object  and new Context
func (o *Object) GetObjCtx(ctx iface.IContext, obj iface.IObject) iface.IObject {
	obj.SetContext(ctx)
	return obj
}

// GetObjPoolCtx Get Object from pool and new Context
// Recommended: Use GetObjBoxCtx instead.
func (o *Object) GetObjPoolCtx(ctx iface.IContext, className string, funcName iface.IObjPoolFunc, params ...interface{}) iface.IObject {
	var obj iface.IObject
	if funcName != nil {
		obj = App().GetObjPool(className, ctx)
		return funcName(obj, params...)
	}else{
		obj = App().GetObjPool(className, ctx, params...)
	}
	return obj
}

// GetObjBoxCtx Get Object from box and new Context,Have the function of the pool
// params: Parameter passed into Prepare
func (o *Object) GetObjBoxCtx(ctx iface.IContext, className string, params ...interface{}) iface.IObject {
	return App().GetObjPool(className, ctx, params...)
}

// GetObject Get single object and new Context
func (o *Object) GetObjSingleCtx(ctx iface.IContext, name string, funcName iface.IObjSingleFunc, params ...interface{}) iface.IObject {
	// obj := funcName(params...)
	obj := App().GetObjSingle(name, funcName, params...)
	obj.SetContext(ctx)
	return obj
}
