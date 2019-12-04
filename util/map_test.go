package util

import "testing"

func TestMapClear(t *testing.T) {
    m := map[string]interface{}{"dd": "gg", "aa": "bb"}
    MapClear(m)
    if len(m) != 0 {
        t.FailNow()
    }
}

func TestMapGetter(t *testing.T) {
    m := map[string]interface{}{"dd": "gg", "aa": "bb"}
    MapSet(m, "cc.cck", "cvv")
    v := MapGet(m, "cc.cck")
    if v != "cvv" {
        t.FailNow()
    }

    if MapGet(m, "cc1") != nil {
        t.FailNow()
    }
}

func TestMapMerge(t *testing.T) {
    m := map[string]interface{}{"dd": "gg", "aa": "bb"}
    m1 := map[string]interface{}{"dd": "gg1", "aa": "bb", "cc": "ccv"}
    MapMerge(m, m1)
    if m["dd"] != "gg1" {
        t.Fatal(`m["dd"] != "gg1"`)
    }

    if _, has := m["cc"]; has == false {
        t.Fatal(`_,has := m["cc"];has==false`)
    }
}
