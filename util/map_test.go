package util

import (
	"reflect"
	"testing"
)

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
		t.Fatal(`v != "cvv"`)
	}

	MapSet(m, "cc1-cck1", "cvv", "-")
	vv := MapGet(m, "cc1-cck1", "-")
	if vv != "cvv" {
		t.Fatal(`vv != "cvv"`)
	}

	vvv := MapGet(m, "cc1.cck1")
	if vvv != "cvv" {
		t.Fatal(`vvv != "cvv"`)
	}

	if MapGet(m, "cc3") != nil {
		t.Fatal(`MapGet(m, "cc1") != nil`)
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

func TestParamsToMapSlice(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		data := map[string]interface{}{
			"advs[0][advId]":             "advid0",
			"advs[0][defaultLanguage]":   2,
			"advs[1][advId]":             "advid1",
			"advs[1][defaultLanguage]":   2,
			"advs[0][appVersionData][0]": "9.0.0",
			"advs[0][appVersionData][1]": "9.0.1",
		}

		ret := ParamsToMapSlice(data)
		if advs, ok := ret["advs"].([]interface{}); ok == false {
			t.Fatal(`_,ok:=ret["advs"].([]interface{});ok==false`)
		} else {
			if advsV, ok := advs[0].(map[string]interface{}); ok == false || advsV["advId"] != "advid0" {
				t.Fatal(`advsV,ok:=advs[0].(map[string]string);ok==false || advsV["advId"]!="advid0"`)
			}
		}
	})

	t.Run("int+string", func(t *testing.T) {
		data := map[string]interface{}{
			"advs[0][advId]":             "advid0",
			"advs[0][defaultLanguage]":   2,
			"advs[a][advId]":             "advid1",
			"advs[a][defaultLanguage]":   2,
			"advs[0][appVersionData][0]": "9.0.0",
			"advs[0][appVersionData][1]": "9.0.1",
			"aaa":                        "aaa",
		}

		ret := ParamsToMapSlice(data)
		if _, ok := ret["advs"].([]interface{}); ok == true {
			t.Fatal(`_,ok:=ret["advs"].([]interface{});ok==true`)
		}
	})

}

func TestMapToSliceString(t *testing.T) {
	data := map[string]interface{}{
		"0": "0",
		"1": "1",
	}

	ret := MapToSliceString(data)
	t.Log(reflect.TypeOf(ret))
	if ret[0] != "0" {
		t.Fatal(`ret[0] != "0"`)
	}
}

func TestMapToSliceInt(t *testing.T) {
	data := map[string]interface{}{
		"0": 0,
		"1": 1,
	}

	ret := MapToSliceInt(data)
	t.Log(reflect.TypeOf(ret))
	if ret[0] != 0 {
		t.Fatal(`ret[0] != 0`)
	}
}
