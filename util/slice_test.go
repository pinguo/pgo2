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


func TestSliceFilterInt(t *testing.T) {


	t.Run("callback=empty", func(t *testing.T) {
		list := []int{0,2}
		if len(SliceFilterInt(list)) != 1{
			t.Fatal(`len(SliceFilterInt(list)) != 1`)
		}
	})

	t.Run("callback=filter2", func(t *testing.T) {
		list := []int{1,0,2}
		if len(SliceFilterInt(list, func(s int) bool {
			if s==1 {
				return false
			}
			return true
		})) != 2{
			t.Fatal(`len(SliceFilterInt(list)) != 2`)
		}
	})

}

func TestSliceFilterString(t *testing.T) {


	t.Run("callback=empty", func(t *testing.T) {
		list := []string{"","2"}
		if len(SliceFilterString(list)) != 1{
			t.Fatal(`len(SliceFilterInt(list)) != 1`)
		}
	})

	t.Run("callback=filter2", func(t *testing.T) {
		list := []string{"1","","2"}
		if len(SliceFilterString(list, func(s string) bool {
			if s=="1" {
				return false
			}
			return true
		})) != 2{
			t.Fatal(`len(SliceFilterString(list)) != 2`)
		}
	})

	t.Run("list=empty", func(t *testing.T) {
		list := []string{}
		if len(SliceFilterString(list)) != 0{
			t.Fatal(`len(SliceFilterInt(list)) != 0`)
		}
	})

}