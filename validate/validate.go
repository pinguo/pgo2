package validate

import (
    "net/http"
    "strings"

    "github.com/pinguo/pgo2/perror"
    "github.com/pinguo/pgo2/util"
)

// validate bool value
func BoolData(data interface{}, name string, dft ...interface{}) *Bool {
    value, useDft := Value(data, name, dft...)
    return &Bool{name, useDft, util.ToBool(value)}
}

// validate int value
func IntData(data interface{}, name string, dft ...interface{}) *Int {
    value, useDft := Value(data, name, dft...)
    return &Int{name, useDft, util.ToInt(value)}
}

// validate float value
func FloatData(data interface{}, name string, dft ...interface{}) *Float {
    value, useDft := Value(data, name, dft...)
    return &Float{name, useDft, util.ToFloat(value)}
}

// validate string value
func StringData(data interface{}, name string, dft ...interface{}) *String {
    value, useDft := Value(data, name, dft...)
    return &String{name, useDft, util.ToString(value)}
}

// get validate value, four situations:
// 1. data: map, name: field, dft[0]: default
// 2. data: map, name: field, dft: empty
// 3. data: value, name: field, dft[0]: default
// 4. data: value, name: field, dft: empty
func Value(data interface{}, name string, dft ...interface{}) (interface{}, bool) {
    var value interface{}
    var useDft = false

    switch v := data.(type) {
    case map[string]interface{}:
        if mv, ok := v[name]; ok {
            value = mv
        }
    case map[string]string:
        if mv, ok := v[name]; ok {
            value = mv
        }
    case map[string][]string:
        sliceValue, sliceOk := v[name]
        if sliceOk && len(sliceValue) > 0 {
            value = sliceValue[0]
        }
    default:
        value = data
    }

    if value == nil {
        if len(dft) == 1 {
            value = dft[0]
            useDft = true
        } else {
            panic(perror.New(http.StatusBadRequest, "%s is required", name))
        }
    } else if strValue, strOk := value.(string); strOk {
        strValue = strings.Trim(strValue, " \r\n\t")
        if len(strValue) > 0 {
            value = strValue
        } else if len(dft) == 1 {
            value = dft[0]
            useDft = true
        } else {
            panic(perror.New(http.StatusBadRequest, "%s can't be empty", name))
        }
    }

    return value, useDft
}
