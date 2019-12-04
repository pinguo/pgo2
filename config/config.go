package config

import (
    "os"
    "path/filepath"
    "strings"
    "sync"

    "github.com/pinguo/pgo2/util"
)

func New(basePath, env string) *Config {
    config := &Config{
        basePath: basePath,
        env:      env,
        parsers:  make(map[string]IConfigParser),
        data:     make(map[string]interface{}),
        paths:    make([]string, 0),
    }

    config.Init()

    return config
}

// Config the config component
type Config struct {
    parsers  map[string]IConfigParser
    data     map[string]interface{}
    paths    []string
    lock     sync.RWMutex
    basePath string
    env      string
}

// Initialize the
func (c *Config) Init() {
    confPath := filepath.Join(c.basePath, "conf")
    if f, _ := os.Stat(confPath); f != nil && f.IsDir() {
        c.paths = append(c.paths, confPath)
    }

    envPath := filepath.Join(confPath, c.env)
    if f, _ := os.Stat(envPath); f != nil && f.IsDir() {
        c.paths = append(c.paths, envPath)
    }

    c.AddParser("json", &JsonParser{})
    c.AddParser("yaml", &YamlParser{})
}

func (c *Config) CheckPath() {
    confPath := filepath.Join(c.basePath, "conf")
    if f, _ := os.Stat(confPath); f == nil || !f.IsDir() {
        panic("Config: invalid conf path, " + confPath)
    }

    envPath := filepath.Join(confPath, c.env)
    if f, _ := os.Stat(envPath); f == nil || !f.IsDir() {
        panic("Config: invalid env path, " + envPath)
    }
}

// AddParser add parser for file with ext extension
func (c *Config) AddParser(ext string, parser IConfigParser) {
    c.parsers[ext] = parser
}

// AddPath add path to end of search paths
func (c *Config) AddPath(path string) {
    paths := make([]string, 0)
    for _, v := range c.paths {
        if v != path {
            paths = append(paths, v)
        }
    }

    c.paths = append(paths, path)
}

// GetBool get bool value from config,
// key is dot separated config key,
// dft is default value if key not exists.
func (c *Config) GetBool(key string, dft bool) bool {
    if v := c.Get(key); v != nil {
        return util.ToBool(v)
    }

    return dft
}

// GetInt get int value from config,
// key is dot separated config key,
// dft is default value if key not exists.
func (c *Config) GetInt(key string, dft int) int {
    if v := c.Get(key); v != nil {
        return util.ToInt(v)
    }

    return dft
}

// GetFloat get float value from config,
// key is dot separated config key,
// dft is default value if key not exists.
func (c *Config) GetFloat(key string, dft float64) float64 {
    if v := c.Get(key); v != nil {
        return util.ToFloat(v)
    }

    return dft
}

// GetString get string value from config,
// key is dot separated config key,
// dft is default value if key not exists.
func (c *Config) GetString(key string, dft string) string {
    if v := c.Get(key); v != nil {
        return util.ToString(v)
    }

    return dft
}

// GetSliceBool get []bool value from config,
// key is dot separated config key,
// nil is default value if key not exists.
func (c *Config) GetSliceBool(key string) []bool {
    var ret []bool
    if v := c.Get(key); v != nil {
        if vI, ok := v.([]interface{}); ok == true {
            for _, vv := range vI {
                ret = append(ret, util.ToBool(vv))
            }
        }
    }

    return ret
}

// GetSliceInt get []int value from config,
// key is dot separated config key,
// nil is default value if key not exists.
func (c *Config) GetSliceInt(key string) []int {
    var ret []int
    if v := c.Get(key); v != nil {
        if vI, ok := v.([]interface{}); ok == true {
            for _, vv := range vI {
                ret = append(ret, util.ToInt(vv))
            }
        }
    }

    return ret
}

// GetSliceFloat get []float value from config,
// key is dot separated config key,
// nil is default value if key not exists.
func (c *Config) GetSliceFloat(key string) []float64 {
    var ret []float64
    if v := c.Get(key); v != nil {
        if vI, ok := v.([]interface{}); ok == true {
            for _, vv := range vI {
                ret = append(ret, util.ToFloat(vv))
            }
        }
    }

    return ret
}

// GetSliceString get []string value from config,
// key is dot separated config key,
// nil is default value if key not exists.
func (c *Config) GetSliceString(key string) []string {
    var ret []string
    if v := c.Get(key); v != nil {
        if vI, ok := v.([]interface{}); ok == true {
            for _, vv := range vI {
                ret = append(ret, util.ToString(vv))
            }
        }
    }

    return ret
}

// Get get value by dot separated key,
// the first part of key is file name
// without extension. if key is empty,
// all loaded config will be returned.
func (c *Config) Get(key string) interface{} {
    ks := strings.Split(key, ".")
    if _, ok := c.data[ks[0]]; !ok {
        c.load(ks[0])
    }

    c.lock.RLock()
    defer c.lock.RUnlock()

    return util.MapGet(c.data, key)
}

// Set set value by dot separated key,
// if key is empty, the value will set
// to root, if val is nil, the key will
// be deleted.
func (c *Config) Set(key string, val interface{}) {
    c.lock.Lock()
    defer c.lock.Unlock()

    util.MapSet(c.data, key, val)
}

// Load load config file under the search paths.
// file under env sub path will be merged.
func (c *Config) load(name string) {
    c.lock.Lock()
    defer c.lock.Unlock()

    // avoid repeated loading
    _, ok := c.data[name]
    if ok || len(name) == 0 {
        return
    }

    for _, path := range c.paths {
        files, _ := filepath.Glob(filepath.Join(path, name+".*"))
        for _, f := range files {
            ext := strings.ToLower(filepath.Ext(f))
            if parser, ok := c.parsers[ext[1:]]; ok {
                conf, err := parser.Parse(f)
                if err != nil {
                    panic(err.Error())
                }
                if conf != nil {
                    util.MapMerge(c.data, map[string]interface{}{name: conf})
                }
            }
        }
    }
}
