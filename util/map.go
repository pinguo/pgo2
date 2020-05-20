package util

import (
	"fmt"
	"strings"
)


var mapToMapsRK = strings.NewReplacer("][", ".", "[", ".", "]", ".")

// MapClear clear map to empty
func MapClear(m map[string]interface{}) {
	for k := range m {
		delete(m, k)
	}
}

// MapMerge merge map recursively
func MapMerge(a map[string]interface{}, m ...map[string]interface{}) {
	for _, b := range m {
		for k := range b {
			va, oa := a[k].(map[string]interface{})
			vb, ob := b[k].(map[string]interface{})

			if oa && ob {
				MapMerge(va, vb)
			} else {
				a[k] = b[k]
			}
		}
	}
}

// MapGet get value by dot separated key, empty key for m itself
func MapGet(m map[string]interface{}, key string, dftSplit ...string) interface{} {
	var data interface{} = m
	split := "."
	if len(dftSplit) > 0 {
		split = dftSplit[0]
	}
	ks := strings.Split(key, split)

	for _, k := range ks {
		// skip empty key segment
		if k = strings.TrimSpace(k); len(k) == 0 {
			continue
		}

		if v, ok := data.(map[string]interface{}); ok {
			if data, ok = v[k]; ok {
				continue
			}
		}

		// not found
		return nil
	}

	return data
}

// MapSet set value by dot separated key, empty key for root, nil val for clear
func MapSet(m map[string]interface{}, key string, val interface{}, dftSplit ...string) {
	data := m
	last := ""

	split := "."
	if len(dftSplit) > 0 {
		split = dftSplit[0]
	}

	ks := strings.Split(key, split)
	for _, k := range ks {
		// skip empty key segment
		if k = strings.TrimSpace(k); len(k) == 0 {
			continue
		}

		if len(last) > 0 {
			if _, ok := data[last].(map[string]interface{}); !ok {
				data[last] = make(map[string]interface{})
			}

			data = data[last].(map[string]interface{})
		}

		last = k
	}

	if len(last) > 0 {
		if nil == val {
			delete(data, last)
		} else {
			data[last] = val
		}
	} else {
		MapClear(m)
		if v, ok := val.(map[string]interface{}); ok {
			MapMerge(m, v)
		} else if nil != val {
			panic(fmt.Sprintf("MapSet: invalid type: %T", val))
		}
	}
}




// 转换map为map/slice {"[aaa][0][bbb]":"sss"} => {"aaa":[{"bbb":"sss"}]}
// 主要是转换map 中key为数字的map 转换为slice
// 转换后的slice顺序不保证
func ParamsToMapSlice(m map[string]interface{}) map[string]interface{} {
	ks := make([]string, 0, len(m))
	mm := make(map[string]interface{})
	for k, v := range m {
		kk := strings.Trim(mapToMapsRK.Replace(k), ".")
		ks = append(ks, kk)
		MapSet(mm, kk, v)
	}

	changeData(mm)

	return mm

}

func changeData(mm map[string]interface{}) {
	data := mm
	for k, v := range data {
		vv, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		allNum := true

		for kkk, _ := range vv {
			if IsAllDigit(kkk) == false {
				allNum = false
				break
			}
		}

		if allNum == false {
			changeData(vv)
			continue
		}

		sliceData := make([]interface{}, 0, len(vv))
		for _, vvv := range vv {
			sliceData = append(sliceData, vvv)
		}
		data[k] = sliceData
		for _, v4 := range sliceData {
			if v5, ok := v4.(map[string]interface{}); ok {
				changeData(v5)
			}
		}

	}
}

// 强制转换map为[]string
// {"1":"ss"} => ["ss"]
// 转换后的slice顺序不保证
func MapToSliceString(sI interface{}) []string {
	s ,ok:=sI.(map[string]interface{})
	if ok == false{
		panic(fmt.Sprintf("MapToSliceString: invalid type: %T", sI))
	}

	sliceData := make([]string, 0, len(s))
	for _,v:=range s{
		sliceData = append(sliceData,ToString(v))
	}

	return sliceData
}

// 强制转换map为[]int
// {"1":1} => [1]
// 转换后的slice顺序不保证
func MapToSliceInt(sI interface{}) []int {
	s ,ok:=sI.(map[string]interface{})
	if ok == false{
		panic(fmt.Sprintf("MapToSliceString: invalid type: %T", sI))
	}

	sliceData := make([]int, 0, len(s))
	for _,v:=range s{
		if vInt,ok:= v.(int);ok == false{
			panic(fmt.Sprintf("MapToSliceInt: invalid type: %T", v))
		}else{
			sliceData = append(sliceData,vInt)
		}
	}

	return sliceData
}