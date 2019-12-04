package util

import (
    "bytes"
    "testing"
)

func TestIsAllDigit(t *testing.T) {
    t.Run("true", func(t *testing.T) {
        if IsAllDigit("12323") != true {
            t.FailNow()
        }
    })

    t.Run("false", func(t *testing.T) {
        if IsAllDigit("12323a") != false {
            t.FailNow()
        }
    })
}

func TestIsAllLetter(t *testing.T) {
    t.Run("true", func(t *testing.T) {
        if IsAllLetter("aadddfgdd") != true {
            t.FailNow()
        }
    })

    t.Run("false", func(t *testing.T) {
        if IsAllLetter("12323add") != false {
            t.FailNow()
        }
    })
}

func TestIsAllLower(t *testing.T) {
    t.Run("true", func(t *testing.T) {
        if IsAllLower("aadddfgdd") != true {
            t.FailNow()
        }
    })

    t.Run("false", func(t *testing.T) {
        if IsAllLower("12323aAAA") != false {
            t.FailNow()
        }
    })
}

func TestIsAllUpper(t *testing.T) {
    t.Run("true", func(t *testing.T) {
        if IsAllUpper("ADGGC") != true {
            t.FailNow()
        }
    })

    t.Run("false", func(t *testing.T) {
        if IsAllUpper("12323aAAA") != false {
            t.FailNow()
        }
    })
}

func TestMd5String(t *testing.T) {
    t.Run("string", func(t *testing.T) {
        if Md5String("aa") != "4124bc0a9335c27f086f24ba207a4912" {
            t.FailNow()
        }
    })

    t.Run("json", func(t *testing.T) {
        if Md5String([]string{"aa"}) != "ca9691df1c54652c15d56b9cee80fb5d" {
            t.FailNow()
        }
    })

    t.Run("other", func(t *testing.T) {
        aa := func() {}
        if Md5String(aa) != "0aa68bb7d22c92f23dcce05ae16dc6b0" {
            t.FailNow()
        }
    })
}

func TestMd5Bytes(t *testing.T) {
    var bRet []byte
    bRet = []byte{65, 36, 188, 10, 147, 53, 194, 127, 8, 111, 36, 186, 32, 122, 73, 18}
    if bytes.Equal(Md5Bytes("aa"), bRet) == false {
        t.FailNow()
    }
}
