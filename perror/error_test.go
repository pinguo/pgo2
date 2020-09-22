package perror

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	var obj interface{}

	obj = New(400)
	eObj, ok := obj.(*Error)
	if ok == false {
		t.FailNow()
	}

	if eObj.IsErrTypeErr() == false {
		t.FailNow()
	}
}

func TestNewWarn(t *testing.T) {
	var obj interface{}

	obj = NewWarn(400)
	eObj, ok := obj.(*Error)
	if ok == false {
		t.FailNow()
	}

	if eObj.IsErrTypeWarn() == false {
		t.FailNow()
	}
}

func TestNewIgnore(t *testing.T) {
	var obj interface{}

	obj = NewIgnore(400)
	eObj, ok := obj.(*Error)
	if ok == false {
		t.FailNow()
	}

	if eObj.IsErrTypeIgnore() == false {
		t.FailNow()
	}
}

func TestError_Status(t *testing.T) {
	e := New(400)
	if e.Status() != 400 {
		t.FailNow()
	}
}

func TestError_Message(t *testing.T) {
	e := New(400, "err")
	if e.Message() != "err" {
		t.FailNow()
	}
}

func TestError_Error(t *testing.T) {
	e := New(400, "err%s", "err")
	assertMsg := fmt.Sprintf("errCode: %d, errMsg: %s", 400, fmt.Sprintf("err%s", "err"))
	if e.Error() != assertMsg {
		t.FailNow()
	}
}
