package logs

type IFormatter interface {
	Format(item *LogItem) string
}

type ITarget interface {
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