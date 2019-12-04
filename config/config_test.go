package config

import (
    "path/filepath"
    "testing"
)

var mockTestBasePath, _ = filepath.Abs("../test/data")
var mockTestEnv = "docker"

func TestNew(t *testing.T) {
    var config interface{}
    config = New(mockTestBasePath, mockTestEnv)
    if _, ok := config.(IConfig); ok == false {
        t.FailNow()
    }

    c, ok := config.(*Config)
    if ok == false {
        t.Fatal(` c,ok := config.(*Config);ok==false`)
    }

    if _, has := c.parsers["json"]; has == false {
        t.Fatal(`_,has:=c.parsers["json"];has==false`)
    }

    if _, has := c.parsers["yaml"]; has == false {
        t.Fatal(`_,has:=c.parsers["yaml"];has==false`)
    }
}

func TestConfig_CheckPath(t *testing.T) {
    t.Run("invalidPath", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        config := New(mockTestBasePath+"fffggggg", mockTestEnv)
        config.CheckPath()
    })

    t.Run("invalidEnv", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        config := New(mockTestBasePath, mockTestEnv+"fffffggg")
        config.CheckPath()
    })

    t.Run("normal", func(t *testing.T) {
        config := New(mockTestBasePath, mockTestEnv)
        config.CheckPath()
    })
}

func TestConfig_AddPath(t *testing.T) {
    config := New(mockTestBasePath, mockTestEnv)
    config.AddPath(mockTestBasePath + "/conf")

    if len(config.paths) != 2 {
        t.Fatal(`len(config.paths) != 2`)
    }

    config.AddPath(mockTestBasePath + "/conf/testing")

    if len(config.paths) != 3 {
        t.Fatal(`len(config.paths) != 3`)
    }
}

type mockTestXmlParser struct {
    base
}

func (m *mockTestXmlParser) Parse(path string) (parseData map[string]interface{}, err error) {
    return nil, nil
}

func TestConfig_AddParser(t *testing.T) {
    config := New(mockTestBasePath, mockTestEnv)
    config.AddParser("xml", &mockTestXmlParser{})
    if _, has := config.parsers["xml"]; has == false {
        t.FailNow()
    }
}

func TestConfig_Getter(t *testing.T) {
    config := New(mockTestBasePath, mockTestEnv)
    t.Run("GetBool", func(t *testing.T) {
        if config.GetBool("testjson.testBool", false) != true {
            t.Fatal(`config.GetBool("testjson.testBool", false) !=true`)
        }

        if config.GetBool("testyaml.testBool", false) != true {
            t.Fatal(`config.GetBool("testyaml.testBool", false) !=true`)
        }
    })

    t.Run("GetBoolDefault", func(t *testing.T) {
        if config.GetBool("testjson.testBoolDft", true) != true {
            t.Fatal(`config.GetBool("testjson.testBoolDft", false) !=true`)
        }

        if config.GetBool("testyaml.testBoolDft", true) != true {
            t.Fatal(`config.GetBool("testyaml.testBoolDft", false) !=true`)
        }
    })

    t.Run("GetString", func(t *testing.T) {
        if config.GetString("testjson.testString", "") == "" {
            t.Fatal(`config.GetString("testjson.GetString", "") ==""`)
        }

        if config.GetString("testyaml.testString", "") == "" {
            t.Fatal(`config.GetString("testyaml.GetString", "") ==""`)
        }
    })

    t.Run("GetStringDefault", func(t *testing.T) {
        if config.GetString("testjson.testStringDft", "dft") != "dft" {
            t.Fatal(`config.GetString("testjson.testStringDft", "dft") !="dft"`)
        }

        if config.GetString("testyaml.testStringDft", "dft") != "dft" {
            t.Fatal(`config.GetString("testyaml.testStringDft", "dft") !="dft"`)
        }
    })

    t.Run("GetInt", func(t *testing.T) {
        if config.GetInt("testjson.testInt", 0) == 0 {
            t.Fatal(`config.GetInt("testjson.testInt", 0) ==0`)
        }

        if config.GetInt("testyaml.testInt", 0) == 0 {
            t.Fatal(`config.GetInt("testyaml.testInt", 0) ==0`)
        }
    })

    t.Run("GetIntDefault", func(t *testing.T) {
        if config.GetInt("testjson.testIntDft", 100) != 100 {
            t.Fatal(`config.GetInt("testjson.testIntDft", 100) !=100`)
        }

        if config.GetInt("testyaml.testIntDft", 100) != 100 {
            t.Fatal(`config.GetInt("testyaml.testIntDft", 100) !=100`)
        }
    })

    t.Run("GetFloat", func(t *testing.T) {
        if config.GetFloat("testjson.testFloat", 0) == 0 {
            t.Fatal(`config.GetFloat("testjson.testFloat", 0) ==0`)
        }

        if config.GetFloat("testyaml.testFloat", 0) == 0 {
            t.Fatal(`config.GetFloat("testyaml.testFloat", 0) ==0`)
        }
    })

    t.Run("GetFloatDefault", func(t *testing.T) {
        if config.GetFloat("testjson.testFloatDft", 100) != 100 {
            t.Fatal(`config.GetFloat("testjson.testFloatDft", 100) !=100`)
        }

        if config.GetFloat("testyaml.testFloatDft", 100) != 100 {
            t.Fatal(`config.GetFloat("testyaml.testFloatDft", 100) !=100`)
        }
    })
}

func TestConfig_GetSlice(t *testing.T) {
    config := New(mockTestBasePath, mockTestEnv)
    t.Run("GetSliceBool", func(t *testing.T) {
        if len(config.GetSliceBool("testjson.testSBool")) != 2 {
            t.Fatal(`len(config.GetSliceBool("testjson.testSBool"))  != 2`)
        }

        if len(config.GetSliceBool("testyaml.testSBool")) != 2 {
            t.Fatal(`len(config.GetSliceBool("testyaml.testSBool")) !=2`)
        }
    })

    t.Run("GetSliceString", func(t *testing.T) {
        if len(config.GetSliceString("testjson.testSString")) != 2 {
            t.Fatal(`len(config.GetSliceString("testjson.testSString")) !=2`)
        }

        if len(config.GetSliceString("testyaml.testSString")) != 2 {
            t.Fatal(`len(config.GetSliceString("testyaml.testSString", ""))!=2`)
        }
    })

    t.Run("GetSliceInt", func(t *testing.T) {
        if len(config.GetSliceInt("testjson.testSInt")) != 2 {
            t.Fatal(`len(config.GetSliceInt("testjson.testSInt")) !=2`)
        }

        if len(config.GetSliceInt("testyaml.testSInt")) != 2 {
            t.Fatal(`len(config.GetSliceInt("testyaml.testInt")) !=2`)
        }
    })

    t.Run("GetSliceFloat", func(t *testing.T) {
        if len(config.GetSliceFloat("testjson.testSFloat")) != 2 {
            t.Fatal(`len(config.GetSliceFloat("testjson.testSFloat")) !=2`)
        }

        if len(config.GetSliceFloat("testyaml.testSFloat")) != 2 {
            t.Fatal(`len(config.GetSliceFloat("testyaml.testSFloat")) != 2`)
        }
    })
}

func TestConfig_Set(t *testing.T) {
    config := New(mockTestBasePath, mockTestEnv)
    config.Set("mockTestSet", "vv")
    if config.GetString("mockTestSet", "") != "vv" {
        t.FailNow()
    }
}
