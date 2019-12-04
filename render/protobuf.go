package render

import (
    "fmt"
    "net/http"

    "github.com/golang/protobuf/proto"
)

func NewProtoBuf(data interface{}) Render {
    return &ProtoBuf{data: data, httpCode: http.StatusOK}
}

type ProtoBuf struct {
    data     interface{}
    httpCode int
}

func (p *ProtoBuf) SetHttpCode(code int) {
    p.httpCode = code
}

func (p *ProtoBuf) Content() []byte {
    bytes, e := proto.Marshal(p.data.(proto.Message))

    if e != nil {
        panic(fmt.Sprintf("failed to marshal ProtoBuf, %s", e))
    }

    return bytes
}

func (p *ProtoBuf) HttpCode() int {
    return p.httpCode
}

func (p *ProtoBuf) ContentType() string {
    return "application/x-protobuf"
}
