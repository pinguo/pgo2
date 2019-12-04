package render

import (
    "encoding/xml"
    "fmt"
    "net/http"
)

func NewXml(data interface{}) Render {
    return &Xml{data: data, httpCode: http.StatusOK}
}

type Xml struct {
    data        interface{}
    httpCode    int
    contentType string
}

func (x *Xml) SetHttpCode(code int) {
    x.httpCode = code
}

func (x *Xml) Content() []byte {
    output, e := xml.Marshal(x.data)

    if e != nil {
        panic(fmt.Sprintf("failed to marshal xml, %s", e))
    }

    return output
}

func (x *Xml) HttpCode() int {
    return x.httpCode
}

func (x *Xml) ContentType() string {
    return "application/xml; charset=utf-8"
}
