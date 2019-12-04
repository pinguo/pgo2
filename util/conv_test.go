package util

import (
    "testing"
)

type mockTestToBool struct {
}

func TestToBool(t *testing.T) {
    t.Run("bool", func(t *testing.T) {
        if ToBool(true) != true {
            t.FailNow()
        }
    })

    t.Run("float32, float64", func(t *testing.T) {
        if ToBool(float32(1)) != true {
            t.FailNow()
        }

        if ToBool(float64(1)) != true {
            t.FailNow()
        }
    })

    t.Run("int, int8, int16, int32, int64", func(t *testing.T) {
        if ToBool(int(1)) != true {
            t.FailNow()
        }

        if ToBool(int8(1)) != true {
            t.FailNow()
        }

        if ToBool(int16(1)) != true {
            t.FailNow()
        }

        if ToBool(int32(1)) != true {
            t.FailNow()
        }

        if ToBool(int64(1)) != true {
            t.FailNow()
        }

    })

    t.Run("uint, uint8, uint16, uint32, uint64", func(t *testing.T) {
        if ToBool(uint(1)) != true {
            t.FailNow()
        }

        if ToBool(uint8(1)) != true {
            t.FailNow()
        }

        if ToBool(uint16(1)) != true {
            t.FailNow()
        }

        if ToBool(uint32(1)) != true {
            t.FailNow()
        }

        if ToBool(uint64(1)) != true {
            t.FailNow()
        }
    })

    t.Run("string", func(t *testing.T) {
        if ToBool("true") != true {
            t.FailNow()
        }
    })

    t.Run("[]byte", func(t *testing.T) {
        if ToBool([]byte("true")) != true {
            t.FailNow()
        }
    })

    t.Run("nil", func(t *testing.T) {
        if ToBool(nil) != false {
            t.FailNow()
        }
    })

    t.Run("default_obj", func(t *testing.T) {

        if ToBool(&mockTestToBool{}) != true {
            t.Fatal(`oBool(&mockTestToBool{}) != true`)
        }

        if ToBool([]string{"ddd"}) != true {
            t.Fatal(`ToBool([]string{"ddd"}) != true`)
        }
    })
}

func TestToInt(t *testing.T) {
    t.Run("bool", func(t *testing.T) {
        if ToInt(true) != 1 {
            t.Fatal(`ToInt(true) !=1`)
        }

        if ToInt(false) != 0 {
            t.Fatal(`ToInt(false) !=0`)
        }
    })

    t.Run("float32, float64", func(t *testing.T) {
        if ToInt(float32(12.11)) != 12 {
            t.Fatal(`oInt(float32(12.11)) != 12`)
        }

        if ToInt(float64(12.11)) != 12 {
            t.Fatal(`oInt(float32(12.11)) != 12`)
        }
    })

    t.Run("int, int8, int16, int32, int64", func(t *testing.T) {
        if ToInt(int(1)) != 1 {
            t.FailNow()
        }

        if ToInt(int8(1)) != 1 {
            t.FailNow()
        }

        if ToInt(int16(1)) != 1 {
            t.FailNow()
        }

        if ToInt(int32(1)) != 1 {
            t.FailNow()
        }

        if ToInt(int64(1)) != 1 {
            t.FailNow()
        }

    })

    t.Run("uint, uint8, uint16, uint32, uint64", func(t *testing.T) {
        if ToInt(uint(1)) != 1 {
            t.FailNow()
        }

        if ToInt(uint8(1)) != 1 {
            t.FailNow()
        }

        if ToInt(uint16(1)) != 1 {
            t.FailNow()
        }

        if ToInt(uint32(1)) != 1 {
            t.FailNow()
        }

        if ToInt(uint64(1)) != 1 {
            t.FailNow()
        }
    })

    t.Run("string", func(t *testing.T) {
        if ToInt("true") != 0 {
            t.FailNow()
        }
    })

    t.Run("[]byte", func(t *testing.T) {
        if ToInt([]byte("true")) != 0 {
            t.FailNow()
        }
    })

    t.Run("nil", func(t *testing.T) {
        if ToInt(nil) != 0 {
            t.FailNow()
        }
    })

    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }

            t.FailNow()
        }()
        ToInt(&mockTestToBool{})
    })
}

