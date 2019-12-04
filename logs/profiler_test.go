package logs

import (
    "testing"
    "time"
)

func TestProfiler_PushLog(t *testing.T) {
    p := NewProfiler()
    p.PushLog("k", "v")
    if p.PushLogString() != "k=v" {
        t.FailNow()
    }
}

func TestProfiler_Counting(t *testing.T) {
    p := NewProfiler()
    key := "testCounting"
    p.Counting(key, 1, 1)
    p.Counting(key, 0, 1)

    if p.CountingString() != "testCounting=1/2" {
        t.FailNow()
    }
}

func TestProfiler_ProfileStart(t *testing.T) {
    p := NewProfiler()
    p.ProfileStart("test")
    if _, has := p.profileStack["test"]; has == false {
        t.FailNow()
    }
}

func TestProfiler_ProfileStop(t *testing.T) {
    p := NewProfiler()
    p.ProfileStart("test")

    p.ProfileStop("test")

    if _, has := p.profile["test"]; has == false {
        t.FailNow()
    }
}

func TestProfiler_ProfileAdd(t *testing.T) {
    p := NewProfiler()
    elapse := 10 * time.Second
    p.ProfileAdd("test", elapse)
    if _, has := p.profile["test"]; has == false {
        t.FailNow()
    }

    if p.ProfileString() == "" {
        t.FailNow()
    }
}


func TestProfiler_Reset(t *testing.T) {
    p := NewProfiler()
    p.ProfileStart("test")
    p.Reset()
    if p.profileStack != nil {
        t.FailNow()
    }
}
