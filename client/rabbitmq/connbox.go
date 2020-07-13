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

	return connBox, err
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
			errMsg := err.Error()
			retErr = errors.New("Failed to connect to RabbitMQ:" + errMsg)
			return
		}

		c.disable = false
		c.useConnCount = 0
		c.useChannelCount = 0
		c.channelList = make(chan *ChannelBox, c.maxChannelNum)
		c.notifyClose = make(chan *amqp.Error, 1)
		c.startTime = time.Now()
		//c.connection.NotifyClose(c.notifyClose)
	}()

	// go c.check(c.startTime)

	return retErr
}

func (c *ConnBox) check(startTime time.Time) {
	defer func() {
		if err := recover(); err != nil {
			pgo2.GLogger().Error("Rabbit ConnBox.check err:" + util.ToString(err))
		}
	}()

	timeTicker := time.NewTicker(time.Second)

	for {
		select {
		case notifyErr, ok := <-c.notifyClose:
			if ok == false {
				return
			}

			if notifyErr != nil {
				func() {
					defer func() {
						if err := recover(); err != nil {
							pgo2.GLogger().Error("Rabbit ConnBox.check start initConn err1:" + util.ToString(err))
						}
					}()

					pgo2.GLogger().Error("Rabbit ConnBox.check notifyErr != nil  err1:" + util.ToString(notifyErr))
					c.setDisable()

					if err := c.initConn(); err != nil {
						pgo2.GLogger().Error("Rabbit ConnBox.check start initConn err:" + util.ToString(err))
					}
				}()
				return
			}
		case <-timeTicker.C:
			if c.startTime != startTime {
				// 自毁
				timeTicker.Stop()
				return
			}

		}
	}
}

func (c *ConnBox) Disable() bool {
	return c.disable
}

func (c *ConnBox) isClosed() bool {
	if c.disable || c.connection.IsClosed() {
		// pgo2.GLogger().Info("disable",c.disable,"c.connection.IsClosed()",c.connection.IsClosed())
		return true
	}
	return false
}

func (c *ConnBox) close() {
	// pgo2.GLogger().Info("ConnBox.close begin")
	if c.connection != nil && c.connection.IsClosed() == false {
		// pgo2.GLogger().Info("ConnBox.close 1 len(channelList)=",len(c.channelList))
		l := len(c.channelList)
		// pgo2.GLogger().Info("sl:", l)
		for i := 0; i < l; i++ {
			select {
			case v := <-c.channelList:
				// pgo2.GLogger().Info("ssss:", i)
				v.Close(true)
			default:
			}
		}

		// pgo2.GLogger().Info("ConnBox.close len(channelList)=",len(c.channelList))
		err := c.connection.Close()
		if err != nil {
			pgo2.GLogger().Warn("Rabbit ConnBox.close err:" + err.Error())
		}

	}
}
