package logs

const (
    LevelNone   = 0x00
    LevelDebug  = 0x01
    LevelInfo   = 0x02
    LevelNotice = 0x04
    LevelWarn   = 0x08
    LevelError  = 0x10
    LevelFatal  = 0x20
    LevelAll    = 0xFF

    rotateNone   = 0
    rotateHourly = 1
    rotateDaily  = 2

    TargetConsole = "console"
    TargetFile    = "file"

    SkipKey = "__"
)
