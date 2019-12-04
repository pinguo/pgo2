package redis

import (
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "io"
    "net"
    "strconv"
    "time"

    "github.com/pinguo/pgo2/value"
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
    ret,err := c.Do("PING")
    if err != nil {
        return false, err
    }
    payload, ok := ret.([]byte)
    return ok && bytes.Equal(payload, replyPong), nil
}

func (c *Conn) ExtendDeadLine(deadLine ...time.Duration) bool {
    deadLine = append(deadLine, c.pool.netTimeout)
    return c.nc.SetDeadline(time.Now().Add(deadLine[0])) == nil
}

func (c *Conn) Do(cmd string, args ...interface{}) (interface{},error) {
    err := c.WriteCmd(cmd, args...)
    if err != nil{
        return nil, err
    }

    return c.ReadReply()
}

func (c *Conn) WriteCmd(cmd string, args ...interface{}) error{
    fmt.Fprintf(c.rw, "*%d\r\n$%d\r\n%s\r\n", len(args)+1, len(cmd), cmd)
    for _, arg := range args {
        argBytes := value.Encode(arg)
        fmt.Fprintf(c.rw, "$%d\r\n", len(argBytes))
        c.rw.Write(argBytes)
        c.rw.Write(lineEnding)
    }

    if e := c.rw.Flush(); e != nil {
        return c.parseError(errSendFailed+e.Error(), true)
    }

    return nil
}

// read reply from server,
// return []byte, int, nil or slice of these types
func (c *Conn) ReadReply() (interface{}, error) {
    line, e := c.rw.ReadSlice('\n')
    if e != nil {
        return nil, c.parseError(errReadFailed+e.Error(), true)
    }

    if !bytes.HasSuffix(line, lineEnding) {
        return nil, c.parseError(errCorrupted+"unexpected line ending", true)
    }

    payload := line[1 : len(line)-2]

    switch line[0] {
    case '+':
        // status response: +<data bytes>\r\n, eg. +OK\r\n, +PONG\r\n
        if bytes.Equal(payload, replyOK) {
            return replyOK, nil
        } else if bytes.Equal(payload, replyPong) {
            return replyPong, nil
        } else {
            data := make([]byte, len(payload))
            copy(data, payload)
            return data, nil
        }

    case '-':
        // error response:-<data bytes>\r\n, eg. -Err unknown command\r\n
        return nil, c.parseError(errBase+string(payload), false)

    case ':':
        // integer response: :<integer>\r\n, eg. :99\r\n
        if n, e := strconv.Atoi(string(payload)); e != nil {
            return nil, c.parseError(errCorrupted+e.Error(), true)
        } else {
            return n, nil
        }

    case '$':
        // bulk string response: $<bytes of data>\r\n<binary data>\r\n,
        // -1 for nil response. eg. $7\r\nfoo bar\r\n
        if size, e := strconv.Atoi(string(payload)); e != nil {
            return nil, c.parseError(errCorrupted+e.Error(), true)
        } else if size >= 0 {
            data := make([]byte, size+2)
            if _, e := io.ReadFull(c.rw, data); e != nil {
                return nil, c.parseError(errCorrupted+e.Error(), true)
            }
            return data[:size], nil
        }

    case '*':
        // multi response: *<argc>\r\n$<bytes of arg1>\r\n<data of arg1>\r\n[argN...],
        // -1 for nil response. eg. *2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n
        if argc, e := strconv.Atoi(string(payload)); e != nil {
            c.parseError(errCorrupted+e.Error(), true)
        } else if argc >= 0 {
            argv := make([]interface{}, argc)
            for i := range argv {
                var err error
                argv[i],err = c.ReadReply()
                if err != nil {
                    return nil,err
                }
            }
            return argv, nil
        }

    default:
        c.parseError(errInvalidResp+string(line[:1]), true)
    }
    return nil,nil
}

func (c *Conn) parseError(err string, fatal bool) error {
    if fatal {
        c.down = true
    }

    return errors.New(err)
}
