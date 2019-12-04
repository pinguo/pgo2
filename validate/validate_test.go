package validate

import (
    "testing"
)

func TestBoolData(t *testing.T) {
    var obj interface{}
    obj = BoolData(false, "name", false)
    if _, ok := obj.(*Bool); ok == false {
        t.FailNow()
    }
}

func TestFloatData(t *testing.T) {
    var obj interface{}
    obj = FloatData("12.1", "name")
    if _, ok := obj.(*Float); ok == false {
        t.FailNow()
    }
}

func TestIntData(t *testing.T) {
    var obj interface{}
    obj = IntData("12", "name")
    if _, ok := obj.(*Int); ok == false {
        t.FailNow()
    }
}

func TestStringData(t *testing.T) {
    var obj interface{}
    obj = StringData("saaa", "name")
    if _, ok := obj.(*String); ok == false {
        t.FailNow()
    }
}

func TestValue(t *testing.T) {
    t.Run("value=map[string]interface{}", func(t *testing.T) {
        iv, useDft := Value(map[string]interface{}{"name": "v1"}, "name")
        if useDft != false {
            t.Fatal(`useDft!=false`)
        }
        if iv != "v1" {
            t.Fatal(`iv != "v1"`)
        }
    })

    t.Run("value=map[string]string", func(t *testing.T) {
        iv, useDft := Value(map[string]string{"name": "v1"}, "name")
        if useDft != false {
            t.Fatal(`useDft!=false`)
        }
        if iv != "v1" {
            t.Fatal(`iv != "v1"`)
        }
    })

    t.Run("value=map[string][]string", func(t *testing.T) {
        iv, useDft := Value(map[string][]string{"name": {"v1", "v2"}}, "name")
        if useDft != false {
            t.Fatal(`useDft!=false`)
        }

        if iv != "v1" {
            t.Fatal(`iv != "v1"`)
        }
    })

    t.Run("value=other", func(t *testing.T) {
        iv, useDft := Value("v1", "name")
        if useDft != false {
            t.Fatal(`useDft!=false`)
        }

        if iv != "v1" {
            t.Fatal(`iv != "v1"`)
        }
    })

    t.Run("value = nil use default", func(t *testing.T) {
        iv, useDft := Value(nil, "name", "dftV1")
        if useDft != true {
            t.Fatal(`useDft != true`)
        }

        if iv != "dftV1" {
            t.Fatal(`iv != "dftV1"`)
        }
    })

    t.Run("value = nil panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        Value(nil, "name")

    })

    t.Run("use default", func(t *testing.T) {
        iv, useDft := Value("", "name", "dftV1")
        if useDft != true {
            t.Fatal(`useDft != true`)
        }

        if iv != "dftV1" {
            t.Fatal(`iv != "dftV1"`)
        }
    })

    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        Value(" ", "name")
    })
}
