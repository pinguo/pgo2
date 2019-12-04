package autoid

import (
    "testing"
)

func TestAuto_setParam(t *testing.T) {
    a := &AutoId{
        start:      0,
        step:       1,
        queueNum:   100,
        remakesNum: 999999999,
        onOff:      true,
    }

    config := make(map[string]interface{})

    config["start"] = int64(2)
    config["step"] = int64(2)
    config["queueNum"] = int64(10)
    config["remakesNum"] = int64(100)

    a.setParam(config)

    if a.start != 2 {
        t.Fatal(`a.start != 2`)
    }

    if a.step != 2 {
        t.Fatal(`a.step != 2`)
    }

    if a.queueNum != 10 {
        t.Fatal(`a.queueNum != 10`)
    }

    if a.remakesNum != 100 {
        t.Fatal(`a.remakesNum != 100`)
    }

}

func TestAuto(t *testing.T) {
    config := make(map[string]interface{})

    config["start"] = int64(0)
    config["step"] = int64(1)
    config["queueNum"] = int64(5)
    config["remakesNum"] = int64(4)
    a := New(config).(*AutoId)
    for i := 0; i <= 10; i++ {
        id := a.Id()

        if i == 4 && id != 4 {
            t.Fatal(`i==4 && id !=4`)
        }

        if i == 5 && id != 0 {
            t.Fatal(`i == 5 && id!=0`)
        }

        if i >= 10 {
            a.Close()
            if a.onOff != false {
                t.Fatal(` a.onOff !=false`)
            }
        }

    }
}
