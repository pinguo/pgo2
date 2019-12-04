package value

import (
    "bytes"
    "encoding/json"
    "reflect"
    "strconv"
    "testing"
)

func TestNew(t *testing.T) {
    var obj interface{}
    obj = New("dd")
    if _, ok := obj.(*Value); ok == false {
        t.FailNow()
    }
}

func TestEncode(t *testing.T) {
    str := "dd"
    if bytes.Equal(Encode(str), []byte(str)) == false {
        t.FailNow()
    }
}

func TestDecode(t *testing.T) {
    str := "dd"
    ret := Encode(str)
    var dStr string
    Decode(ret, &dStr)
    if dStr != str {
        t.FailNow()
    }
}

func TestValue_Valid(t *testing.T) {
    v := New(nil)
    if v.Valid() != false {
        t.FailNow()
    }
}

func TestValue_TryEncode(t *testing.T) {
    v := New(TestValue_TryEncode)
    _, err := v.TryEncode()
    if err == nil {
        t.FailNow()
    }
}

func TestValue_TryDecode(t *testing.T) {
    v := New(TestValue_TryDecode)
    _, err := v.TryEncode()
    if err == nil {
        t.FailNow()
    }
}

func TestValue_Encode(t *testing.T) {
    t.Run("[]byte", func(t *testing.T) {
        in := []byte("dd")
        v := New(in)
        if bytes.Equal(v.Encode(), in) == false {
            t.FailNow()
        }
    })

    t.Run("string", func(t *testing.T) {
        in := "dd"
        v := New(in)
        if bytes.Equal(v.Encode(), []byte(in)) == false {
            t.FailNow()
        }
    })

    t.Run("bool", func(t *testing.T) {
        in := true
        v := New(in)
        var assert []byte
        assert = strconv.AppendBool(assert, in)
        if bytes.Equal(v.Encode(), assert) == false {
            t.FailNow()
        }
    })

    t.Run("float32, float64", func(t *testing.T) {
        var output []byte
        in := float64(123.11)
        f64 := reflect.ValueOf(in).Float()
        output = strconv.AppendFloat(output, f64, 'g', -1, 64)

        v := New(in)
        if bytes.Equal(v.Encode(), output) == false {
            t.Fatal("float64 err")
        }

        var assert []byte
        in1 := float32(123.11)
        f32 := reflect.ValueOf(in1).Float()
        assert = strconv.AppendFloat(assert, f32, 'g', -1, 64)

        v1 := New(in1)
        if bytes.Equal(v1.Encode(), assert) == false {
            t.Fatal("float32 err")
        }
    })

    t.Run(" int int8 int16 int32 int64", func(t *testing.T) {
        var v *Value
        var in interface{}
        var intOut []byte
        in = int(123)
        intOut = strconv.AppendInt(intOut, reflect.ValueOf(in).Int(), 10)
        v = New(in)
        if bytes.Equal(v.Encode(), intOut) == false {
            t.Fatal("int err")
        }

        intOut = nil
        in = int8(123)
        intOut = strconv.AppendInt(intOut, reflect.ValueOf(in).Int(), 10)
        v = New(in)
        if bytes.Equal(v.Encode(), intOut) == false {
            t.Fatal("int8 err")
        }

        intOut = nil
        in = int16(123)
        intOut = strconv.AppendInt(intOut, reflect.ValueOf(in).Int(), 10)
        v = New(in)
        if bytes.Equal(v.Encode(), intOut) == false {
            t.Fatal("int16 err")
        }

        intOut = nil
        in = int32(123)
        intOut = strconv.AppendInt(intOut, reflect.ValueOf(in).Int(), 10)
        v = New(in)
        if bytes.Equal(v.Encode(), intOut) == false {
            t.Fatal("int32 err")
        }

        intOut = nil
        in = int64(123)
        intOut = strconv.AppendInt(intOut, reflect.ValueOf(in).Int(), 10)
        v = New(in)
        if bytes.Equal(v.Encode(), intOut) == false {
            t.Fatal("int64 err")
        }
    })

    t.Run("uint uint8 uint16 uint32 uint64", func(t *testing.T) {
        var v *Value
        var in interface{}
        var intOut []byte
        in = uint(123)
        intOut = strconv.AppendUint(intOut, reflect.ValueOf(in).Uint(), 10)
        v = New(in)
        if bytes.Equal(v.Encode(), intOut) == false {
            t.Fatal("uint err")
        }

        intOut = nil
        in = uint8(123)
        intOut = strconv.AppendUint(intOut, reflect.ValueOf(in).Uint(), 10)
        v = New(in)
        if bytes.Equal(v.Encode(), intOut) == false {
            t.Fatal("uint8 err")
        }

        intOut = nil
        in = uint16(123)
        intOut = strconv.AppendUint(intOut, reflect.ValueOf(in).Uint(), 10)
        v = New(in)
        if bytes.Equal(v.Encode(), intOut) == false {
            t.Fatal("uint16 err")
        }

        intOut = nil
        in = uint32(123)
        intOut = strconv.AppendUint(intOut, reflect.ValueOf(in).Uint(), 10)
        v = New(in)
        if bytes.Equal(v.Encode(), intOut) == false {
            t.Fatal("uint32 err")
        }

        intOut = nil
        in = uint64(123)
        intOut = strconv.AppendUint(intOut, reflect.ValueOf(in).Uint(), 10)
        v = New(in)
        if bytes.Equal(v.Encode(), intOut) == false {
            t.Fatal("uint64 err")
        }
    })

    t.Run("other panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        v := New(TestValue_Encode)
        v.Encode()
    })

    t.Run("other json", func(t *testing.T) {
        v := New([]string{"a"})
        ret := v.Encode()
        var pm []string
        if err := json.Unmarshal(ret, &pm); err != nil {
            t.FailNow()
        }
    })

}

