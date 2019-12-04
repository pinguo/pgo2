package pgo2

import (
    "bytes"
    "html/template"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestNewView(t *testing.T) {
    var obj interface{}
    obj = NewView(nil)
    if _, ok := obj.(*View); ok == false {
        t.FailNow()
    }
}

func TestView_SetSuffix(t *testing.T) {
    obj := NewView(nil)
    obj.SetSuffix("html")
    if obj.suffix != ".html" {
        t.FailNow()
    }
}

func TestView_SetCommons(t *testing.T) {
    App(true)
    defer App(true)
    App().viewPath, _ = filepath.Abs("./test/data")
    SetAlias("@view", App().viewPath)
    obj := NewView(nil)
    t.Run("errStyle", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        obj.SetCommons([]interface{}{1})
    })

    t.Run("normal", func(t *testing.T) {
        obj.SetCommons([]interface{}{"header", "footer"})
        if len(obj.commons) != 2 {
            t.Fatal(`len(obj.commons) !=2`)
        }

        if _, err := os.Stat(obj.commons[0]); err != nil {
            t.Fatal(`_,err:=os.Stat(obj.commons[0]);err!=nil`)
        }
    })
}

func TestView_AddFuncMap(t *testing.T) {
    obj := NewView(nil)
    obj.AddFuncMap(template.FuncMap{"fName": func(name string) string { return "f" + name }})
    tmpl := `{{ . | fName }}`

    temp := template.New("view.html").Funcs(obj.funcMap)
    temp, _ = temp.Parse(tmpl)

    ww := bytes.NewBuffer(nil)

    temp.Execute(ww, "name")

    if ww.String() != "fname" {
        t.FailNow()
    }

}

func TestView_Display(t *testing.T) {
    App(true)
    defer App(true)
    App().viewPath, _ = filepath.Abs("./test/data")
    SetAlias("@view", App().viewPath)
    obj := NewView(nil)
    obj.AddFuncMap(template.FuncMap{"fName": func(name string) string { return "f" + name }})
    obj.SetCommons([]interface{}{"header", "footer"})
    w := bytes.NewBuffer(nil)
    obj.Display(w, "view", nil)

    if strings.Index(w.String(), "test view") == -1 {
        t.FailNow()
    }
}

func TestView_Render(t *testing.T) {
    App(true)
    defer App(true)
    App().viewPath, _ = filepath.Abs("./test/data")
    SetAlias("@view", App().viewPath)
    obj := NewView(nil)
    obj.AddFuncMap(template.FuncMap{"fName": func(name string) string { return "f" + name }})
    obj.SetCommons([]interface{}{"header", "footer"})
    ret := obj.Render("view", nil)

    if strings.Index(string(ret), "test view") == -1 {
        t.FailNow()
    }
}
