package pgo2

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    _ "net/http/pprof"
    "os"
    "os/signal"
    "path/filepath"
    "reflect"
    "runtime"
    "sync"
    "sync/atomic"
    "syscall"
    "time"

    "github.com/pinguo/pgo2/core"
    "github.com/pinguo/pgo2/iface"
)

type ServerConfig struct {
    httpAddr  string // address for http
    httpsAddr string // address for https
    debugAddr string // address for pprof

    crtFile         string        // https certificate file
    keyFile         string        // https private key file
    maxHeaderBytes  int           // max http header bytes
    readTimeout     time.Duration // timeout for reading request
    writeTimeout    time.Duration // timeout for writing response
    statsInterval   time.Duration // interval for output server stats
    enableAccessLog bool
    pluginNames     []string
    maxPostBodySize int64 // max post body size
}

// Server the server component, configuration:
// server:
//     httpAddr:  "0.0.0.0:8000"
//     debugAddr: "0.0.0.0:8100"
//     httpsAddr: "0.0.0.0:8443"
//     crtFile: "@app/conf/site.crt"
//     keyFile: "@app/conf/site.key"
//     maxHeaderBytes: 1048576
//     readTimeout:   "30s"
//     writeTimeout:  "30s"
//     statsInterval: "60s"
//     enableAccessLog: true
//     maxPostBodySize: 1048576
func NewServer(config map[string]interface{}) *Server {
    server := &Server{
        maxHeaderBytes:  DefaultHeaderBytes,
        readTimeout:     DefaultTimeout,
        writeTimeout:    DefaultTimeout,
        statsInterval:   60 * time.Second,
        enableAccessLog: true,
    }

    server.pool.New = func() interface{} {
        return new(Context)
    }

    core.Configure(server, config)

    return server
}

type Server struct {
    httpAddr  string // address for http
    httpsAddr string // address for https
    debugAddr string // address for pprof

    crtFile         string        // https certificate file
    keyFile         string        // https private key file
    maxHeaderBytes  int           // max http header bytes
    readTimeout     time.Duration // timeout for reading request
    writeTimeout    time.Duration // timeout for writing response
    statsInterval   time.Duration // interval for output server stats
    enableAccessLog bool
    pluginNames     []string

    numReq          uint64          // request num handled
    plugins         []iface.IPlugin // server plugin list
    servers         []*http.Server  // http server list
    pool            sync.Pool       // Context pool
    maxPostBodySize int64           // max post body size
}

// SetHttpAddr set http addr, if both httpAddr and httpsAddr
// are empty, "0.0.0.0:8000" will be used as httpAddr.
func (s *Server) SetHttpAddr(addr string) {
    s.httpAddr = addr
}

// SetHttpsAddr set https addr.
func (s *Server) SetHttpsAddr(addr string) {
    s.httpsAddr = addr
}

// SetDebugAddr set debug and pprof addr.
func (s *Server) SetDebugAddr(addr string) {
    s.debugAddr = addr
}

// SetCrtFile set certificate file for https
func (s *Server) SetCrtFile(certFile string) {
    s.crtFile, _ = filepath.Abs(GetAlias(certFile))
}

// SetKeyFile set private key file for https
func (s *Server) SetKeyFile(keyFile string) {
    s.keyFile, _ = filepath.Abs(GetAlias(keyFile))
}

// SetMaxHeaderBytes set max header bytes
func (s *Server) SetMaxHeaderBytes(maxBytes int) {
    s.maxHeaderBytes = maxBytes
}

// SetMaxPostBodySize set max header bytes
func (s *Server) SetMaxPostBodySize(maxBytes int64) {
    s.maxPostBodySize = maxBytes
}

// SetReadTimeout set timeout to read request
func (s *Server) SetReadTimeout(v string) {
    if timeout, err := time.ParseDuration(v); err != nil {
        panic(fmt.Sprintf("Server: SetReadTimeout failed, val:%s, err:%s", v, err.Error()))
    } else {
        s.readTimeout = timeout
    }
}

// SetWriteTimeout set timeout to write response
func (s *Server) SetWriteTimeout(v string) {
    if timeout, err := time.ParseDuration(v); err != nil {
        panic(fmt.Sprintf("Server: SetWriteTimeout failed, val:%s, err:%s", v, err.Error()))
    } else {
        s.writeTimeout = timeout
    }
}

