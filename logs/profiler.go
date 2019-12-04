package logs

import (
    "fmt"
    "strings"
    "time"

    "github.com/pinguo/pgo2/util"
)

func NewProfiler() *Profiler {
    return &Profiler{}
}

// Profiler
type Profiler struct {
    pushLog      []string
    counting     map[string][2]int
    profile      map[string][2]int
    profileStack map[string]time.Time
}

func (p *Profiler) Reset() {
    p.pushLog = nil
    p.counting = nil
    p.profile = nil
    p.profileStack = nil
}

// PushLog add push log, the push log string is key=Util.ToString(v)
func (p *Profiler) PushLog(key string, v interface{}) {
    if p.pushLog == nil {
        p.pushLog = make([]string, 0)
    }

    pl := key + "=" + util.ToString(v)
    p.pushLog = append(p.pushLog, pl)
}

// Counting add counting info, the counting string is key=sum(hit)/sum(total)
func (p *Profiler) Counting(key string, hit, total int) {
    if p.counting == nil {
        p.counting = make(map[string][2]int)
    }

    v := p.counting[key]

    if hit > 0 {
        v[0] += hit
    }

    if total <= 0 {
        total = 1
    }

    v[1] += total
    p.counting[key] = v
}

// ProfileStart mark start of profile
func (p *Profiler) ProfileStart(key string) {
    if p.profileStack == nil {
        p.profileStack = make(map[string]time.Time)
    }

    p.profileStack[key] = time.Now()
}

// ProfileStop mark stop of profile
func (p *Profiler) ProfileStop(key string) {
    if startTime, ok := p.profileStack[key]; ok {
        delete(p.profileStack, key)
        p.ProfileAdd(key, time.Now().Sub(startTime))
    }
}

// ProfileAdd add profile info, the profile string is key=sum(elapse)/count
func (p *Profiler) ProfileAdd(key string, elapse time.Duration) {
    if p.profile == nil {
        p.profile = make(map[string][2]int)
    }

    v, _ := p.profile[key]
    v[0] += int(elapse.Nanoseconds() / 1e6)
    v[1] += 1

    p.profile[key] = v
}

// GetPushLogString get push log string
func (p *Profiler) PushLogString() string {
    if len(p.pushLog) == 0 {
        return ""
    }

    return strings.Join(p.pushLog, " ")
}

// GetCountingString get counting info string
func (p *Profiler) CountingString() string {
    if len(p.counting) == 0 {
        return ""
    }

    cs := make([]string, 0)
    for k, v := range p.counting {
        cs = append(cs, fmt.Sprintf("%s=%d/%d", k, v[0], v[1]))
    }

    return strings.Join(cs, " ")
}

// GetProfileString get profile info string
func (p *Profiler) ProfileString() string {
    if len(p.profile) == 0 {
        return ""
    }

    ps := make([]string, 0)
    for k, v := range p.profile {
        ps = append(ps, fmt.Sprintf("%s=%dms/%d", k, v[0], v[1]))
    }

    return strings.Join(ps, " ")
}
