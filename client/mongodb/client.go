package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/pinguo/pgo2/core"
	"github.com/qiniu/qmgo"
	qopts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo Client component, configuration:
// components:
//      mongo:
//          dsn: "mongodb://host1:port1/[db][?options]"
//          connectTimeout: "1s"
//          readTimeout: "10s"
//          writeTimeout: "10s"
//
// https://docs.mongodb.com/manual/reference/connection-string/
// query options, default:
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

func New(config map[string]interface{}) (interface{}, error) {
	c := &Client{}
	c.dsn = defaultDsn
	c.connectTimeout = defaultConnectTimeout
	c.readTimeout = defaultReadTimeout
	c.writeTimeout = defaultWriteTimeout

	if err := core.ClientConfigure(c, config); err != nil {
		return nil, err
	}

	var clientOption qopts.ClientOptions
	var okClientOption bool
	if iClientOption, has := config["clientOption"]; has {
		if clientOption, okClientOption = iClientOption.(qopts.ClientOptions); !okClientOption {
			return nil, fmt.Errorf(errInvalidConfig, "clientOption", "need options.ClientOptions")
		}
	}


	if err := c.Init(clientOption); err != nil {
		return nil, err
	}

	return c, nil
}

type Client struct {
	MClient        *qmgo.Client
	dsn            string
	connectTimeout time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
}

func (c *Client) Init(clientOption qopts.ClientOptions) error {
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

	ctx := context.Background()
	if reflect.DeepEqual(clientOption, qopts.ClientOptions{}){
		clientOption = qopts.ClientOptions{ClientOptions:&options.ClientOptions{}}
	}

	if clientOption.ClientOptions.ConnectTimeout == nil {
		clientOption.ClientOptions.ConnectTimeout = &c.connectTimeout
	}



	var err error
	c.MClient, err = qmgo.NewClient(ctx, &qmgo.Config{
		Uri: c.dsn,
	}, clientOption)

	if err != nil {
		return fmt.Errorf(errInvalidDsn, c.dsn, err.Error())
	}

	return nil
}

func (c *Client) SetDsn(dsn string) {
	c.dsn = dsn
}

func (c *Client) SetConnectTimeout(v string) error {
	if connectTimeout, e := time.ParseDuration(v); e != nil {
		return fmt.Errorf(errSetProp, "connectTimeout", e.Error())
	} else {
		c.connectTimeout = connectTimeout
	}

	return nil
}

func (c *Client) SetReadTimeout(v string) error {
	if readTimeout, e := time.ParseDuration(v); e != nil {
		return fmt.Errorf(errSetProp, "readTimeout", e.Error())
	} else {
		c.readTimeout = readTimeout
	}

	return nil
}

func (c *Client) SetWriteTimeout(v string) error {
	if writeTimeout, e := time.ParseDuration(v); e != nil {
		return fmt.Errorf(errSetProp, "writeTimeout", e.Error())
	} else {
		c.writeTimeout = writeTimeout
	}

	return nil
}

func (c *Client) WriteTimeout() time.Duration{
	return c.writeTimeout
}

func (c *Client) ReadTimeout() time.Duration{
	return c.readTimeout
}