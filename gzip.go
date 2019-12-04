package pgo2

import (
    "compress/gzip"
    "io"
    "io/ioutil"
    "net/http"
    "path/filepath"
    "strings"
    "sync"

    "github.com/pinguo/pgo2/iface"
)

// Gzip gzip compression plugin
func NewGzip() *Gzip {
    gzipObj := &Gzip{}
    gzipObj.pool.New = func() interface{} {
        return &gzipWrite{writer: gzip.NewWriter(ioutil.Discard)}
    }

    return gzipObj
}

type Gzip struct {
    pool sync.Pool
}

func (g *Gzip) HandleRequest(ctx iface.IContext) {
    ae := ctx.Header("Accept-Encoding", "")
    if !strings.Contains(ae, "gzip") {
        return
    }

    ext := filepath.Ext(ctx.Path())
    switch strings.ToLower(ext) {
    case ".png", ".gif", ".jpeg", ".jpg", ".ico":
        return
    }

    gw := g.pool.Get().(*gzipWrite)
    gw.reset(ctx)

    defer func() {
        gw.finish()
        g.pool.Put(gw)
    }()

    ctx.Next()
}

type gzipWrite struct {
    http.ResponseWriter
    writer *gzip.Writer
    ctx    iface.IContext
    size   int
}

func (g *gzipWrite) reset(ctx iface.IContext) {
    g.ResponseWriter = ctx.Output()
    g.ctx = ctx
    g.size = -1
    ctx.SetOutput(g)
}

func (g *gzipWrite) finish() {
    if g.size > 0 {
        g.writer.Close()
    }
}

func (g *gzipWrite) start() {
    if g.size == -1 {
        g.size = 0
        g.writer.Reset(g.ResponseWriter)
        g.ctx.SetHeader("Content-Encoding", "gzip")
    }
}

func (g *gzipWrite) Flush() {
    if g.size > 0 {
        g.writer.Flush()
    }

    if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
        flusher.Flush()
    }
}

func (g *gzipWrite) Write(data []byte) (n int, e error) {
    if len(data) == 0 {
        return 0, nil
    }

    g.start()

    n, e = g.writer.Write(data)
    g.size += n
    return
}

func (g *gzipWrite) WriteString(data string) (n int, e error) {
    if len(data) == 0 {
        return 0, nil
    }

    g.start()

    n, e = io.WriteString(g.writer, data)
    g.size += n
    return
}