func TestToFloat(t *testing.T) {
    t.Run("bool", func(t *testing.T) {
        if ToFloat(true) != 1 {
            t.Fatal(`ToFloat(true) !=1`)
        }

        if ToFloat(false) != 0 {
            t.Fatal(`ToFloat(false) !=0`)
        }
    })

    t.Run("float32, float64", func(t *testing.T) {
        if len(ToString(ToFloat(float32(12.11)))) < 5 {
            t.Fatal(`len(ToString(ToFloat(float32(12.11)))) < 5`)
        }

        if len(ToString(ToFloat(float64(12.11)))) < 5 {
            t.Fatal(`len(ToString(ToFloat(float64(12.11)))) < 5`)
        }
    })

    t.Run("int, int8, int16, int32, int64", func(t *testing.T) {
        if ToFloat(int(1)) != 1 {
            t.FailNow()
        }

        if ToFloat(int8(1)) != 1 {
            t.FailNow()
        }

        if ToFloat(int16(1)) != 1 {
            t.FailNow()
        }

        if ToFloat(int32(1)) != 1 {
            t.FailNow()
        }

        if ToFloat(int64(1)) != 1 {
            t.FailNow()
        }

    })

    t.Run("uint, uint8, uint16, uint32, uint64", func(t *testing.T) {
        if ToFloat(uint(1)) != 1 {
            t.FailNow()
        }

        if ToFloat(uint8(1)) != 1 {
            t.FailNow()
        }

        if ToFloat(uint16(1)) != 1 {
            t.FailNow()
        }

        if ToFloat(uint32(1)) != 1 {
            t.FailNow()
        }

        if ToFloat(uint64(1)) != 1 {
            t.FailNow()
        }
    })

    t.Run("string", func(t *testing.T) {
        if ToFloat("true") != 0 {
            t.FailNow()
        }
    })

    t.Run("[]byte", func(t *testing.T) {
        if ToFloat([]byte("true")) != 0 {
            t.FailNow()
        }
    })

    t.Run("nil", func(t *testing.T) {
        if ToFloat(nil) != 0 {
            t.FailNow()
        }
    })

    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }

            t.FailNow()
        }()
        ToFloat(&mockTestToBool{})
    })
}

func TestToString(t *testing.T) {
    t.Run("bool", func(t *testing.T) {
        if ToString(true) != "true" {
            t.Fatal(`ToString(true) !="true"`)
        }

        if ToString(false) != "false" {
            t.Fatal(`ToString(false) !="false"`)
        }
    })

    t.Run("float32, float64", func(t *testing.T) {
        if len(ToString(float32(12.11))) < 5 {
            t.Fatal(`len(ToString(float32(12.11)))<5`)
        }

        if len(ToString(float64(12.11))) < 5 {
            t.Fatal(`len(ToString(float64(12.11)))<5`)
        }
    })

    t.Run("int, int8, int16, int32, int64", func(t *testing.T) {
        if ToString(int(1)) != "1" {
            t.FailNow()
        }

        if ToString(int8(1)) != "1" {
            t.FailNow()
        }

        if ToString(int16(1)) != "1" {
            t.FailNow()
        }

        if ToString(int32(1)) != "1" {
            t.FailNow()
        }

        if ToString(int64(1)) != "1" {
            t.FailNow()
        }

    })

    t.Run("uint, uint8, uint16, uint32, uint64", func(t *testing.T) {
        if ToString(uint(1)) != "1" {
            t.FailNow()
        }

        if ToString(uint8(1)) != "1" {
            t.FailNow()
        }

        if ToString(uint16(1)) != "1" {
            t.FailNow()
        }

        if ToString(uint32(1)) != "1" {
            t.FailNow()
        }

        if ToString(uint64(1)) != "1" {
            t.FailNow()
        }
    })

    t.Run("string", func(t *testing.T) {
        if ToString("true") != "true" {
            t.FailNow()
        }
    })

    t.Run("[]byte", func(t *testing.T) {
        if ToString([]byte("true")) != "true" {
            t.FailNow()
        }
    })

    t.Run("default", func(t *testing.T) {
        if ToString(&mockTestToBool{}) != "{}" {
            t.FailNow()
        }
    })
}
