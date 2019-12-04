package util

import "testing"

func TestSliceSearchInt(t *testing.T) {
    list := []int{1, 1, 2}
    if SliceSearchInt(list, 1) != 0 {
        t.Fatal(`SliceSearchInt(list,1) !=0`)
    }

    if SliceSearchInt(list, 3) != -1 {
        t.Fatal(`SliceSearchInt(list,3) !=-1`)
    }
}

func TestSliceSearchString(t *testing.T) {
    list := []string{"a", "a", "c"}
    if SliceSearchString(list, "c") != 2 {
        t.Fatal(`SliceSearchString(list,"c") !=2`)
    }

    if SliceSearchString(list, "d") != -1 {
        t.Fatal(SliceSearchString(list, "d") != -1)
    }
}

func TestSliceUniqueInt(t *testing.T) {
    list := []int{1, 1, 2}
    if len(SliceUniqueInt(list)) != 2 {
        t.FailNow()
    }
}

func TestSliceUniqueString(t *testing.T) {
    list := []string{"a", "a", "c"}
    if len(SliceUniqueString(list)) != 2 {
        t.FailNow()
    }
}
