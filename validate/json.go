package validate

import (
    "net/http"

    "github.com/pinguo/pgo2/perror"
    "github.com/pinguo/pgo2/util"
)

// Json validator for json value
type Json struct {
    Name   string
    UseDft bool
    Value  map[string]interface{}
}

func (j *Json) Has(key string) *Json {
    if v := util.MapGet(j.Value, key); !j.UseDft && v == nil {
        panic(perror.New(http.StatusBadRequest, "%s json field missing", j.Name))
    }
    return j
}

func (j *Json) Do() map[string]interface{} {
    return j.Value
}
