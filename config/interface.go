package config

type IConfigParser interface {
	Parse(path string) (map[string]interface{}, error)
}

type IConfig interface {
	AddParser(ext string, parser IConfigParser)
	AddPath(path string)
	GetBool(key string, dft bool, dftSplit ...string) bool
	GetInt(key string, dft int, dftSplit ...string) int
	GetFloat(key string, dft float64, dftSplit ...string) float64
	GetString(key string, dft string, dftSplit ...string) string
	GetSliceBool(key string, dftSplit ...string) []bool
	GetSliceInt(key string, dftSplit ...string) []int
	GetSliceFloat(key string, dftSplit ...string) []float64
	GetSliceString(key string, dftSplit ...string) []string
	Get(key string, dftSplit ...string) interface{}
	Set(key string, val interface{}, dftSplit ...string)
	CheckPath() error
}
