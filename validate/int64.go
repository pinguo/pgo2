package validate

import (
	"net/http"

	"github.com/pinguo/pgo2/perror"
)

// int64 validator for int64 value
type Int64 struct {
	Name   string
	UseDft bool
	Value  int64
}

func (i *Int64) Min(v int64) *Int64 {
	if !i.UseDft && i.Value < v {
		panic(perror.New(http.StatusBadRequest, "%s is too small", i.Name))
	}
	return i
}

func (i *Int64) Max(v int64) *Int64 {
	if !i.UseDft && i.Value > v {
		panic(perror.New(http.StatusBadRequest, "%s is too large", i.Name))
	}
	return i
}

func (i *Int64) Enum(enums ...int64) *Int64 {
	found := false
	for _, v := range enums {
		if v == i.Value {
			found = true
			break
		}
	}

	if !i.UseDft && !found {
		panic(perror.New(http.StatusBadRequest, "%s is invalid", i.Name))
	}
	return i
}

func (i *Int64) Do() int64 {
	return i.Value
}

// Int64Slice validator for int64 slice value
type Int64Slice struct {
	Name  string
	Value []int64
}

func (i *Int64Slice) Do() []int64 {
	return i.Value
}
