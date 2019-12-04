package logs

import (
    "fmt"
    "time"
)

func NewLogger(name, logId string, log *Log) *Logger {
    return &Logger{name: name, logId: logId, log: log}
}

// Logger
type Logger struct {
    name  string
    logId string
    log   *Log
}

func (l *Logger) Init(name, logId string, log *Log) {
    l.name, l.logId, l.log = name, logId, log
}

func (l *Logger) logMsg(level int, format string, v ...interface{}) {
    if !l.log.isHandling(level) {
        return
    }

    item := &LogItem{
        When:  time.Now(),
        Level: level,
        Name:  l.name,
        LogId: l.logId,
    }

    if len(v) == 0 {
        item.Message = format
    } else {
        item.Message = fmt.Sprintf(format, v...)
    }

    l.log.addItem(item)
}

func (l *Logger) LogId() string {
    return l.logId
}

func (l *Logger) Debug(format string, v ...interface{}) {
    l.logMsg(LevelDebug, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
    l.logMsg(LevelInfo, format, v...)
}

func (l *Logger) Notice(format string, v ...interface{}) {
    l.logMsg(LevelNotice, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
    l.logMsg(LevelWarn, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
    l.logMsg(LevelError, format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
    l.logMsg(LevelFatal, format, v...)
}
