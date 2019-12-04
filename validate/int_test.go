package validate

import "testing"

func TestInt_Min(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        f := &Int{Name: "name", UseDft: false, Value: 2}
        f.Min(3)
    })

    t.Run("normal", func(t *testing.T) {

        f := &Int{Name: "name", UseDft: false, Value: 2}
        if f.Min(1) != f {
            t.FailNow()
        }
    })

}

func TestInt_Max(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        f := &Int{Name: "name", UseDft: false, Value: 2}
        f.Max(1)
    })

    t.Run("normal", func(t *testing.T) {

        f := &Int{Name: "name", UseDft: false, Value: 2}
        if f.Max(3) != f {
            t.FailNow()
        }
    })
}

func TestInt_Enum(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        f := &Int{Name: "name", UseDft: false, Value: 2}
        f.Enum(1, 3)
    })

    t.Run("normal", func(t *testing.T) {

        f := &Int{Name: "name", UseDft: false, Value: 2}
        if f.Enum(3, 2) != f {
            t.FailNow()
        }
    })
}

func TestInt_Do(t *testing.T) {
    f := &Int{Name: "name", UseDft: false, Value: 2}
    if f.Do() != 2 {
        t.FailNow()
    }
}

func TestIntSlice_Do(t *testing.T) {
    f := IntSlice{Name: "name", Value: []int{1, 2}}
    if len(f.Do()) != 2 {
        t.FailNow()
    }
}
