package config

type IConfigParser interface {
    Parse(path string) (map[string]interface{}, error)
}

type IConfig interface {
    AddParser(ext string, parser IConfigParser)
    AddPath(path string)
    GetBool(key string, dft bool) bool
    GetInt(key string, dft int) int
    GetFloat(key string, dft float64) float64
    GetString(key string, dft string) string
    GetSliceBool(key string) []bool
    GetSliceInt(key string) []int
    GetSliceFloat(key string) []float64
    GetSliceString(key string) []string
    Get(key string) interface{}
    Set(key string, val interface{})
    CheckPath()
}
