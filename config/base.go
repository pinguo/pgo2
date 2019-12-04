package config

import (
    "bytes"
    "os"
    "regexp"
)

// env expand regexp: ${env}, ${env||default}
var envRe = regexp.MustCompile(`\$\{[^\}\|]+(\|\|[^\$\{\}]+?)?\}`)

type base struct {
}

// ExpandEnv expand env variables, format: ${env}, ${env||default}
func (b *base) expandEnv(data []byte) []byte {
    rf := func(s []byte) []byte {
        tmp := bytes.Split(s[2:len(s)-1], []byte{'|', '|'})
        env := bytes.TrimSpace(tmp[0])

        if val, ok := os.LookupEnv(string(env)); ok {
            // return env value
            return []byte(val)
        } else if len(tmp) > 1 {
            // return default value
            return bytes.TrimSpace(tmp[1])
        }

        // return original
        return s
    }

    return envRe.ReplaceAllFunc(data, rf)
}
