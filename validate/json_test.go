package validate

import "testing"

func TestJson_Has(t *testing.T) {
    v := map[string]interface{}{"name": "v"}
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()

        j := &Json{Name: "name", UseDft: false, Value: v}
        j.Has("name1")
    })

    t.Run("normal", func(t *testing.T) {
        j := &Json{Name: "name", UseDft: false, Value: v}
        if j.Has("name") != j {
            t.FailNow()
        }
    })

}

func TestJson_Do(t *testing.T) {
    v := map[string]interface{}{"name": "v"}
    j := &Json{Name: "name", UseDft: false, Value: v}
    if j.Do()["name"] != "v" {
        t.FailNow()
    }
}
