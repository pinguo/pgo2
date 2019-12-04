package pgo2

import (
    "os"
    "path/filepath"
    "strings"
    "testing"

    "github.com/pinguo/pgo2/iface"
)

func TestApplication_InitArgs(t *testing.T) {

    t.Run("space", func(t *testing.T) {
        args := []string{"xxxx", "--name", "nameTest0", "xxxx"}
        ret := App().initArgs(args)
        if ret["name"] != args[2] {
            t.FailNow()
        }
    })

    t.Run("eq", func(t *testing.T) {
        args := []string{"--aa", "--name=nameTest1", "--nameTest1", "--xxxx"}
        ret := App().initArgs(args)
        if ret["name"] != "nameTest1" {
            t.FailNow()
        }

        if App().HasArg("name") == false {
            t.FailNow()
        }

        if App().Arg("name") == "" {
            t.FailNow()
        }
    })

}

func TestApplication_Getter(t *testing.T) {
    exeBase := filepath.Base(os.Args[0])
    exeExt := filepath.Ext(os.Args[0])
    exeDir := filepath.Dir(os.Args[0])
    app := &Application{
        mode:       ModeWeb,
        env:        DefaultEnv,
        name:       strings.TrimSuffix(exeBase, exeExt),
        components: make(map[string]interface{}),
    }

    app.basePath = app.genBasePath(exeDir)
    osArgs := os.Args
    osArgs = append(osArgs, "--cmd=/test/index")
    app.Init(osArgs)
    app.Name()

    if app.Mode() != ModeCmd {
        t.Fatal("app.Mode() != ModeCmd")
    }

    if app.RuntimePath() == "" {
        t.Fatal("app.RuntimePath() == \"\"")
    }

    if app.BasePath() == "" {
        t.Fatal("app.BasePath() = \"\"")
    }

    if app.PublicPath() == "" {
        t.Fatal("app.PublicPath() == \"\" ")
    }

    if app.ViewPath() == "" {
        t.Fatal("app.ViewPath() == \"\" ")
    }

    if app.Config() == nil {
        t.Fatal("app.Config() == nil  ")
    }

    if app.Container() == nil {
        t.Fatal("app.Container() == nil ")
    }

    if app.Server() == nil {
        t.Fatal("app.Server() == nil ")
    }

    if app.Router() == nil {
        t.Fatal("app.Router() == nil ")
    }

    if app.Log() == nil {
        t.Fatal("app.Log() == nil ")
    }

    if app.Status() == nil {
        t.Fatal("app.Status() == nil ")
    }

    if app.I18n() == nil {
        t.Fatal("app.I18n() == nil ")
    }

    if app.View() == nil {
        t.Fatal("app.View() == nil ")
    }

    if app.StopBefore() == nil {
        t.Fatal("app.StopBefore() == nil ")
    }

}

func newComponentTest(config map[string]interface{}) (interface{}, error) {
    return &componentTest{}, nil
}

type componentTest struct {
}

func TestApplication_Component(t *testing.T) {

    t.Run("haveParam", func(t *testing.T) {
        params := make(map[string]interface{})
        params["name"] = "ddd"
        if _, ok := App().Component("testId", newComponentTest, params).(*componentTest); ok == false {
            t.FailNow()
        }

    })

    t.Run("nilParam", func(t *testing.T) {
        if _, ok := App().Component("testId1", newComponentTest).(*componentTest); ok == false {
            t.FailNow()
        }

    })

    t.Run("exitsParam", func(t *testing.T) {
        if _, ok := App().Component("testId1", newComponentTest).(*componentTest); ok == false {
            t.FailNow()
        }

    })
}

func newObjSingleTest(params ...interface{}) iface.IObject {
    return &objSingleTest{}
}

type objSingleTest struct {
    Object
}

func TestApplication_GetObjSingle(t *testing.T) {
    t.Run("new", func(t *testing.T) {
        if _, ok := App().GetObjSingle("testSingleId", newObjSingleTest).(iface.IObject); ok == false {
            t.FailNow()
        }
    })

    t.Run("exits", func(t *testing.T) {
        if _, ok := App().GetObjSingle("testSingleId", newObjSingleTest).(iface.IObject); ok == false {
            t.FailNow()
        }
    })

}
