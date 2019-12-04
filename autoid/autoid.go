package autoid

import "fmt"

func New(config map[string]interface{}) interface{} {

    a := &AutoId{
        start:      0,
        step:       1,
        queueNum:   100,
        remakesNum: 999999999,
        onOff:      true,
    }
    a.setParam(config)
    a.queue = make(chan int64, a.queueNum)
    go a.makeId()

    return a
}

type AutoId struct {
    start      int64
    step       int64
    queueNum   int64
    queue      chan int64
    remakesNum int64
    onOff      bool
}

func (a *AutoId) setParam(config map[string]interface{}) {
    if config == nil {
        return
    }

    if cStart, ok := config["start"]; ok {
        a.start = cStart.(int64)
    }

    if cStep, ok := config["step"]; ok {
        a.step = cStep.(int64)
    }

    if cQueueNum, ok := config["queueNum"]; ok {
        a.queueNum = cQueueNum.(int64)
    }

    if cRemakesNum, ok := config["remakesNum"]; ok {
        a.remakesNum = cRemakesNum.(int64)
    }
}

func (a *AutoId) makeId() {
    defer func() {
        if err := recover(); err != nil {
            fmt.Println("makeId err", err)
        }
    }()

    for i := a.start; a.onOff; i += a.step {
        if i > a.remakesNum {
            i = a.start
        }
        a.queue <- i
    }
}

func (a *AutoId) Id() int64 {
    return <-a.queue
}

func (a *AutoId) Close() {
    a.onOff = false
    close(a.queue)
}
