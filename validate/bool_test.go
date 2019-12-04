package validate

import "testing"

func TestBool_Must(t *testing.T) {

    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        b := &Bool{Name: "name", UseDft: false, Value: true}
        b.Must(false)
    })

    t.Run("normal", func(t *testing.T) {

        b := &Bool{Name: "name", UseDft: false, Value: true}
        if b.Must(true) != b {
            t.FailNow()
        }
    })
}

func TestBool_Do(t *testing.T) {
    b := &Bool{Name: "name", UseDft: false, Value: true}
    if b.Do() != b.Value {
        t.FailNow()
    }
}
