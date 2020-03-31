package util

// SliceSearch search in a slice has the length of n,
// return the first position where f(i) is true,
// -1 is returned if nothing found.
func SliceSearch(n int, f func(int) bool) int {
	for i := 0; i < n; i++ {
		if f(i) {
			return i
		}
	}

	return -1
}

// SliceSearchInt search x in an int slice, return the first position of x,
// -1 is returned if nothing found.
func SliceSearchInt(a []int, x int) int {
	return SliceSearch(len(a), func(i int) bool { return a[i] == x })
}

// SliceSearchString search x in a string slice, return the first position of x,
// -1 is returned if nothing found.
func SliceSearchString(a []string, x string) int {
	return SliceSearch(len(a), func(i int) bool { return a[i] == x })
}

// SliceUniqueInt retrieve unique items item from int slice.
func SliceUniqueInt(a []int) []int {
	l := int(len(a)/2) + 1
	result, exists := make([]int, 0, l), make(map[int]bool)
	for _, v := range a {
		if !exists[v] {
			result = append(result, v)
			exists[v] = true
		}
	}

	return result
}

// SliceUniqueString retrieve unique string items from string slice.
func SliceUniqueString(a []string) []string {
	l := int(len(a)/2) + 1
	result, exists := make([]string, 0, l), make(map[string]bool)
	for _, v := range a {
		if !exists[v] {
			result = append(result, v)
			exists[v] = true
		}
	}

	return result
}

type SliceFilterIntFunc func(s int) bool

// SliceFilterInt 依次将 slice 中的每个值传递到 callback 函数。
// 如果 callback 函数返回 true，则 slice 的当前值会被包含在返回的结果slice中。
// 如果callback 未传 就会过滤零值
func SliceFilterInt(a []int, callbacks ...SliceFilterIntFunc) []int {
	l := int(len(a)/2) + 1
	result := make([]int, 0, l)
	var callback SliceFilterIntFunc
	if len(callbacks) > 0 {
		callback = callbacks[0]
	}
	for _, v := range a {
		if callback == nil {
			if v == 0 {
				continue
			}
		} else {
			if callback(v) == false {
				continue
			}
		}

		result = append(result, v)
	}

	return result
}

type SliceFilterStringFunc func(s string) bool

// SliceFilterString 依次将 slice 中的每个值传递到 callback 函数。
// 如果 callback 函数返回 true，则 slice 的当前值会被包含在返回的结果slice中。
// 如果callback 未传 就会过滤零值
func SliceFilterString(a []string, callbacks ...SliceFilterStringFunc) []string {
	l := int(len(a)/2) + 1
	result := make([]string, 0, l)
	var callback SliceFilterStringFunc
	if len(callbacks) > 0 {
		callback = callbacks[0]
	}
	for _, v := range a {
		if callback == nil {
			if v == "" {
				continue
			}
		} else {
			if callback(v) == false {
				continue
			}
		}

		result = append(result, v)
	}

	return result
}