package logs

import (
    "fmt"
    "runtime"
    "strings"
    "sync"
    "time"

    "github.com/pinguo/pgo2/core"
)

// LevelToString convert int level to string
func LevelToString(level int) string {
    switch level {
    case LevelNone:
        return "NONE"
    case LevelDebug:
        return "DEBUG"
    case LevelInfo:
        return "INFO"
    case LevelNotice:
        return "NOTICE"
    case LevelWarn:
        return "WARN"
    case LevelError:
        return "ERROR"
    case LevelFatal:
        return "FATAL"
    case LevelAll:
        return "ALL"
    default:
        panic(fmt.Sprintf("unknown log level: %x", level))
    }
}

// StringToLevel convert string to int level
func StringToLevel(level string) int {
    switch strings.ToUpper(level) {
    case "NONE":
        return LevelNone
    case "DEBUG":
        return LevelDebug
    case "INFO":
        return LevelInfo
    case "NOTICE":
        return LevelNotice
    case "WARN":
        return LevelWarn
    case "ERROR":
        return LevelError
    case "FATAL":
        return LevelFatal
    case "ALL":
        return LevelAll
    default:
        panic(fmt.Sprintf("unknown log level: %s", level))
    }
}

// parse comma separated level string to int format
// eg. `debug,info` => 0x03
func parseLevels(str string) int {
    levels := LevelNone
    parts := strings.Split(str, ",")

    for _, v := range parts {
        v = strings.TrimSpace(v)
        levels |= StringToLevel(v)
    }

    return levels
}

// LogItem represent an item of log
type LogItem struct {
    When    time.Time
    Level   int
    Name    string
    LogId   string
    Trace   string
    Message string
}

// Log the log component, configuration:
// log:
//     levels: "ALL"
//     traceLevels: "DEBUG"
//     chanLen: 1000
//     flushInterval: "60s"
//     targets:
//         info:
//             name: "file"
//             levels: "DEBUG,INFO,NOTICE"
//             filePath: "@runtime/info.log"
//             maxLogFile: 10
//             rotate: "daily"
//         error: {
//             name: "file"
//             levels: "WARN,ERROR,FATAL"
//             filePath: "@runtime/error.log"
//             maxLogFile: 10
//             rotate: "daily"
func NewLog(runtimePath string, config map[string]interface{}) *Log {
    log := &Log{
        levels:        LevelAll,
        chanLen:       1000,
        traceLevels:   LevelDebug,
        flushInterval: 60 * time.Second,
        runtimePath:   runtimePath,
    }

    core.Configure(log, config)

    log.Init()

    return log
}

type Log struct {
    levels        int
    chanLen       int
    traceLevels   int
    flushInterval time.Duration
    targets       map[string]ITarget
    msgChan       chan *LogItem
    msgChanClosed bool
    wg            sync.WaitGroup

    runtimePath string
}

func (d *Log) Init() {
    d.msgChan = make(chan *LogItem, d.chanLen)

    if len(d.targets) == 0 {
        // use console target as default
        d.targets = make(map[string]ITarget)
        d.targets[TargetConsole] = NewConsole()
    }

    // start loop
    d.wg.Add(1)
    go d.loop()
}

func (d *Log) Target(name string) ITarget {
    return d.targets[name]
}

// SetLevels set levels to handle, default "ALL"
func (d *Log) SetLevels(v interface{}) {
    if _, ok := v.(string); ok {
        d.levels = parseLevels(v.(string))
    } else if _, ok := v.(int); ok {
        d.levels = v.(int)
    } else {
        panic(fmt.Sprintf("Log: invalid levels: %v", v))
    }
}

// SetChanLen set length of log channel, default 1000
func (d *Log) SetChanLen(len int) {
    d.chanLen = len
}

// SetTraceLevels set levels to trace, default "DEBUG"
func (d *Log) SetTraceLevels(v interface{}) {
    if _, ok := v.(string); ok {
        d.traceLevels = parseLevels(v.(string))
    } else if _, ok := v.(int); ok {
        d.traceLevels = v.(int)
    } else {
        panic(fmt.Sprintf("Log: invalid trace levels: %v", v))
    }
}

// SetFlushInterval set interval to flush log, default "60s"
func (d *Log) SetFlushInterval(v string) {
    if flushInterval, err := time.ParseDuration(v); err != nil {
        panic(fmt.Sprintf("Log: parse flushInterval error, val:%s, err:%s", v, err.Error()))
    } else {
        d.flushInterval = flushInterval
    }
}

// SetTargets set output target, ConsoleTarget will be used if no targets specified
func (d *Log) SetTargets(targets map[string]interface{}) {
    d.targets = make(map[string]ITarget)

    for name, val := range targets {
        config, ok := val.(map[string]interface{})
        if ok {
            if _, ok := config["name"]; !ok {
                config["name"] = TargetConsole
            }
        }

        class, _ := config["name"].(string)
        switch class {
        case TargetConsole:
            d.targets[name] = NewConsole(config)
        case TargetFile:
            d.targets[name] = NewFile(d.runtimePath, config)

        }

    }
}

// SetTarget set output target
func (d *Log) SetTarget(name string, target ITarget) {
    d.targets[name] = target
}

// GetLogger get a new logger with name and id specified
func (d *Log) Logger(name, logId string) *Logger {
    return NewLogger(name, logId, d)
}

// GetProfiler get a new profiler
func (d *Log) Profiler() *Profiler {
    return NewProfiler()
}

// Flush close msg chan and wait loop end
func (d *Log) Flush() {
    d.msgChanClosed = true
    close(d.msgChan)
    d.wg.Wait()
}

func (d *Log) isHandling(level int) bool {
    return true
    return level&d.levels != 0
}

func (d *Log) addItem(item *LogItem) {
    if d.traceLevels&item.Level != 0 {
        if _, file, line, ok := runtime.Caller(3); ok {
            if pos := strings.LastIndex(file, "src/"); pos > 0 {
                file = file[pos+4:]
            }

            item.Trace = fmt.Sprintf("[%s:%d]", file, line)
        }
    }

    if d.msgChanClosed {
        // msgChan is closed
        return
    }

    d.msgChan <- item
}

func (d *Log) loop() {
    flushTimer := time.Tick(d.flushInterval)

    for {
        select {
        case item, ok := <-d.msgChan:
            for _, target := range d.targets {
                if ok {
                    target.Process(item)
                } else {
                    target.Flush(true)
                }
            }

            if !ok {
                goto end
            }
        case <-flushTimer:
            for _, target := range d.targets {
                target.Flush(false)
            }
        }
    }

end:
    d.wg.Done()
}
