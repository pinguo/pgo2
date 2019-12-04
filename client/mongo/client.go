package mongo

import (
    "fmt"
    "net/url"
    "strings"
    "time"

    "github.com/globalsign/mgo"
    "github.com/pinguo/pgo2/core"
)


// Mongo Client component, configuration:
// components:
//      mongo:
//          dsn: "mongodb://host1:port1/[db][?options]"
//          connectTimeout: "1s"
//          readTimeout: "10s"
//          writeTimeout: "10s"
//
// see Dial() for query options, default:
// replicaSet=
// connect=replicaSet
// maxPoolSize=100
// minPoolSize=1
// maxIdleTimeMS=300000
// ssl=false
// w=1
// j=false
// wtimeoutMS=10000
// readPreference=secondaryPreferred
func New(config map[string]interface{}) (interface{}, error){
    c := &Client{}
    c.dsn = defaultDsn
    c.connectTimeout = defaultConnectTimeout
    c.readTimeout = defaultReadTimeout
    c.writeTimeout = defaultWriteTimeout

    if err:=core.ClientConfigure(c, config); err != nil {
        return nil, err
    }

    if err :=c.Init(); err != nil {
        return nil, err
    }

    return c, nil
}

type Client struct {
    session        *mgo.Session
    dsn            string
    connectTimeout time.Duration
    readTimeout    time.Duration
    writeTimeout   time.Duration
}

func (c *Client) Init() error{
    server, query := c.dsn, defaultOptions
    if pos := strings.IndexByte(c.dsn, '?'); pos > 0 {
        dsnOpts, _ := url.ParseQuery(c.dsn[pos+1:])
        options, _ := url.ParseQuery(defaultOptions)

        for k, v := range dsnOpts {
            if len(v) > 0 && len(v[0]) > 0 {
                options.Set(k, v[0])
            }
        }
        server = c.dsn[:pos]
        query = options.Encode()
    }

    c.dsn = server + "?" + query
    dialInfo, e := mgo.ParseURL(c.dsn)
    if e != nil {
        return fmt.Errorf(errInvalidDsn, c.dsn, e.Error())
    }

    dialInfo.Timeout = c.connectTimeout
    dialInfo.ReadTimeout = c.readTimeout
    dialInfo.WriteTimeout = c.writeTimeout

    if c.session, e = mgo.DialWithInfo(dialInfo); e != nil {
        return fmt.Errorf(errDialFailed, c.dsn, e.Error())
    }

    c.session.SetMode(mgo.Monotonic, true)

    return nil
}

func (c *Client) SetDsn(dsn string) {
    c.dsn = dsn
}

func (c *Client) SetConnectTimeout(v string) error{
    if connectTimeout, e := time.ParseDuration(v); e != nil {
        return fmt.Errorf(errSetProp, "connectTimeout", e.Error())
    } else {
        c.connectTimeout = connectTimeout
    }

    return nil
}

func (c *Client) SetReadTimeout(v string) error{
    if readTimeout, e := time.ParseDuration(v); e != nil {
        return fmt.Errorf(errSetProp, "readTimeout", e.Error())
    } else {
        c.readTimeout = readTimeout
    }

    return nil
}

func (c *Client) SetWriteTimeout(v string) error{
    if writeTimeout, e := time.ParseDuration(v); e != nil {
        return fmt.Errorf(errSetProp, "writeTimeout", e.Error())
    } else {
        c.writeTimeout = writeTimeout
    }

    return nil
}

func (c *Client) GetSession() *mgo.Session {
    return c.session.Copy()
}