func TestValue_Decode(t *testing.T) {
    t.Run("[]byte", func(t *testing.T) {
        in := []byte("dd")
        v := New(in)

        var ret []byte
        v.Decode(&ret)
        if ret == nil {
            t.FailNow()
        }

    })

    t.Run("string", func(t *testing.T) {
        in := "dd"
        v := New(in)

        var ret string
        v.Decode(&ret)
        if ret == "" {
            t.FailNow()
        }

    })

    t.Run("bool", func(t *testing.T) {
        in := true
        v := New(in)
        var ret bool
        v.Decode(&ret)
        if ret == false {
            t.FailNow()
        }
    })

    t.Run("float32 float64", func(t *testing.T) {
        in := 123.11
        v := New(in)
        var ret float32
        v.Decode(&ret)
        if ret < 1 {
            t.FailNow()
        }

        in = 123.11
        v = New(in)
        var ret1 float32
        v.Decode(&ret1)
        if ret1 < 1 {
            t.FailNow()
        }
    })

    t.Run("int int8 int16 int32 int64", func(t *testing.T) {
        var v *Value
        var in interface{}

        in = 123

        var ret int
        v = New(in)
        v.Decode(&ret)
        if ret != in {
            t.FailNow()
        }

        var ret1 int8
        v = New(in)
        v.Decode(&ret1)
        if ret1 < 1 {
            t.FailNow()
        }

        var ret2 int16
        v = New(in)
        v.Decode(&ret2)
        if ret2 < 1 {
            t.FailNow()
        }

        var ret3 int32
        v = New(in)
        v.Decode(&ret3)
        if ret3 < 1 {
            t.FailNow()
        }

        var ret4 int64
        v = New(in)
        v.Decode(&ret4)
        if ret4 < 1 {
            t.FailNow()
        }
    })

    t.Run("uint uint8 uint16 uint32 uint64", func(t *testing.T) {
        var v *Value
        var in interface{}
        in = 123

        var ret uint
        v = New(in)
        v.Decode(&ret)
        if ret < 1 {
            t.FailNow()
        }

        var ret1 uint8
        v = New(in)
        v.Decode(&ret1)
        if ret1 < 1 {
            t.FailNow()
        }

        var ret2 uint16
        v = New(in)
        v.Decode(&ret2)
        if ret2 < 1 {
            t.FailNow()
        }

        var ret3 uint32
        v = New(in)
        v.Decode(&ret3)
        if ret3 < 1 {
            t.FailNow()
        }

        var ret4 uint64
        v = New(in)
        v.Decode(&ret4)
        if ret4 < 1 {
            t.FailNow()
        }
    })

    t.Run("other panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        v := New(TestValue_Encode)
        var ret int
        v.Decode(&ret)
    })

    t.Run("other json", func(t *testing.T) {
        v := New([]string{"a"})
        v.Encode()
        var pm []string
        v.Decode(&pm)
        if pm == nil {
            t.FailNow()
        }
    })
}

func TestValue_Data(t *testing.T) {
    str := "dd"
    v := New(str)
    if v.Data() != str {
        t.FailNow()
    }
}

func TestValue_Bool(t *testing.T) {
    in := true
    v := New(in)
    if v.Bool() != in {
        t.FailNow()
    }
}

func TestValue_Int(t *testing.T) {
    in := 11
    v := New(in)
    if v.Int() != in {
        t.FailNow()
    }
}

func TestValue_Float(t *testing.T) {
    in := 11.1
    v := New(in)
    if v.Float() < in {
        t.FailNow()
    }
}

func TestValue_String(t *testing.T) {
    in := "11.1"
    v := New(in)
    if v.String() != in {
        t.FailNow()
    }
}

func TestValue_Bytes(t *testing.T) {
    in := "aaa"
    v := New(in)
    if bytes.Equal(v.Bytes(), []byte(in)) == false {
        t.FailNow()
    }
}

func TestValue_MarshalJSON(t *testing.T) {
    str := `{"n":"v"}`
    v := New(str)
    _, err := v.MarshalJSON()
    if err != nil {
        t.FailNow()
    }

}

func TestValue_UnmarshalJSON(t *testing.T) {
    str := `{"n":"v"}`
    v := New("")
    err := v.UnmarshalJSON([]byte(str))
    if err != nil {
        t.FailNow()
    }
}
