package validate

import (
    "net/http"

    "github.com/pinguo/pgo2/perror"
)

// Bool validator for bool value
type Bool struct {
    Name   string
    UseDft bool
    Value  bool
}

func (b *Bool) Must(v bool) *Bool {
    if !b.UseDft && b.Value != v {
        panic(perror.New(http.StatusBadRequest, "%s must be %v", b.Name, v))
    }
    return b
}

func (b *Bool) Do() bool {
    return b.Value
}
