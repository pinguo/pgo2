package phttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/pinguo/pgo2/core"
	"github.com/pinguo/pgo2/util"
)

// Http Client component, configuration:
// components:
//      http:
//          verifyPeer: false
//          userAgent: "PGO Framework"
//          timeout: "10s"
func New(config map[string]interface{}) (interface{}, error) {
	client := &Client{
		verifyPeer: false,
		userAgent:  defaultUserAgent,
		timeout:    defaultTimeout,
	}

	err := core.ClientConfigure(client, config)
	if err != nil {
		return nil, err
	}
	client.Init()

	return client, nil
}

type Client struct {
	verifyPeer bool          // verify https peer or not
	userAgent  string        // default User-Agent header
	timeout    time.Duration // default request timeout

	client *http.Client
}

func (c *Client) Init() {
	// reused client and transport, transport will cache
	// connections for future reuse, if transport created
	// on demand, net poll goroutine on connection per
	// transport will be leaked.
	c.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !c.verifyPeer,
			},
		},
	}
}

func (c *Client) SetVerifyPeer(verifyPeer bool) {
	c.verifyPeer = verifyPeer
}

func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

func (c *Client) SetTimeout(v string) error {
	if timeout, err := time.ParseDuration(v); err != nil {
		return fmt.Errorf("http parse timeout failed, " + err.Error())
	} else {
		c.timeout = timeout
	}

	return nil
}

// Get perform a get request, and return a response pointer.
// addr is the request url. data is the params associated
// and will be append to addr if not empty, data type can be
// url.Values or map with string key. option is an optional
// configuration object to specify header, cookie etc.
func (c *Client) Get(addr string, data interface{}, option ...*Option) (*http.Response, error) {
	var query url.Values
	switch v := data.(type) {
	case nil:
		// pass

	case url.Values:
		query = v

	case map[string]interface{}, map[string]string, core.Map:
		query = make(url.Values)
		rv := reflect.ValueOf(data)
		keys := rv.MapKeys()
		for _, key := range keys {
			val := rv.MapIndex(key)
			query.Set(key.String(), util.ToString(val.Interface()))
		}

	default:
		return nil, fmt.Errorf("http get invalid data type: %T", data)
	}

	if len(query) != 0 {
		if pos := strings.IndexByte(addr, '?'); pos == -1 {
			addr = addr + "?" + query.Encode()
		} else {
			addr = addr + "&" + query.Encode()
		}
	}

	req, err := c.HttpNewRequest("GET", addr, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, option...)
}

func (c *Client) HttpNewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

// Post perform a post request, and return a response pointer.
// addr is the request url. data is the params will be sent,
// data type can be url.Values, map with string key, string,
// []byte or io.Reader, if url.Values or map is specified,
// Content-Type header will be set to "application/x-www-form-urlencoded".
// option is an optional configuration object to specify header, cookie etc.
func (c *Client) Post(addr string, data interface{}, option ...*Option) (*http.Response, error) {
	var body io.Reader
	var contentType string

	switch v := data.(type) {
	case nil:
		// pass

	case url.Values:
		body = strings.NewReader(v.Encode())
		contentType = "application/x-www-form-urlencoded"

	case map[string]interface{}, map[string]string, core.Map:
		query, rv := make(url.Values), reflect.ValueOf(data)
		keys := rv.MapKeys()
		for _, key := range keys {
			val := rv.MapIndex(key)
			query.Set(key.String(), util.ToString(val.Interface()))
		}

		body = strings.NewReader(query.Encode())
		contentType = "application/x-www-form-urlencoded"

	case string:
		body = strings.NewReader(v)

	case []byte:
		body = bytes.NewReader(v)

	case io.Reader:
		body = v

	default:
		return nil, fmt.Errorf("http post invalid data type: %T", data)
	}

	req, err := c.HttpNewRequest("POST", addr, body)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return c.Do(req, option...)
}

// Do perform a request specified by req param, and return response pointer.
func (c *Client) Do(req *http.Request, option ...*Option) (*http.Response, error) {
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	timeout := c.timeout
	if len(option) > 0 && option[0] != nil {
		opt := option[0]
		if opt.Timeout > 0 {
			timeout = opt.Timeout
		}

		for key, val := range opt.Header {
			if len(val) > 0 {
				req.Header.Set(key, val[0])
			}
		}

		for _, cookie := range opt.Cookies {
			req.AddCookie(cookie)
		}
	}

	ctx, _ := context.WithTimeout(req.Context(), timeout)

	return c.ClientDo(req.WithContext(ctx))
}

func (c *Client) ClientDo(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
