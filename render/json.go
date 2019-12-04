package render

import (
    "encoding/json"
    "fmt"
    "net/http"
)

func NewJson(data interface{}) Render {
    return &Json{data: data, httpCode: http.StatusOK}
}

type Json struct {
    data     interface{}
    httpCode int
}

func (j *Json) SetHttpCode(code int) {
    j.httpCode = code
}

func (j *Json) Content() []byte {
    output, e := json.Marshal(j.data)

    if e != nil {
        panic(fmt.Sprintf("failed to marshal json, %s", e))
    }

    return output
}

func (j *Json) HttpCode() int {
    return j.httpCode
}

func (j *Json) ContentType() string {
    return "application/json; charset=utf-8"
}
