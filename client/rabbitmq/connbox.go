package rabbitmq

import (
    "crypto/tls"
    "crypto/x509"
    "errors"
    "io/ioutil"
    "sync"
    "time"

    "github.com/pinguo/pgo2"
    "github.com/pinguo/pgo2/util"
    "github.com/streadway/amqp"
)

func newConnBox(id, addr, dsn string, maxChannelNum int, tlsDft ...string) (*ConnBox, error) {
    tlsCert, tlsCertKey, tlsRootCAs := "", "", ""
    if len(tlsDft) > 0 {
        tlsCert = tlsDft[0]
        tlsCertKey = tlsDft[1]
        tlsRootCAs = tlsDft[2]
    }

    connBox := &ConnBox{id: id, addr: addr, dsn: dsn, tlsCert: tlsCert, tlsCertKey: tlsCertKey, tlsRootCAs: tlsRootCAs, maxChannelNum: maxChannelNum}
    err := connBox.initConn()

    return connBox,err
}

type ConnBox struct {
    id              string
    addr            string
    useConnCount    int
    useChannelCount int
    channelList     chan *ChannelBox
    lock            sync.RWMutex
    newChannelLock  sync.RWMutex
    startTime       time.Time

    connection *amqp.Connection
    tlsCert    string
    tlsCertKey string
    tlsRootCAs string
    dsn        string

    maxChannelNum int

    notifyClose chan *amqp.Error

    disable bool
}

func (c *ConnBox) setEnable() {
    c.lock.Lock()
    defer c.lock.Unlock()
    c.disable = false
}

func (c *ConnBox) setDisable() {
    c.lock.Lock()
    defer c.lock.Unlock()
    c.disable = true
    c.close()
}

func (c *ConnBox) initConn() (retErr error) {

    func() {
        c.lock.Lock()
        defer c.lock.Unlock()
        var err error

        if c.tlsCert != "" && c.tlsCertKey != "" {
            cfg := new(tls.Config)
            if c.tlsRootCAs != "" {
                cfg.RootCAs = x509.NewCertPool()
                if ca, err := ioutil.ReadFile(c.tlsRootCAs); err == nil {
                    cfg.RootCAs.AppendCertsFromPEM(ca)
                }
            }

            if cert, err := tls.LoadX509KeyPair(c.tlsCert, c.tlsCertKey); err == nil {
                cfg.Certificates = append(cfg.Certificates, cert)
            }

            c.connection, err = amqp.DialTLS(c.dsn, cfg)
        } else {
            c.connection, err = amqp.Dial(c.dsn)
        }

        if err != nil {
            errMsg:= err.Error()
            retErr = errors.New("Failed to connect to RabbitMQ:" + errMsg)
        }

        c.disable = false
        c.channelList = make(chan *ChannelBox, c.maxChannelNum)
        c.notifyClose = make(chan *amqp.Error)
        c.startTime = time.Now()
        c.connection.NotifyClose(c.notifyClose)
    }()

    go c.check(c.startTime)

    return retErr
}

func (c *ConnBox) check(startTime time.Time) {
    defer func() {
        if err := recover(); err != nil {
            pgo2.GLogger().Error("Rabbit ConnBox.check err:" + util.ToString(err))
        }
    }()

    for {
        if c.startTime != startTime {
            // 自毁
            return
        }

        select {
        case err, ok := <-c.notifyClose:
            if ok == false {
                return
            }

            if err != nil {
                func() {
                    defer func() {
                        if err := recover(); err != nil {
                            pgo2.GLogger().Error("Rabbit ConnBox.check start initConn err:" + util.ToString(err))
                        }
                    }()

                    c.setDisable()

                    if err := c.initConn();err!=nil{
                        pgo2.GLogger().Error("Rabbit ConnBox.check start initConn err:" + util.ToString(err))
                    }
                }()
                return
            }

        default:
            time.Sleep(100 * time.Microsecond)

        }
    }
}

func (c *ConnBox) isClosed() bool {
    if c.disable || c.connection.IsClosed() {
        return true
    }
    return false
}

func (c *ConnBox) close() {
    if c.connection != nil && c.connection.IsClosed() == false {
        err := c.connection.Close()
        if err != nil {
            pgo2.GLogger().Warn("Rabbit ConnBox.close err:" + err.Error())
        }
    }
}
