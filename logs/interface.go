package logs

type IFormatter interface {
	Format(item *LogItem) string
}

type ITarget interface {
	SetLevels(v interface{})
	SetFormatter(v interface{})
	IsHandling(level int) bool
	Format(item *LogItem) string
	Process(item *LogItem)
	Flush(final bool)
}

type ILogger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Notice(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}