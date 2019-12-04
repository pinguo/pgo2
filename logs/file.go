package logs

import (
    "bytes"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/pinguo/pgo2/core"
)

// File target for file, configuration:
// info:
//     class: "@pgo/logs/File"
//     levels: "DEBUG,INFO,NOTICE"
//     filePath: "@runtime/info.log"
//     maxLogFile: 10
//     maxBufferByte: 10485760
//     maxBufferLine: 10000
//     rotate: "daily"
func NewFile(runtimePath string, dftConfig ...map[string]interface{}) ITarget {
    file := &File{
        filePath:      "@runtime/app.log",
        maxLogFile:    10,
        maxBufferByte: 10 * 1024 * 1024,
        maxBufferLine: 10000,
        rotate:        rotateDaily,
        runtimePath:   runtimePath,
    }

    file.levels = LevelAll

    dftConfig = append(dftConfig, make(map[string]interface{}))
    config := dftConfig[0]
    core.Configure(file, config)

    return file

}

type File struct {
    Target
    filePath      string
    maxLogFile    int
    maxBufferByte int
    maxBufferLine int
    rotate        int

    buffer        bytes.Buffer
    lastRotate    time.Time
    curBufferLine int

    runtimePath string
}

func (f *File) Init() {
    f.filePath = f.completePath()
    h, e := os.OpenFile(f.filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
    if e != nil {
        panic(fmt.Sprintf("File: failed to open file: %s, e: %s", f.filePath, e))
    }

    defer h.Close()

    stat, e := h.Stat()
    if e != nil {
        panic(fmt.Sprintf("File: failed to stat file: %s, e: %s", f.filePath, e))
    }

    f.curBufferLine = 0
    f.lastRotate = stat.ModTime()
    f.buffer.Grow(f.maxBufferByte)
}

// SetFilePath set file path, default "@runtime/app.log"
func (f *File) SetFilePath(filePath string) {
    f.filePath = filePath
}

// SetMaxLogFile set max log backups, default 10
func (f *File) SetMaxLogFile(maxLogFile int) {
    f.maxLogFile = maxLogFile
}

// SetMaxBufferByte set max buffer bytes, default 10MB
func (f *File) SetMaxBufferByte(maxBufferByte int) {
    f.maxBufferByte = maxBufferByte
}

// SetMaxBufferLine set max buffer lines, default 10000
func (f *File) SetMaxBufferLine(maxBufferLine int) {
    f.maxBufferLine = maxBufferLine
}

// SetRotate set rotate policy(none, hourly, daily), default "daily"
func (f *File) SetRotate(rotate string) {
    switch strings.ToUpper(rotate) {
    case "NONE":
        f.rotate = rotateNone
    case "HOURLY":
        f.rotate = rotateHourly
    case "DAILY":
        f.rotate = rotateDaily
    default:
        panic("File: invalid rotate:" + rotate)
    }
}

//  resolve path, eg. @runtime/app.log => /path/to/runtime/app.log
func (f *File) completePath() string {
    if strings.Index(f.filePath, "@runtime") >= 0 {
        return strings.Replace(f.filePath, "@runtime", f.runtimePath, 1)
    }

    return f.filePath
}

// Process check and rotate log file if rotate is enable,
// write log to buffer, flush buffer to file if buffer is full.
func (f *File) Process(item *LogItem) {
    if !f.IsHandling(item.Level) {
        return
    }

    // rotate log file
    if f.shouldRotate(item.When) {
        f.rotateLog(item.When)
    }

    // write log to buffer
    f.buffer.WriteString(f.Format(item))
    f.curBufferLine++

    // flush buffer to file
    if f.curBufferLine >= f.maxBufferLine || f.buffer.Len() >= f.maxBufferByte {
        f.Flush(false)
    }

    return
}

// Flush flush log buffer to file
func (f *File) Flush(final bool) {
    // nothing to flush
    if f.curBufferLine == 0 {
        return
    }

    // open log file to write
    h, e := os.OpenFile(f.filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
    if e != nil {
        panic(fmt.Sprintf("File: failed to open file: %s, e: %s", f.filePath, e))
    }

    defer h.Close()

    // write log buffer to file
    f.buffer.WriteTo(h)
    f.buffer.Reset()
    f.curBufferLine = 0
}

func (f *File) shouldRotate(now time.Time) bool {
    if f.rotate == rotateHourly {
        return now.Hour() != f.lastRotate.Hour() || now.Day() != f.lastRotate.Day()
    } else if f.rotate == rotateDaily {
        return now.Day() != f.lastRotate.Day()
    }

    return false
}

func (f *File) rotateLog(now time.Time) {
    layout, interval := "", time.Duration(0)
    if f.rotate == rotateHourly {
        layout = "2006010215"
        interval = time.Hour
    } else if f.rotate == rotateDaily {
        layout = "20060102"
        interval = time.Hour * 24
    } else {
        return
    }

    // flush and close file
    f.Flush(true)

    // move current file to backup file
    suffix := f.lastRotate.Format(layout)
    newPath := fmt.Sprintf("%s.%s", f.filePath, suffix)
    os.Rename(f.filePath, newPath)

    // update last rotate time
    f.lastRotate = now

    // clean backup file
    backups, _ := filepath.Glob(f.filePath + ".*")
    if len(backups) > 0 {
        for _, backup := range backups {
            ext := filepath.Ext(backup)
            d, e := time.ParseInLocation(layout, ext[1:], now.Location())
            if e == nil && int(now.Sub(d)/interval) > f.maxLogFile {
                os.Remove(backup)
            }
        }
    }
}
