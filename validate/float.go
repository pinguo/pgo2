package validate

import (
    "net/http"

    "github.com/pinguo/pgo2/perror"
)

// Float validator for float value
type Float struct {
    Name   string
    UseDft bool
    Value  float64
}

func (f *Float) Min(v float64) *Float {
    if !f.UseDft && f.Value < v {
        panic(perror.New(http.StatusBadRequest, "%s is too small", f.Name))
    }
    return f
}

func (f *Float) Max(v float64) *Float {
    if !f.UseDft && f.Value > v {
        panic(perror.New(http.StatusBadRequest, "%s is too large", f.Name))
    }
    return f
}

func (f *Float) Do() float64 {
    return f.Value
}

// FloatSliceValidator validator for float slice value
type FloatSlice struct {
    Name  string
    Value []float64
}

func (f *FloatSlice) Do() []float64 {
    return f.Value
}
