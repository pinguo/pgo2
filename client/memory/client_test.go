package memory

import "testing"

func TestNew(t *testing.T) {
    var obj interface{}
    obj, err := New(nil)
    if _, ok := obj.(*Client); ok == false {
        t.FailNow()
    }

    if err != nil {
        t.FailNow()
    }
}

func TestClient_Add(t *testing.T) {

}

func TestClient_Del(t *testing.T) {

}

func TestClient_Exists(t *testing.T) {

}

func TestClient_Get(t *testing.T) {

}

func TestClient_Incr(t *testing.T) {

}

func TestClient_MAdd(t *testing.T) {

}

func TestClient_MDel(t *testing.T) {

}

func TestClient_MGet(t *testing.T) {

}

func TestClient_MSet(t *testing.T) {

}

func TestClient_Set(t *testing.T) {

}

func TestClient_SetGcInterval(t *testing.T) {

}

func TestClient_SetGcMaxItems(t *testing.T) {

}
