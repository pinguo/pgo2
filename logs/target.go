package logs

import "fmt"

// Target base class of output
type Target struct {
    levels    int
    formatter IFormatter
}

// SetLevels set levels for target, eg. "DEBUG,INFO,NOTICE"
func (t *Target) SetLevels(v interface{}) {
    if _, ok := v.(string); ok {
        t.levels = parseLevels(v.(string))
    } else if _, ok := v.(int); ok {
        t.levels = v.(int)
    } else {
        panic(fmt.Sprintf("Target: invalid levels: %v", v))
    }
}

// SetFormatter set user-defined log formatter, eg. "Lib/Log/Formatter"
func (t *Target) SetFormatter(v interface{}) {
    if ptr, ok := v.(IFormatter); ok {
        t.formatter = ptr
        //} else if class, ok := v.(string); ok {
        //    t.formatter = CreateObject(class).(IFormatter)
        //} else if config, ok := v.(map[string]interface{}); ok {
        //    t.formatter = CreateObject(config).(IFormatter)
    } else {
        panic(fmt.Sprintf("Target: invalid formatter: %v", v))
    }
}

// IsHandling check whether this target is handling the log item
func (t *Target) IsHandling(level int) bool {
    return t.levels&level != 0
}

// Format format log item to string
func (t *Target) Format(item *LogItem) string {
    // call user-defined formatter if exists
    if t.formatter != nil {
        return t.formatter.Format(item)
    }

    // default log format: [time][logId][name][level][trace]: message\n
    return fmt.Sprintf("[%s][%s][%s][%s]%s: %s\n",
        item.When.Format("2006/01/02 15:04:05.000"),
        item.LogId,
        item.Name,
        LevelToString(item.Level),
        item.Trace,
        item.Message,
    )
}