// SetStatsInterval set interval to output stats
func (s *Server) SetStatsInterval(v string) {
    if interval, err := time.ParseDuration(v); err != nil {
        panic(fmt.Sprintf("Server: SetStatsInterval failed, val:%s, err:%s", v, err.Error()))
    } else {
        s.statsInterval = interval
    }
}

// SetEnableAccessLog set access log enable or not
func (s *Server) SetEnableAccessLog(v bool) {
    s.enableAccessLog = v
}

// SetPlugins set plugin by names
func (s *Server) SetPlugins(v []interface{}) {
    for _, vv := range v {
        name := vv.(string)
        switch name {
        case "gzip":
            s.plugins = append(s.plugins, NewGzip())
        case "file":
            s.plugins = append(s.plugins, NewFile(App().componentConf("file")))
        default:
            panic("For the defined plug-in:" + name)
        }

    }
}

// AddPlugins add plugin
func (s *Server) AddPlugin(v iface.IPlugin) {
    s.plugins = append(s.plugins, v)
}

// ServerStats server stats
type ServerStats struct {
    MemMB   uint   // memory obtained from os
    NumReq  uint64 // number of handled requests
    NumGO   uint   // number of goroutines
    NumGC   uint   // number of gc runs
    TimeGC  string // total time of gc pause
    TimeRun string // total time of app runs
}

// TimeRun time duration since app run
func (s *Server) timeRun() time.Duration {
    d := time.Since(appTime)
    d -= d % time.Second
    return d
}

// GetStats get server stats
func (s *Server) GetStats() *ServerStats {
    memStats := runtime.MemStats{}
    runtime.ReadMemStats(&memStats)

    timeGC := time.Duration(memStats.PauseTotalNs)
    if timeGC > time.Minute {
        timeGC -= timeGC % time.Second
    } else {
        timeGC -= timeGC % time.Millisecond
    }

    return &ServerStats{
        MemMB:   uint(memStats.Sys / (1 << 20)),
        NumReq:  atomic.LoadUint64(&s.numReq),
        NumGO:   uint(runtime.NumGoroutine()),
        NumGC:   uint(memStats.NumGC),
        TimeGC:  timeGC.String(),
        TimeRun: s.timeRun().String(),
    }
}

// Serve request processing entry
func (s *Server) Serve() {
    // flush log when app end
    defer App().Log().Flush()
    // exec stopBefore when app end
    defer App().StopBefore().Exec()

    // add server plugins
    s.addServerPlugin()

    if App().cmdList() {
        s.cmdList()
        return
    }

    // process command request
    if App().Mode() == ModeCmd {
        s.ServeCMD()
        return
    }

    // process http request
    if s.httpAddr == "" && s.httpsAddr == "" {
        s.httpAddr = DefaultHttpAddr
    }

    wg := sync.WaitGroup{}
    s.handleHttp(&wg)
    s.handleHttps(&wg)
    s.handleDebug(&wg)
    s.handleSignal(&wg)
    s.handleStats(&wg)
    wg.Wait()
}

func (s *Server) cmdList() {
    list := App().Router().CmdHandlers()
    fmt.Println("System parameters:")
    fmt.Println("set running env (requested), eg. --env=online")
    fmt.Println("set running cmd (optional), eg. --cmd=/foo/bar")
    fmt.Println("set base path (optional), eg. --base=/base/path")
    fmt.Println("Displays a list of CMD controllers used (optional), eg. --cmdList")
    fmt.Println("")
    fmt.Println("The path list:")
    for uri, _ := range list {
        fmt.Println("--cmd=" + uri)
    }
}

// ServeCMD serve command request
func (s *Server) ServeCMD() {
    ctx := Context{enableAccessLog: s.enableAccessLog}
    // only apply the last plugin for command
    ctx.Process(s.plugins[len(s.plugins)-1:])
}

