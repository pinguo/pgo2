package pgo2

import (
    "net/http"
    "os"
    "path/filepath"

    "github.com/pinguo/pgo2/core"
    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/util"
)

// File file plugin, this plugin only handle file in @public directory,
// request url with empty or excluded extension will not be handled.
func NewFile(config map[string]interface{}) *File {
    f := &File{}

    core.Configure(f, config)

    return f
}

type File struct {
    excludeExtensions []string
}

func (f *File) SetExcludeExtensions(v []interface{}) {
    for _, vv := range v {
        f.excludeExtensions = append(f.excludeExtensions, vv.(string))
    }
}

func (f *File) HandleRequest(ctx iface.IContext) {
    // if extension is empty or excluded, pass
    if ext := filepath.Ext(ctx.Path()); ext == "" {
        return
    } else if len(f.excludeExtensions) != 0 {
        if util.SliceSearchString(f.excludeExtensions, ext) != -1 {
            return
        }
    }

    // skip other plugins
    defer ctx.Abort()

    // GET or HEAD method is required
    method := ctx.Method()
    if method != http.MethodGet && method != http.MethodHead {
        http.Error(ctx.Output(), "", http.StatusMethodNotAllowed)
        return
    }

    // file in @public directory is required
    path := filepath.Join(App().PublicPath(), util.CleanPath(ctx.Path()))
    h, e := os.Open(path)
    if e != nil {
        http.Error(ctx.Output(), "", http.StatusNotFound)
        return
    }

    defer h.Close()
    info, _ := h.Stat()
    http.ServeContent(ctx.Output(), ctx.Input(), path, info.ModTime(), h)
}
