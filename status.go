package pgo2

import (
    "fmt"
    "net/http"

    "github.com/pinguo/pgo2/core"
    "github.com/pinguo/pgo2/util"
)

// Status the status component, configuration:
// status:
//     useI18n: false
//     mapping:
//         11002: "Verify Sign Error"
func NewStatus(config map[string]interface{}) *Status {
    s := &Status{
        useI18n: false,
        mapping: make(map[int]string),
    }

    core.Configure(s, config)

    return s
}

type Status struct {
    useI18n bool
    mapping map[int]string
}

// SetUseI18n set whether to use i18n translation
func (s *Status) SetUseI18n(useI18n bool) {
    s.useI18n = useI18n
}

// SetMapping set mapping from status code to text
func (s *Status) SetMapping(m map[string]interface{}) {
    for k, v := range m {
        s.mapping[util.ToInt(k)] = util.ToString(v)
    }
}

// GetText get status text
func (s *Status) Text(status int, lang string, dft ...string) string {
    txt, ok := s.mapping[status]
    if !ok {
        if len(dft) == 0 || len(dft[0]) == 0 {
            if txt = http.StatusText(status); len(txt) == 0 {
                panic(fmt.Sprintf("unknown status code: %d", status))
            }
        } else {
            txt = dft[0]
        }
    }

    if s.useI18n && lang != "" {
        txt = App().I18n().Translate(txt, lang)
    }

    return txt
}
