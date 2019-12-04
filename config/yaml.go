package config

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"

    "github.com/go-yaml/yaml"
)

// YamlParser parser for yaml config
type YamlParser struct {
    base
}

// Parse parse yaml config, environment value like ${env||default} will expand
func (y *YamlParser) Parse(path string) (parseData map[string]interface{}, err error) {
    h, e := os.Open(path)
    if e != nil {
        return nil, fmt.Errorf("YamlParser: failed to open file: " + path)
    }

    defer h.Close()

    content, e := ioutil.ReadAll(h)
    if e != nil {
        return nil, fmt.Errorf("YamlParser: failed to read file: " + path)
    }

    // expand env: ${env||default}
    content = y.expandEnv(content)

    if e := y.YamlUnmarshal(content, &parseData); e != nil {
        panic(fmt.Sprintf("YamlParser: failed to parse file: %s, %s", path, e.Error()))
    }

    return parseData, nil
}

// YamlMarshal wrapper for yaml.Marshal.
func (y *YamlParser) YamlMarshal(in interface{}) ([]byte, error) {
    return yaml.Marshal(in)
}

// YamlUnmarshal wrapper for yaml.Unmarshal or yaml.UnmarshalStrict,
// if type of out is map[string]interface{}, *map[string]interface{},
// the inner map[interface{}]interface{} will fix to map[string]interface{}
// recursively. if type of out is *interface{}, the underlying type of
// out will change to *map[string]interface{}.
func (y *YamlParser) YamlUnmarshal(in []byte, out interface{}, strict ...bool) error {
    var err error
    if len(strict) > 0 && strict[0] {
        err = yaml.UnmarshalStrict(in, out)
    } else {
        err = yaml.Unmarshal(in, out)
    }

    if err == nil {
        y.yamlFixOut(out)
    }

    return err
}

// YamlEncode wrapper for yaml.Encoder.
func (y *YamlParser) YamlEncode(w io.Writer, in interface{}) error {
    enc := yaml.NewEncoder(w)
    defer enc.Close()
    return enc.Encode(in)
}

// YamlDecode wrapper for yaml.Decoder, strict is for Decoder.SetStrict().
// if type of out is map[string]interface{}, *map[string]interface{},
// the inner map[interface{}]interface{} will fix to map[string]interface{}
// recursively. if type of out is *interface{}, the underlying type of
// out will change to *map[string]interface{}.
func (y *YamlParser) YamlDecode(r io.Reader, out interface{}, strict ...bool) error {
    dec := yaml.NewDecoder(r)
    if len(strict) > 0 && strict[0] {
        dec.SetStrict(true)
    }

    err := dec.Decode(out)
    if err == nil {
        y.yamlFixOut(out)
    }

    return err
}

func (y *YamlParser) yamlFixOut(out interface{}) {
    switch v := out.(type) {
    case *map[string]interface{}:
        for key, val := range *v {
            (*v)[key] = y.yamlCleanValue(val)
        }

    case map[string]interface{}:
        for key, val := range v {
            v[key] = y.yamlCleanValue(val)
        }

    case *interface{}:
        if vv, ok := (*v).(map[interface{}]interface{}); ok {
            *v = y.yamlCleanMap(vv)
        }
    }
}

func (y *YamlParser) yamlCleanValue(v interface{}) interface{} {
    switch vv := v.(type) {
    case map[interface{}]interface{}:
        return y.yamlCleanMap(vv)

    case []interface{}:
        return y.yamlCleanArray(vv)

    default:
        return v
    }
}

func (y *YamlParser) yamlCleanMap(in map[interface{}]interface{}) map[string]interface{} {
    result := make(map[string]interface{}, len(in))
    for k, v := range in {
        result[fmt.Sprintf("%v", k)] = y.yamlCleanValue(v)
    }
    return result
}

func (y *YamlParser) yamlCleanArray(in []interface{}) []interface{} {
    result := make([]interface{}, len(in))
    for k, v := range in {
        result[k] = y.yamlCleanValue(v)
    }
    return result
}
