package pgo2

import (
    "testing"

    "github.com/golang/mock/gomock"
    mock_config "github.com/pinguo/pgo2/test/mock/config"
)

func TestNewI18n(t *testing.T) {
    var obj interface{}
    obj = NewI18n(nil)
    if _, ok := obj.(*I18n); ok == false {
        t.FailNow()
    }
}

func TestI18n_SetSourceLang(t *testing.T) {
    i18n := NewI18n(nil)
    i18n.SetSourceLang("zh-CN")
    if i18n.sourceLang != "zh-CN" {
        t.FailNow()
    }

}

func TestI18n_SetTargetLang(t *testing.T) {
    i18n := NewI18n(nil)
    targetLang := []interface{}{"en", "zh-CN", "zh-TW"}
    i18n.SetTargetLang(targetLang)
    if len(i18n.targetLang) != len(targetLang) {
        t.FailNow()
    }

}

func TestI18n_loadMessage(t *testing.T) {
    t.Run("notExitsLang", func(t *testing.T) {
        i18n := NewI18n(nil)
        if i18n.loadMessage("test", "zh-CN") != "test" {
            t.Fatal(`i18n.loadMessage("test","zh-CN") !="test"`)
        }
    })

    t.Run("normal", func(t *testing.T) {
        mockTest := "mockTest"
        ctrl := gomock.NewController(t)
        defer ctrl.Finish()

        mockConfig := mock_config.NewMockIConfig(ctrl)

        mockConfig.EXPECT().GetString(gomock.Any(), gomock.Any()).Return(mockTest)

        App().SetConfig(mockConfig)
        i18n := NewI18n(nil)
        targetLang := []interface{}{"en", "zh-CN"}
        i18n.SetTargetLang(targetLang)
        if i18n.loadMessage("test", "zh-CN") != mockTest {
            t.Fatal(`i18n.loadMessage("test","zh-CN")!=`, mockTest)
        }

    })

}

func TestI18n_detectLang(t *testing.T) {
    t.Run("inTargetLang", func(t *testing.T) {
        i18n := NewI18n(nil)
        targetLang := []interface{}{"en", "zh-CN"}
        i18n.SetTargetLang(targetLang)
        lang := "zh-CN;q=0.9,en;q=0.8,zh-TW;q=0.7"
        if i18n.detectLang(lang) != "zh-CN" {
            t.FailNow()
        }

    })

    t.Run("notInTargetLang", func(t *testing.T) {
        i18n := NewI18n(nil)
        targetLang := []interface{}{"en", "zh"}
        i18n.SetTargetLang(targetLang)
        lang := "zh-CN;q=0.9,en;q=0.8,zh-TW;q=0.7"
        if i18n.detectLang(lang) != "zh" {
            t.FailNow()
        }

    })

    t.Run("sourceLang", func(t *testing.T) {
        i18n := NewI18n(nil)

        lang := "zh-CN;q=0.9,en;q=0.8,zh-TW;q=0.7"
        if i18n.detectLang(lang) != "en" {
            t.FailNow()
        }

    })
}

func TestI18n_Translate(t *testing.T) {
    mockTest := "mockTest%d"
    mockTestRet := "mockTest2019"
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockConfig := mock_config.NewMockIConfig(ctrl)

    mockConfig.EXPECT().GetString(gomock.Any(), gomock.Any()).Return(mockTest)

    App().SetConfig(mockConfig)
    i18n := NewI18n(nil)
    targetLang := []interface{}{"en", "zh-CN"}
    i18n.SetTargetLang(targetLang)

    if i18n.Translate("test", "zh-CN", 2019) != mockTestRet {
        t.Fatal(`i18n.Translate("test","zh-CN", 2019)!=`, mockTestRet)
    }
}
