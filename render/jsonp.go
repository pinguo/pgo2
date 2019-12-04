package render

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

func NewJsonp(callback string, data interface{}) Render {
    return &Jsonp{callback: callback, data: data, httpCode: http.StatusOK}
}

type Jsonp struct {
    callback string
    data     interface{}
    httpCode int
}

func (j *Jsonp) SetHttpCode(code int) {
    j.httpCode = code
}

func (j *Jsonp) Content() []byte {
    output, e := json.Marshal(j.data)

    if e != nil {
        panic(fmt.Sprintf("failed to marshal json, %s", e))
    }

    buf := &bytes.Buffer{}
    buf.WriteString(j.callback + "(")
    buf.Write(output)
    buf.WriteString(")")

    return buf.Bytes()
}

func (j *Jsonp) HttpCode() int {
    return j.httpCode
}

func (j *Jsonp) ContentType() string {
    return "application/json; charset=utf-8"
}
