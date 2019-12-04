package validate

import (
    "net/http"

    "github.com/pinguo/pgo2/perror"
)

// Int validator for int value
type Int struct {
    Name   string
    UseDft bool
    Value  int
}

func (i *Int) Min(v int) *Int {
    if !i.UseDft && i.Value < v {
        panic(perror.New(http.StatusBadRequest, "%s is too small", i.Name))
    }
    return i
}

func (i *Int) Max(v int) *Int {
    if !i.UseDft && i.Value > v {
        panic(perror.New(http.StatusBadRequest, "%s is too large", i.Name))
    }
    return i
}

func (i *Int) Enum(enums ...int) *Int {
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

func (i *Int) Do() int {
    return i.Value
}

// IntSlice validator for int slice value
type IntSlice struct {
    Name  string
    Value []int
}

func (i *IntSlice) Do() []int {
    return i.Value
}
