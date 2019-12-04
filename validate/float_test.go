package validate

import "testing"

func TestFloat_Min(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        f := &Float{Name: "name", UseDft: false, Value: 2}
        f.Min(3)
    })

    t.Run("normal", func(t *testing.T) {

        f := &Float{Name: "name", UseDft: false, Value: 2}
        if f.Min(1) != f {
            t.FailNow()
        }
    })

}

func TestFloat_Max(t *testing.T) {
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        f := &Float{Name: "name", UseDft: false, Value: 2}
        f.Max(1)
    })

    t.Run("normal", func(t *testing.T) {

        f := &Float{Name: "name", UseDft: false, Value: 2}
        if f.Max(3) != f {
            t.FailNow()
        }
    })
}

func TestFloat_Do(t *testing.T) {
    f := &Float{Name: "name", UseDft: false, Value: 2}
    if f.Do() != 2 {
        t.FailNow()
    }
}

func TestFloatSlice_Do(t *testing.T) {
    f := FloatSlice{Name: "name", Value: []float64{1, 2}}
    if len(f.Do()) != 2 {
        t.FailNow()
    }
}
