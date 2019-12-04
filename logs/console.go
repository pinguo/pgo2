package logs

import (
    "os"

    "github.com/pinguo/pgo2/core"
)

func NewConsole(dftConfig ...map[string]interface{}) ITarget {
    target := &Console{}
    target.levels = LevelAll

    dftConfig = append(dftConfig, make(map[string]interface{}))
    config := dftConfig[0]
    core.Configure(target, config)

    return target
}

// ConsoleTarget target for console
type Console struct {
    Target
}

// Process write log to stdout
func (c *Console) Process(item *LogItem) {
    if !c.IsHandling(item.Level) {
        return
    }

    os.Stdout.WriteString(c.Format(item))
    return
}

// Flush flush log to stdout
func (c *Console) Flush(final bool) {
    os.Stdout.Sync()
}