// ServeHTTP serve http request
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Change the maxPostBodySize
    if s.maxPostBodySize > 0 {
        r.Body = http.MaxBytesReader(w, r.Body, s.maxPostBodySize)
    }
    // increase request num
    atomic.AddUint64(&s.numReq, 1)

    ctx := s.pool.Get().(iface.IContext)
    ctx.HttpRW(s.enableAccessLog, r, w)
    ctx.Process(s.plugins)
    s.pool.Put(ctx)
}

// HandleRequest handle request of cmd or http,
// this method called in the last of plugin chain.
func (s *Server) HandleRequest(ctx iface.IContext) {
    // get request path and resolve route
    path := ctx.Path()

    // get new controller bind to this route
    rv, action, params := App().Router().CreateController(path, ctx)
    if !rv.IsValid() {
        ctx.End(http.StatusNotFound, []byte("route not found"))
        return
    }

    actionId := ctx.ActionId()
    controller := rv.Interface().(iface.IController)

    // fill empty string for missing param
    numIn := action.Type().NumIn()
    if len(params) < numIn {
        fill := make([]string, numIn-len(params))
        params = append(params, fill...)
    }

    // prepare params for action call
    callParams := make([]reflect.Value, 0)
    for _, param := range params {
        callParams = append(callParams, reflect.ValueOf(param))
    }

    defer func() {
        if v := recover(); v != nil {

          controller.HandlePanic(v)
        }

        // after action hook
        controller.AfterAction(actionId)
    }()

    // before action hook
    controller.BeforeAction(actionId)

    // call action method
    action.Call(callParams)
}

func (s *Server) handleHttp(wg *sync.WaitGroup) {
    if s.httpAddr == "" {
        return
    }

    svr := s.newHttpServer(s.httpAddr)
    s.servers = append(s.servers, svr)
    wg.Add(1)

    GLogger().Info("start running http at " + svr.Addr)

    go func() {
        if err := svr.ListenAndServe(); err != http.ErrServerClosed {
            panic("ListenAndServe failed, " + err.Error())
        }
    }()
}

func (s *Server) handleHttps(wg *sync.WaitGroup) {
    if s.httpsAddr == "" {
        return
    } else if s.crtFile == "" || s.keyFile == "" {
        panic("https no crtFile or keyFile configured")
    }

    svr := s.newHttpServer(s.httpsAddr)
    s.servers = append(s.servers, svr)
    wg.Add(1)

    GLogger().Info("start running https at " + svr.Addr)

    go func() {
        if err := svr.ListenAndServeTLS(s.crtFile, s.keyFile); err != http.ErrServerClosed {
            panic("ListenAndServeTLS failed, " + err.Error())
        }
    }()
}

func (s *Server) handleDebug(wg *sync.WaitGroup) {
    if s.debugAddr == "" {
        return
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        data, _ := json.Marshal(s.GetStats())
        w.Write(data)
    })

    svr := s.newHttpServer(s.debugAddr)
    svr.Handler = nil // use default handler
    s.servers = append(s.servers, svr)
    wg.Add(1)

    GLogger().Info("start running debug at " + svr.Addr)

    go func() {
        if err := svr.ListenAndServe(); err != http.ErrServerClosed {
            panic("ListenAndServe failed, " + err.Error())
        }
    }()
}

func (s *Server) handleSignal(wg *sync.WaitGroup) {
    sig := make(chan os.Signal)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sig // wait signal
        for _, svr := range s.servers {
            GLogger().Info("stop running " + svr.Addr)
            ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
            svr.Shutdown(ctx)
            wg.Done()
        }
    }()
}

func (s *Server) handleStats(wg *sync.WaitGroup) {
    timer := time.Tick(s.statsInterval)

    go func() {
        for {
            <-timer // wait timer
            data, _ := json.Marshal(s.GetStats())
            GLogger().Info("app stats: " + string(data))
        }
    }()
}

func (s *Server) newHttpServer(addr string) *http.Server {
    return &http.Server{
        Addr:           addr,
        ReadTimeout:    s.readTimeout,
        WriteTimeout:   s.writeTimeout,
        MaxHeaderBytes: s.maxHeaderBytes,
        Handler:        s,
    }
}

func (s *Server) addServerPlugin() {
    // server is the last plugin
    s.AddPlugin(s)

    if len(s.plugins) > MaxPlugins {
        panic("Server: too many plugins")
    }
}
