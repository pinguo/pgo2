package memcache

import (
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "io"
    "net"
    "strconv"
    "strings"
    "time"
)

func newConn(addr string, nc net.Conn, pool *Pool) *Conn {
    r := bufio.NewReader(nc)
    w := bufio.NewWriter(nc)

    return &Conn{
        addr: addr,
        nc:   nc,
        rw:   bufio.NewReadWriter(r, w),
        pool: pool,
    }
}

type Item struct {
    Key   string // the key under which to store data
    Data  []byte // the bytes stored under the key
    Flags uint16 // a uint16 associated with the data
    CasId uint64 // a unique uint64 of an existing item
}

type Conn struct {
    prev       *Conn
    next       *Conn
    lastActive time.Time

    addr string
    nc   net.Conn
    rw   *bufio.ReadWriter
    pool *Pool
    down bool
}

func (c *Conn) Close(force bool) {
    if force || c.down || !c.pool.putFreeConn(c) {
        c.nc.Close()
    } else {
        c.lastActive = time.Now()
    }
}

func (c *Conn) CheckActive() (bool, error) {
    if time.Since(c.lastActive) < c.pool.maxIdleTime {
        return true, nil
    }

    c.ExtendDeadLine()
    v, err := c.Version()
    if err != nil {
        return false, err
    }
    return len(v) != 0, nil
}

func (c *Conn) ExtendDeadLine(deadLine ...time.Duration) bool {
    deadLine = append(deadLine, c.pool.netTimeout)
    return c.nc.SetDeadline(time.Now().Add(deadLine[0])) == nil
}

// execute store command, cmd is set, add, replace, append, prepend or cas,
// expire is expiration time, either unix timestamp or offset in seconds from now,
// 0 means never expires, negative value means immediately expired
// return true if succeed, otherwise false
func (c *Conn) Store(cmd string, item *Item, expire int) (bool, error) {
    switch cmd {
    case CmdCas:
        fmt.Fprintf(c.rw, "%s %s %d %d %d %d\r\n", cmd, item.Key, item.Flags, expire, len(item.Data), item.CasId)
    case CmdSet, CmdAdd, CmdReplace, CmdAppend, CmdPrepend:
        fmt.Fprintf(c.rw, "%s %s %d %d %d\r\n", cmd, item.Key, item.Flags, expire, len(item.Data))
    default:
        return false, errors.New(errInvalidCmd + cmd)
    }

    c.rw.Write(item.Data)
    c.rw.Write(lineEnding)
    if e := c.rw.Flush(); e != nil {
        return false, c.parseError(errSendFailed+e.Error(), true)
    }

    if line, e := c.rw.ReadSlice('\n'); e != nil {
        return false, c.parseError(errReadFailed+e.Error(), true)
    } else if bytes.Equal(line, replyStored) {
        return true, nil
    } else if bytes.Equal(line, replyNotStored) {
        return false, nil
    } else if bytes.Equal(line, replyExists) {
        return false, nil
    } else if bytes.Equal(line, replyNotFound) {
        return false, nil
    } else {
        return false, c.parseError(errBase+string(line), false)
    }
    return false, nil
}

// execute retrieve command, cmd is get or gets,
// item retrieved by gets cmd has a unique uint64 value,
// result expects zero or more items, if some of the keys
// not exists or expired, then corresponding item do not
// present in the result
func (c *Conn) Retrieve(cmd string, keys ...string) (items []*Item, err error) {
    if cmd != CmdGet && cmd != CmdGets {
        return nil, errors.New(errInvalidCmd + cmd)
    }

    if len(keys) < 1 {
        return nil, errors.New(errEmptyKeys)
    }

    c.rw.WriteString(cmd)
    for _, v := range keys {
        c.rw.WriteByte(' ')
        c.rw.WriteString(v)
    }

    c.rw.Write(lineEnding)
    if e := c.rw.Flush(); e != nil {
        return nil, c.parseError(errSendFailed+e.Error(), true)
    }

    for {
        line, e := c.rw.ReadSlice('\n')
        if e != nil {
            return nil, c.parseError(errReadFailed+e.Error(), true)
        }

        if bytes.Equal(line, replyEnd) {
            break
        }

        rd, item, size := bytes.NewReader(line), new(Item), 0
        if cmd == CmdGet {
            _, e = fmt.Fscanf(rd, "VALUE %s %d %d\r\n", &item.Key, &item.Flags, &size)
        } else {
            _, e = fmt.Fscanf(rd, "VALUE %s %d %d %d\r\n", &item.Key, &item.Flags, &size, &item.CasId)
        }

        if e != nil {
            return nil, c.parseError(errBase+string(line), false)
        } else {
            item.Data = make([]byte, size+2)
            if _, e = io.ReadFull(c.rw, item.Data); e != nil {
                return nil, c.parseError(errCorrupted+e.Error(), true)
            }

            item.Data = item.Data[:size]
            items = append(items, item)
        }
    }

    return items, nil
}

