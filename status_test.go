package pgo2

import (
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/pinguo/pgo2/iface"
    mock_iface "github.com/pinguo/pgo2/test/mock/iface"
)

func TestNewStatus(t *testing.T) {
    var obj interface{}

    obj = NewStatus(nil)

    if _, ok := obj.(iface.IStatus); ok == false {
        t.FailNow()
    }
}

func TestStatus_SetUseI18n(t *testing.T) {
    status := NewStatus(nil)
    status.SetUseI18n(true)
    if status.useI18n == false {
        t.FailNow()
    }
}

func TestStatus_SetMapping(t *testing.T) {
    status := NewStatus(nil)
    status.SetMapping(map[string]interface{}{"11": "test"})
    if status.mapping[11] != "test" {
        t.FailNow()
    }
}

func TestStatus_Text(t *testing.T) {
    status := NewStatus(nil)
    status.SetMapping(map[string]interface{}{"11": "test"})
    t.Run("panic", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        status.Text(111, "zh-CN")
    })

    t.Run("default", func(t *testing.T) {
        if status.Text(111, "zh-CN", "dftTest") != "dftTest" {
            t.FailNow()
        }
    })

    t.Run("useI18n", func(t *testing.T) {
        defer App(true)
        status.SetUseI18n(true)
        status.SetMapping(map[string]interface{}{"111": "test"})
        mockTest := "mockTest"
        ctrl := gomock.NewController(t)
        defer ctrl.Finish()

        mockI18n := mock_iface.NewMockII18n(ctrl)
        mockI18n.EXPECT().Translate(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(message, lang interface{}, params ...interface{}) {
            t.Log("mock I18n.Translate")
        }).Return(mockTest)

        App().SetI18n(mockI18n)

        if status.Text(111, "zh-Cn") != mockTest {
            t.FailNow()
        }

    })
}
