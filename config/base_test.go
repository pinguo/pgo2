package config

import (
    "bytes"
    "os"
    "testing"
)

func TestBase_expandEnv(t *testing.T) {

    t.Run("haveEnv", func(t *testing.T) {
        data := "${MockUnitTest||dftDd}"
        os.Setenv("MockUnitTest", "dd")
        base := &base{}
        ret := base.expandEnv([]byte(data))

        if bytes.Equal(ret, []byte("dd")) == false {
            t.FailNow()
        }

    })

    t.Run("default", func(t *testing.T) {
        data := "${MockUnitTest1||dftDd}"
        base := &base{}
        ret := base.expandEnv([]byte(data))

        if bytes.Equal(ret, []byte("dftDd")) == false {
            t.FailNow()
        }
    })

    t.Run("origin", func(t *testing.T) {
        data := "${MockUnitTest1}"
        base := &base{}
        ret := base.expandEnv([]byte(data))
        if bytes.Equal(ret, []byte(data)) == false {
            t.FailNow()
        }
    })
}
