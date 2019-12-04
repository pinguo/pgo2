package perror

import (
    "fmt"
    "testing"
)

func TestNew(t *testing.T) {
    var obj interface{}

    obj = New(100)

    if _, ok := obj.(*Error); ok == false {
        t.FailNow()
    }
}

func TestError_Status(t *testing.T) {
    e := New(100)
    if e.Status() != 100 {
        t.FailNow()
    }
}

func TestError_Message(t *testing.T) {
    e := New(100, "err")
    if e.Message() != "err" {
        t.FailNow()
    }
}

func TestError_Error(t *testing.T) {
    e := New(100, "err%s", "err")
    assertMsg := fmt.Sprintf("errCode: %d, errMsg: %s", 100, fmt.Sprintf("err%s", "err"))
    if e.Error() != assertMsg {
        t.FailNow()
    }
}
