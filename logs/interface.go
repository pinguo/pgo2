package logs

type IFormatter interface {
    Format(item *LogItem) string
}

type ITarget interface {
    Process(item *LogItem)
    Flush(final bool)
}
