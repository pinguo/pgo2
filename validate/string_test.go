package validate

import (
    "regexp"
    "testing"
)

func TestString_Min(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        s.Min(4)
    })

    t.Run("normal", func(t *testing.T) {
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        if s.Min(2) != s {
            t.FailNow()
        }
    })
}

func TestString_Max(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        s.Max(2)
    })

    t.Run("normal", func(t *testing.T) {
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        if s.Max(4) != s {
            t.FailNow()
        }
    })
}

func TestString_Len(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        s.Len(4)
    })

    t.Run("normal", func(t *testing.T) {
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        if s.Len(3) != s {
            t.FailNow()
        }
    })
}

func TestString_Enum(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        s.Enum("bbb")
    })

    t.Run("normal", func(t *testing.T) {
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        if s.Enum("aaa", "bbb") != s {
            t.FailNow()
        }
    })
}

func TestString_RegExp(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        s.RegExp("bbb")
    })

    t.Run("normal", func(t *testing.T) {
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        re, _ := regexp.Compile(`(a+)`)
        if s.RegExp(re) != s {
            t.FailNow()
        }
    })
}

func TestString_Filter(t *testing.T) {
    t.Run("panic1", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        s.Filter(func(v, n string) string {
            return ""
        })
    })

    t.Run("panic2", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        s.Filter(func(v, n string) string {
            panic("err")
            return ""
        })
    })

    t.Run("normal", func(t *testing.T) {
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        if s.Filter(func(v, n string) string {
            return v
        }) != s {
            t.FailNow()
        }
    })
}

func TestString_Password(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        s.Password()
    })

    t.Run("normal", func(t *testing.T) {
        s := &String{Name: "name", UseDft: false, Value: "12Aaa,%"}
        if s.Password() != s {
            t.FailNow()
        }
    })
}

func TestString_Email(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        s.Email()
    })

    t.Run("normal", func(t *testing.T) {
        s := &String{Name: "name", UseDft: false, Value: "aa@aa.com"}
        if s.Email() != s {
            t.FailNow()
        }
    })
}

func TestString_Mobile(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "aaa"}
        s.Mobile()
    })

    t.Run("normal", func(t *testing.T) {
        s := &String{Name: "name", UseDft: false, Value: "13111111111"}
        if s.Mobile() != s {
            t.FailNow()
        }
    })
}

func TestString_IPv4(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "127.0.0"}
        s.IPv4()
    })

    t.Run("normal", func(t *testing.T) {
        s := &String{Name: "name", UseDft: false, Value: "127.0.0.1"}
        if s.IPv4() != s {
            t.FailNow()
        }
    })
}

func TestString_Bool(t *testing.T) {
    s := &String{Name: "name", UseDft: false, Value: "false"}
    var obj interface{}
    obj = s.Bool()
    if _, ok := obj.(*Bool); ok == false {
        t.FailNow()
    }
}

func TestString_Int(t *testing.T) {
    s := &String{Name: "name", UseDft: false, Value: "1"}
    var obj interface{}
    obj = s.Int()
    if _, ok := obj.(*Int); ok == false {
        t.FailNow()
    }
}

func TestString_Float(t *testing.T) {
    s := &String{Name: "name", UseDft: false, Value: "1.0"}
    var obj interface{}
    obj = s.Float()
    if _, ok := obj.(*Float); ok == false {
        t.FailNow()
    }

}

func TestString_Slice(t *testing.T) {
    s := &String{Name: "name", UseDft: false, Value: "a,b"}
    var obj interface{}
    obj = s.Slice(",")
    if _, ok := obj.(*StringSlice); ok == false {
        t.FailNow()
    }
}

func TestString_Json(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &String{Name: "name", UseDft: false, Value: "a,b"}
        s.Json()
    })

    t.Run("normal", func(t *testing.T) {
        s := &String{Name: "name", UseDft: false, Value: `{"aa":"ddd"}`}
        var obj interface{}
        obj = s.Json()
        if _, ok := obj.(*Json); ok == false {
            t.FailNow()
        }
    })

}

func TestString_Do(t *testing.T) {
    s := &String{Name: "name", UseDft: false, Value: "a"}
    if s.Do() != s.Value {
        t.FailNow()
    }
}

func TestStringSlice_Float(t *testing.T) {

    s := &StringSlice{Name: "name", UseDft: false, Value: []string{"1.2", "1.3"}}
    var obj interface{}
    obj = s.Float()
    if _, ok := obj.(*FloatSlice); ok == false {
        t.FailNow()
    }
}

func TestStringSlice_Int(t *testing.T) {
    s := &StringSlice{Name: "name", UseDft: false, Value: []string{"1", "2"}}
    var obj interface{}
    obj = s.Int()
    if _, ok := obj.(*IntSlice); ok == false {
        t.FailNow()
    }
}

func TestStringSlice_Min(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &StringSlice{Name: "name", UseDft: false, Value: []string{"1", "2"}}
        s.Min(3)
    })

    t.Run("normal", func(t *testing.T) {
        s := &StringSlice{Name: "name", UseDft: false, Value: []string{"1", "2"}}
        if s.Min(1) != s {
            t.FailNow()
        }
    })
}

func TestStringSlice_Max(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &StringSlice{Name: "name", UseDft: false, Value: []string{"1", "2"}}
        s.Max(1)
    })

    t.Run("normal", func(t *testing.T) {
        s := &StringSlice{Name: "name", UseDft: false, Value: []string{"1", "2"}}
        if s.Max(3) != s {
            t.FailNow()
        }
    })
}

func TestStringSlice_Len(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        s := &StringSlice{Name: "name", UseDft: false, Value: []string{"1", "2"}}
        s.Len(1)
    })

    t.Run("normal", func(t *testing.T) {
        s := &StringSlice{Name: "name", UseDft: false, Value: []string{"1", "2"}}
        if s.Len(2) != s {
            t.FailNow()
        }
    })
}

func TestStringSlice_Do(t *testing.T) {
    s := &StringSlice{Name: "name", UseDft: false, Value: []string{"1", "2"}}
    if len(s.Do()) != len(s.Value) {
        t.FailNow()
    }
}
