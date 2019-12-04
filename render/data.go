package render

import (
    "net/http"
)

func NewData(data []byte) Render {
    return &Data{data: data, httpCode: http.StatusOK}
}

type Data struct {
    data     []byte
    httpCode int
}

func (d *Data) SetHttpCode(code int) {
    d.httpCode = code
}

func (d *Data) Content() []byte {
    return d.data
}

func (d *Data) HttpCode() int {
    return d.httpCode
}

func (d *Data) ContentType() string {
    return ""
}
