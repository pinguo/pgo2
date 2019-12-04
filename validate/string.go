package validate

import (
    "encoding/json"
    "net/http"
    "regexp"
    "strings"
    "unicode"
    "unicode/utf8"

    "github.com/pinguo/pgo2/perror"
    "github.com/pinguo/pgo2/util"
)

var (
    emailRe  = regexp.MustCompile(`(?i)^[a-z0-9_-]+@[a-z0-9_-]+(\.[a-z0-9_-]+)+$`)
    mobileRe = regexp.MustCompile(`^(\+\d{2,3} )?1[35789]\d{9}$`)
    ipv4Re   = regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$`)
)

// String validator for string value
type String struct {
    Name   string
    UseDft bool
    Value  string
}

func (s *String) Min(v int) *String {
    if !s.UseDft && utf8.RuneCountInString(s.Value) < v {
        panic(perror.New(http.StatusBadRequest, "%s is too short", s.Name))
    }
    return s
}

func (s *String) Max(v int) *String {
    if !s.UseDft && utf8.RuneCountInString(s.Value) > v {
        panic(perror.New(http.StatusBadRequest, "%s is too long", s.Name))
    }
    return s
}

func (s *String) Len(v int) *String {
    if !s.UseDft && utf8.RuneCountInString(s.Value) != v {
        panic(perror.New(http.StatusBadRequest, "%s has invalid length", s.Name))
    }
    return s
}

func (s *String) Enum(enums ...string) *String {
    found := false
    for _, v := range enums {
        if v == s.Value {
            found = true
            break
        }
    }

    if !s.UseDft && !found {
        panic(perror.New(http.StatusBadRequest, "%s is invalid", s.Name))
    }
    return s
}

func (s *String) RegExp(v interface{}) *String {
    var re *regexp.Regexp
    if pat, ok := v.(string); ok {
        re = regexp.MustCompile(pat)
    } else {
        re = v.(*regexp.Regexp)
    }

    if !s.UseDft && !re.MatchString(s.Value) {
        panic(perror.New(http.StatusBadRequest, "%s is invalid", s.Name))
    }

    return s
}

func (s *String) Filter(f func(v, n string) string) *String {
    defer func() {
        if v := recover(); !s.UseDft && v != nil {
            panic(perror.New(http.StatusBadRequest, "%s is invalid", s.Name))
        }
    }()

    if v := f(s.Value, s.Name); len(v) > 0 {
        s.Value = v
    } else if !s.UseDft {
        panic(perror.New(http.StatusBadRequest, "%s is invalid", s.Name))
    }

    return s
}

func (s *String) Password() *String {
    length, number, letter, special := false, false, false, false

    if l := len(s.Value); 6 <= l && l <= 32 {
        length = true
        for i := 0; i < l; i++ {
            switch {
            case unicode.IsNumber(rune(s.Value[i])):
                number = true
            case unicode.IsLetter(rune(s.Value[i])):
                letter = true
            case unicode.IsPunct(rune(s.Value[i])):
                special = true
            case unicode.IsSymbol(rune(s.Value[i])):
                special = true
            }
        }
    }

    if !s.UseDft && (!length || !number || !letter || !special) {
        panic(perror.New(http.StatusBadRequest, "%s is invalid password", s.Name))
    }

    return s
}

func (s *String) Email() *String {
    if !s.UseDft && !emailRe.MatchString(s.Value) {
        panic(perror.New(http.StatusBadRequest, "%s is invalid email", s.Name))
    }

    return s
}

func (s *String) Mobile() *String {
    if !s.UseDft && !mobileRe.MatchString(s.Value) {
        panic(perror.New(http.StatusBadRequest, "%s is invalid mobile", s.Name))
    }

    return s
}

func (s *String) IPv4() *String {
    if !s.UseDft && !ipv4Re.MatchString(s.Value) {
        panic(perror.New(http.StatusBadRequest, "%s is invalid ipv4", s.Name))
    }

    return s
}

func (s *String) Bool() *Bool {
    return &Bool{s.Name, s.UseDft, util.ToBool(s.Value)}
}

func (s *String) Int() *Int {
    return &Int{s.Name, s.UseDft, util.ToInt(s.Value)}
}

func (s *String) Float() *Float {
    return &Float{s.Name, s.UseDft, util.ToFloat(s.Value)}
}

func (s *String) Slice(sep string) *StringSlice {
    validator := &StringSlice{s.Name, s.UseDft, make([]string, 0)}

    if len(s.Value) > 0 {
        parts := strings.Split(s.Value, sep)
        for _, v := range parts {
            validator.Value = append(validator.Value, strings.TrimSpace(v))
        }
    }

    return validator
}

func (s *String) Json() *Json {
    validator := &Json{s.Name, s.UseDft, make(map[string]interface{})}
    decoder := json.NewDecoder(strings.NewReader(s.Value))
    if err := decoder.Decode(&validator.Value); !s.UseDft && err != nil {
        panic(perror.New(http.StatusBadRequest, "%s is invalid json", s.Name))
    }

    return validator
}

func (s *String) Do() string {
    return s.Value
}

// StringSlice validator for string slice value
type StringSlice struct {
    Name   string
    UseDft bool
    Value  []string
}

func (s *StringSlice) Min(v int) *StringSlice {
    if !s.UseDft && len(s.Value) < v {
        panic(perror.New(http.StatusBadRequest, "%s has too few elements", s.Name))
    }
    return s
}

func (s *StringSlice) Max(v int) *StringSlice {
    if !s.UseDft && len(s.Value) > v {
        panic(perror.New(http.StatusBadRequest, "%s has too many elements", s.Name))
    }
    return s
}

func (s *StringSlice) Len(v int) *StringSlice {
    if !s.UseDft && len(s.Value) != v {
        panic(perror.New(http.StatusBadRequest, "%s has invalid length", s.Name))
    }
    return s
}

func (s *StringSlice) Int() *IntSlice {
    validator := &IntSlice{s.Name, make([]int, 0)}
    for _, v := range s.Value {
        validator.Value = append(validator.Value, util.ToInt(v))
    }

    return validator
}

func (s *StringSlice) Float() *FloatSlice {
    validator := &FloatSlice{s.Name, make([]float64, 0)}
    for _, v := range s.Value {
        validator.Value = append(validator.Value, util.ToFloat(v))
    }

    return validator
}

func (s *StringSlice) Do() []string {
    return s.Value
}