// execute delete command,
// return true if succeed otherwise false
func (c *Conn) Delete(key string) (bool, error) {
    fmt.Fprintf(c.rw, "%s %s\r\n", CmdDelete, key)
    if e := c.rw.Flush(); e != nil {
        return false, c.parseError(errSendFailed+e.Error(), true)
    }

    if line, e := c.rw.ReadSlice('\n'); e != nil {
        return false, c.parseError(errReadFailed+e.Error(), true)
    } else if bytes.Equal(line, replyDeleted) {
        return true, nil
    } else if !bytes.Equal(line, replyNotFound) {
        return false, c.parseError(errBase+string(line), false)
    }
    return false, nil
}

// execute increment/decrement command,
// negative delta for decrement, otherwise for increment,
// if key not found, treat the data as zero,
// if the data is not uint64 representation, function panic.
// if decrease data below 0, new data will be 0.
func (c *Conn) Increment(key string, delta int) (int, error) {
    if delta > 0 {
        fmt.Fprintf(c.rw, "%s %s %d\r\n", CmdIncr, key, delta)
    } else {
        fmt.Fprintf(c.rw, "%s %s %d\r\n", CmdDecr, key, -delta)
    }

    if e := c.rw.Flush(); e != nil {
        return 0, c.parseError(errSendFailed+e.Error(), true)
    }

    if line, e := c.rw.ReadSlice('\n'); e != nil {
        return 0, c.parseError(errReadFailed+e.Error(), true)
    } else if bytes.Equal(line, replyNotFound) {
        if delta < 0 {
            delta = 0
        }

        data := strconv.AppendInt(nil, int64(delta), 10)
        c.Store(CmdSet, &Item{Key: key, Data: data}, 0)
        return delta, nil
    } else {
        rd, result := bytes.NewReader(line), 0
        if _, e := fmt.Fscanf(rd, "%d\r\n", &result); e != nil {
            return 0, c.parseError(errBase+string(line), false)
        }
        return result, nil
    }
}

// execute touch command
// expire is expiration time, either unix timestamp or offset in seconds from now,
// 0 means never expires, negative value means immediately expired
func (c *Conn) Touch(key string, expire int) (bool, error) {
    fmt.Fprintf(c.rw, "%s %s %d\r\n", CmdTouch, key, expire)
    if e := c.rw.Flush(); e != nil {
        return false, c.parseError(errSendFailed+e.Error(), true)
    }

    if line, e := c.rw.ReadSlice('\n'); e != nil {
        c.parseError(errReadFailed+e.Error(), true)
    } else if bytes.Equal(line, replyTouched) {
        return true, nil
    } else if bytes.Equal(line, replyNotFound) {
        return false, nil
    } else {
        return false, c.parseError(errBase+string(line), false)
    }
    return false, nil
}

// execute stats command
// if args is empty, all statistics will be returned,
// otherwise specified field will be returned
func (c *Conn) Stats(args ...string) (map[string]string, error) {
    c.rw.WriteString(CmdStats)
    for _, v := range args {
        c.rw.WriteByte(' ')
        c.rw.WriteString(v)
    }

    c.rw.Write(lineEnding)
    if e := c.rw.Flush(); e != nil {
        return nil, c.parseError(errSendFailed+e.Error(), true)
    }

    stats := make(map[string]string)

    for {
        line, e := c.rw.ReadSlice('\n')
        if e != nil {
            return nil, c.parseError(errReadFailed+e.Error(), true)
        }

        if bytes.Equal(line, replyEnd) {
            break
        }

        rd, key, value := bytes.NewReader(line), "", ""
        if _, e := fmt.Fscanf(rd, "STAT %s %s\r\n", &key, &value); e != nil {
            return nil, c.parseError(errBase+string(line), true)
        }

        stats[key] = value
    }

    return stats, nil
}

// execute flush_all command
func (c *Conn) FlushAll() (bool, error) {
    c.rw.WriteString(CmdFlushAll)
    c.rw.Write(lineEnding)
    if e := c.rw.Flush(); e != nil {
        return false, c.parseError(errSendFailed+e.Error(), true)
    }

    if line, e := c.rw.ReadSlice('\n'); e != nil {
        return false, c.parseError(errReadFailed+e.Error(), true)
    } else if bytes.Equal(line, replyOK) {
        return true, nil
    } else {
        return false, c.parseError(errBase+string(line), false)
    }
    return false, nil
}

// execute version command
func (c *Conn) Version() (string, error) {
    c.rw.WriteString(CmdVersion)
    c.rw.Write(lineEnding)
    if e := c.rw.Flush(); e != nil {
        c.parseError(errSendFailed+e.Error(), true)
    }

    var ver string
    if line, e := c.rw.ReadSlice('\n'); e != nil {
        return "", c.parseError(errReadFailed+e.Error(), true)
    } else {
        rd := bytes.NewReader(line)
        if _, e := fmt.Fscanf(rd, "VERSION %s\r\n", &ver); e != nil {
            return "", c.parseError(errBase+string(line), false)
        }
    }
    return ver, nil
}

func (c *Conn) parseError(err string, fatal bool) error {
    if fatal {
        c.down = true
    }

    return errors.New(strings.TrimLeft(err, "\r\n"))
}
