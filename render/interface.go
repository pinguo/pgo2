package render

type Render interface {
    // Data() Data to be written
    Content() []byte
    // HttpCode The HTTP status code
    HttpCode() int
    // ContentType The HTTP Content-Type
    ContentType() string
    // SetHttpCode Set The HTTP status code
    SetHttpCode(code int)
}
