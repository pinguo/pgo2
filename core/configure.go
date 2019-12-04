package core

import (
    "errors"
    "reflect"
    "strings"
)

// Configure configure object using the given configuration,
// obj is a pointer or reflect.Value of a pointer,
// config is the configuration map for properties.
func Configure(obj interface{}, config map[string]interface{}) {
    // skip empty configuration
    if config == nil || len(config) == 0 {
        return
    }
    // v refer to the object pointer
    var v reflect.Value
    if _, ok := obj.(reflect.Value); ok {
        v = obj.(reflect.Value)
    } else {
        v = reflect.ValueOf(obj)
    }

    if v.Kind() != reflect.Ptr {
        panic("Configure: obj require a pointer or reflect.Value of a pointer")
    }

    // rv refer to the value of pointer
    rv := v.Elem()

    for key, val := range config {
        // change key to title string
        key = strings.Title(key)

        // check object's setter method
        if method := v.MethodByName("Set" + key); method.IsValid() {
            newVal := reflect.ValueOf(val).Convert(method.Type().In(0))
            method.Call([]reflect.Value{newVal})
            continue
        }

        // check object's public field
        field := rv.FieldByName(key)
        if field.IsValid() && field.CanSet() {
            newVal := reflect.ValueOf(val).Convert(field.Type())
            field.Set(newVal)
            continue
        }
    }
}

// ClientConfigure configure object using the given configuration,
// obj is a pointer or reflect.Value of a pointer,
// config is the configuration map for properties.
func ClientConfigure(obj interface{}, config map[string]interface{}) error {
    // skip empty configuration
    if config == nil || len(config) == 0 {
        return nil
    }
    // v refer to the object pointer
    var v reflect.Value
    if _, ok := obj.(reflect.Value); ok {
        v = obj.(reflect.Value)
    } else {
        v = reflect.ValueOf(obj)
    }

    if v.Kind() != reflect.Ptr {
        return errors.New("ClientConfigure: obj require a pointer or reflect.Value of a pointer")
    }

    // rv refer to the value of pointer
    rv := v.Elem()

    for key, val := range config {
        // change key to title string
        key = strings.Title(key)

        // check object's setter method
        if method := v.MethodByName("Set" + key); method.IsValid() {
            newVal := reflect.ValueOf(val).Convert(method.Type().In(0))
            errValues := method.Call([]reflect.Value{newVal})
            if len(errValues) != 0 && errValues[0].IsNil() == false {
                if errValues[0].MethodByName("Error").IsValid() {
                    return errors.New(errValues[0].MethodByName("Error").Call(nil)[0].String())
                }
            }
            continue
        }

        // check object's public field
        field := rv.FieldByName(key)
        if field.IsValid() && field.CanSet() {
            newVal := reflect.ValueOf(val).Convert(field.Type())
            field.Set(newVal)
            continue
        }
    }

    return nil
}
