package logs

import (
    "testing"
)

func TestTarget_SetLevels(t *testing.T) {
    t.Run("int", func(t *testing.T) {
        target := &Target{}
        target.SetLevels(LevelNone)
        if target.levels != LevelNone {
            t.FailNow()
        }
    })

    t.Run("string", func(t *testing.T) {
        target := &Target{}
        target.SetLevels("NONE")
        if target.levels != LevelNone {
            t.FailNow()
        }
    })

    t.Run("invalid", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        target := &Target{}
        target.SetLevels("12333")

    })
}

type mockTestFormatter struct {
}

func (m *mockTestFormatter) Format(item *LogItem) string {
    return "mock"
}

func TestTarget_SetFormatter(t *testing.T) {
    t.Run("invalid", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        target := &Target{}
        target.SetFormatter("12333")
    })

    t.Run("normal", func(t *testing.T) {

        target := &Target{}
        target.SetFormatter(&mockTestFormatter{})
    })
}

func TestTarget_IsHandling(t *testing.T) {
    target := &Target{}
    target.SetLevels(LevelDebug)
    if target.IsHandling(LevelDebug) != true {
        t.FailNow()
    }
}

func TestTarget_Format(t *testing.T) {
    t.Run("customFormater", func(t *testing.T) {
        target := &Target{}
        target.SetFormatter(&mockTestFormatter{})
        if target.Format(&LogItem{}) != "mock" {
            t.FailNow()
        }
    })

    t.Run("defaultFormater", func(t *testing.T) {
        target := &Target{}
        if target.Format(&LogItem{}) == "" {
            t.FailNow()
        }
    })
}
