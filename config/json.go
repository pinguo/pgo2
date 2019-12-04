package config

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

// JsonParser parser for json config
type JsonParser struct {
    base
}

// Parse parse json config, environment value like ${env||default} will expand
func (j *JsonParser) Parse(path string) (parseData map[string]interface{}, err error) {
    h, e := os.Open(path)
    if e != nil {
        return nil, fmt.Errorf("JsonParser: failed to open file: " + path)
    }

    defer h.Close()

    content, e := ioutil.ReadAll(h)
    if e != nil {
        return nil, fmt.Errorf("JsonParser: failed to read file: " + path)
    }

    // expand env: ${env||default}
    content = j.expandEnv(content)

    if e := json.Unmarshal(content, &parseData); e != nil {
        return nil, fmt.Errorf("jsonParser: failed to parse file: %s, %s", path, e.Error())
    }

    return parseData, nil
}
