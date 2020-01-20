package value

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pinguo/pgo2/util"
)

func New(data interface{}) *Value {
	return &Value{data}
}

// Encode encode data to bytes
func Encode(data interface{}) []byte {
	v := Value{data}
	return v.Encode()
}

// Decode data bytes to ptr
func Decode(data interface{}, ptr interface{}) {
	v := Value{data}
	v.Decode(ptr)
}

// Value adapter for value of any type,
// provide uniform encoding and decoding.
type Value struct {
	data interface{}
}

// Valid check the underlying data is nil
func (v *Value) Valid() bool {
	return v.data != nil
}

// TryEncode try encode data, err is not nil if panic
func (v *Value) TryEncode() (output []byte, err error) {
	defer func() {
		if v := recover(); v != nil {
			output, err = nil, errors.New(util.ToString(v))
		}
	}()
	return v.Encode(), nil
}

// TryDecode try decode data, err is not nil if panic
func (v *Value) TryDecode(ptr interface{}) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = errors.New(util.ToString(v))
		}
	}()
	v.Decode(ptr)
	return nil
}

// Encode encode data to bytes, panic if failed
func (v *Value) Encode() []byte {
	var output []byte
	switch d := v.data.(type) {
	case []byte:
		output = d
	case string:
		output = []byte(d)
	case bool:
		output = strconv.AppendBool(output, d)
	case float32, float64:
		f64 := reflect.ValueOf(v.data).Float()
		output = strconv.AppendFloat(output, f64, 'g', -1, 64)
	case int, int8, int16, int32, int64:
		i64 := reflect.ValueOf(v.data).Int()
		output = strconv.AppendInt(output, i64, 10)
	case uint, uint8, uint16, uint32, uint64:
		u64 := reflect.ValueOf(v.data).Uint()
		output = strconv.AppendUint(output, u64, 10)
	default:
		if j, e := json.Marshal(v.data); e == nil {
			output = j
		} else {
			panic("Value.Encode: " + e.Error())
		}
	}
	return output
}

// Decode decode data bytes to ptr, panic if failed
func (v *Value) Decode(ptr interface{}) {
	switch p := ptr.(type) {
	case *[]byte:
		*p = v.Bytes()
	case *string:
		*p = v.String()
	case *bool:
		*p = util.ToBool(v.data)
	case *float32, *float64:
		fv := util.ToFloat(v.data)
		rv := reflect.ValueOf(ptr).Elem()
		rv.Set(reflect.ValueOf(fv).Convert(rv.Type()))
	case *int, *int8, *int16, *int32, *int64:
		iv := util.ToInt(v.data)
		rv := reflect.ValueOf(ptr).Elem()
		rv.Set(reflect.ValueOf(iv).Convert(rv.Type()))
	case *uint, *uint8, *uint16, *uint32, *uint64:
		iv := util.ToInt(v.data)
		rv := reflect.ValueOf(ptr).Elem()
		rv.Set(reflect.ValueOf(iv).Convert(rv.Type()))
	default:
		if e := json.Unmarshal(v.Bytes(), ptr); e != nil {
			rv := reflect.ValueOf(ptr)
			if rv.Kind() != reflect.Ptr || rv.IsNil() {
				panic("Value.Decode: require a valid pointer")
			}

			if rv = rv.Elem(); rv.Kind() == reflect.Interface {
				rv.Set(reflect.ValueOf(v.data))
			} else {
				panic("Value.Decode: " + e.Error())
			}
		}
	}
}

// Data return underlying data
func (v *Value) Data() interface{} {
	return v.data
}

// Bool return underlying data as bool
func (v *Value) Bool() bool {
	return util.ToBool(v.data)
}

// Int return underlying data as int
func (v *Value) Int() int {
	return util.ToInt(v.data)
}

// Float return underlying data as float64
func (v *Value) Float() float64 {
	return util.ToFloat(v.data)
}

// String return underlying data as string
func (v *Value) String() string {
	switch d := v.data.(type) {
	case []byte:
		return string(d)
	case string:
		return d
	default:
		if j, e := json.Marshal(v.data); e == nil {
			return string(j)
		}
		return fmt.Sprintf("%+v", v.data)
	}
}

// Bytes return underlying data as bytes
func (v *Value) Bytes() []byte {
	switch d := v.data.(type) {
	case []byte:
		return d
	case string:
		return []byte(d)
	default:
		if j, e := json.Marshal(v.data); e == nil {
			return j
		}
		return []byte(fmt.Sprintf("%+v", v.data))
	}
}

func (v *Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v *Value) UnmarshalJSON(b []byte) error {
	var s interface{}
	if e := json.Unmarshal(b, &s); e != nil {
		return e
	}
	v.data = s
	return nil
}
