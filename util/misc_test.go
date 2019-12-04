package util

import (
    "bytes"
    "os"
    "testing"
)

func TestGenUniqueId(t *testing.T) {
    if len(GenUniqueId()) != 24 {
        t.FailNow()
    }
}

func TestExpandEnv(t *testing.T) {
    t.Run("haveEnv", func(t *testing.T) {
        data := "${MockUnitTest||dftDd}"
        os.Setenv("MockUnitTest", "dd")
        ret := ExpandEnv([]byte(data))

        if bytes.Equal(ret, []byte("dd")) == false {
            t.FailNow()
        }

    })

    t.Run("default", func(t *testing.T) {
        data := "${MockUnitTest1||dftDd}"
        ret := ExpandEnv([]byte(data))

        if bytes.Equal(ret, []byte("dftDd")) == false {
            t.FailNow()
        }
    })

    t.Run("origin", func(t *testing.T) {
        data := "${MockUnitTest1}"
        ret := ExpandEnv([]byte(data))
        if bytes.Equal(ret, []byte(data)) == false {
            t.FailNow()
        }
    })
}

func TestFormatLanguage(t *testing.T) {
    t.Run("empty", func(t *testing.T) {
        if FormatLanguage("23") != "" {
            t.FailNow()
        }
    })

    t.Run("matches-1", func(t *testing.T) {
        if FormatLanguage("zh") != "zh" {
            t.FailNow()
        }
    })

    t.Run("map-1", func(t *testing.T) {
        if FormatLanguage("zh-CHS") != "zh-CN" {
            t.FailNow()
        }
    })

    t.Run("map-2", func(t *testing.T) {
        if FormatLanguage("zh-CHT") != "zh-TW" {
            t.FailNow()
        }
    })

}

func TestFormatVersion(t *testing.T) {
    if FormatVersion("v10...2....2.1-alpha", 5) != "v10.2.2.1.0-alpha" {
        t.FailNow()
    }
}

func TestVersionCompare(t *testing.T) {
    t.Run("<", func(t *testing.T) {
        if VersionCompare("12.1.1-dev", "12.1.1-alpha") != -1 {
            t.FailNow()
        }
    })

    t.Run("=", func(t *testing.T) {
        if VersionCompare("12.1.1", "12.1.1") != 0 {
            t.FailNow()
        }
    })

    t.Run(">", func(t *testing.T) {
        if VersionCompare("12.1.10", "12.1.1") != 1 {
            t.FailNow()
        }
    })
}
