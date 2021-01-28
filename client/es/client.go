package es

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pinguo/pgo2/client/phttp"
	"github.com/pinguo/pgo2/core"
	"github.com/pinguo/pgo2/perror"
	"github.com/pinguo/pgo2/util"

	"github.com/pinguo/pgo2"
)

// 数据结构
type DataItem struct {
	Head string
	Body string
}

// Elasticsearch client component,Support for batch and asynchronous writes.
// support read-write splitting, configuration:
// components:
//      es:
//          esHost: "http://127.0.0.1:9200/"
//          httpId:    "http" // 默认http
//          retryNum: 0 // 默认0重试1次
//          batchChanLen: 500 // 队列数量 默认500
//          batchFlushInterval: "50s" // 批量操作刷新间隔时间 默认50秒
//          batchEsTimeout:"30s" // 批量请求es服务器超时时间 默认30秒
//          singleEsTimeout: "500ms" // 单个请求es服务器超时时间 默认500毫秒
//          maxBufferLine: 600 // 缓冲最条数,默认600
//          maxBufferByte: 1048576 // 写缓冲大小 单位bite 默认1M
func New(config map[string]interface{}) (interface{}, error) {
	c := &Client{}
	c.batchChanLen = 500
	// 50s 用于和log刷新间隔
	c.batchFlushInterval = 50 * time.Second
	c.maxBufferByte = 1 * 1024 * 1024
	c.maxBufferLine = 600
	c.httpId = "http"
	c.batchEsTimeout = 30000 * time.Millisecond
	c.singleEsTimeout = 500 * time.Millisecond
	if err := core.ClientConfigure(c, config); err != nil {
		return nil, err
	}

	if err := c.Init(); err != nil {
		return nil, err
	}
	return c, nil
}

type Client struct {
	// http 组建id
	httpId string
	// es host
	esHost string
	// 网络重试次数
	retryNum int
	// 队列数量
	batchChanLen int
	// 批量操作刷新间隔时间 单位：秒
	batchFlushInterval time.Duration
	// 批量请求es服务器超时时间
	batchEsTimeout time.Duration
	// 单个请求es服务器超时时间
	singleEsTimeout time.Duration

	// 缓冲最条数
	maxBufferLine int
	// 写缓冲大小
	maxBufferByte int

	// 队列
	chDataChan chan *DataItem
	// 保存缓冲区状态
	wg sync.WaitGroup
	// 准备写缓冲区
	buffer        bytes.Buffer
	curBufferLine int

}

func (c *Client) Init() error {
	c.chDataChan = make(chan *DataItem, c.batchChanLen)
	c.wg.Add(1)
	// start loop
	go c.loop()
	//
	pgo2.App().StopBefore().Add(c, "Close")
	return nil
}


func (c *Client) SetEsHost(v string) {
	c.esHost = v
}

func (c *Client) SetRetryNum(v int) {
	c.retryNum = v
}

func (c *Client) SetMaxBufferLine(v int) {
	c.batchChanLen = v
}

func (c *Client) SetBatchFlushInterval(v string) error{
	if flushInterval, err := time.ParseDuration(v); err != nil {
		return errors.New(fmt.Sprintf("Log: parse flushInterval error, val:%s, err:%s", v, err.Error()))
	} else {
		c.batchFlushInterval = flushInterval
	}
	return nil
}

func (c *Client) SetBatchEsTimeout(v string) error{
	if timeout, err := time.ParseDuration(v); err != nil {
		return errors.New(fmt.Sprintf("Log: parse batchEsTimeout error, val:%s, err:%s", v, err.Error()))
	} else {
		c.batchEsTimeout = timeout
	}
	return nil
}

func (c *Client) SetSingleEsTimeout(v string) error{
	if timeout, err := time.ParseDuration(v); err != nil {
		return errors.New(fmt.Sprintf("Log: parse singleEsTimeout error, val:%s, err:%s", v, err.Error()))
	} else {
		c.singleEsTimeout = timeout
	}
	return nil
}

func (c *Client) SetBatchFlushNum(v int) {
	c.batchChanLen = v
}

func (c *Client) SetMaxBufferByte(v int) {
	c.maxBufferByte = v
}

// 循环检查
func (c *Client) loop() {
	flushTimer := time.Tick(c.batchFlushInterval)

	for {
		select {
		case item, ok := <-c.chDataChan:
			if ok {
				// 有数据
				c.Process(item)
			} else {
				// 通道被关闭 刷新数据
				c.Flush()
			}

			if !ok {
				goto end
			}
		case <-flushTimer:
			// 定时刷新
			c.Flush()
		}
	}

end:
	c.wg.Done()
}

// 关闭通道，等待处理数据
func (c *Client) Close() {
	close(c.chDataChan)
	c.wg.Wait()
}

func (c *Client) Process(item *DataItem) {
	// write data to buffer
	c.buffer.WriteString(item.Head + "\n" + item.Body + "\n")
	c.curBufferLine++
	// 行数和大小大道限制 刷盘
	if c.curBufferLine >= c.maxBufferLine || c.buffer.Len() >= c.maxBufferByte {
		c.Flush()
	}
}

// es 批量返回结构
type esBatchRet struct {
	Errors interface{} `json:"errors"`
}

// 刷盘
func (c *Client) Flush() {
	if c.curBufferLine == 0 {
		return
	}
	buffer := c.buffer
	c.buffer.Reset()
	flushNum := c.curBufferLine
	c.curBufferLine = 0

	//重试, 只处理IO错误，不处理逻辑错误 ,可能丢操作数据
	sTime := time.Now().UnixNano() / 1e6
	var errEs error
	var i = 0
	for i = 0; i <= c.retryNum; i++ {
		var ret *esBatchRet
		content, err := c.Request("POST", "_bulk", buffer.Bytes(), c.batchEsTimeout, true)
		if err != nil {
			errEs = fmt.Errorf("post es network err:" + err.Error())
		}
		err = json.Unmarshal(content, &ret)
		if err != nil {
			// 继续重试
			continue
		}

		if errors, ok := ret.Errors.(bool); ok == false || errors != false {
			// es内部错误，继续重试, 部分操作具体错误数
			errEs = fmt.Errorf("es inner err:" + util.ToString(content))
			continue
		}

		break
	}

	//
	if errEs != nil {
		pgo2.App().Log().Logger(pgo2.App().Name(), "es_flush").Error("es_request_bulk err:" + errEs.Error())
	}

	eTime := time.Now().UnixNano() / 1e6
	pgo2.App().Log().Logger(pgo2.App().Name(), "es_flush").Info("num:" + util.ToString(flushNum) + " " + c.esHost + "=" + util.ToString(eTime-sTime) + "ms/" + util.ToString(i+1))

}

// 批量增加命令  异步执行
// action :  index，create，delete, update
// head : {“_ index”：“test”，“_ type”：“_ doc”，“_ id”：“1”}
// body : {"filed1":"value1"}
func (c *Client) Batch(action, head, body string) error{
	if body == "" && action != batchActionDelete {
		return perror.New(http.StatusBadGateway, "ES The lack of the body")
	}
	if action == batchActionDelete && body != "" {
		return perror.New(http.StatusBadGateway, "ES Delete unsupported body")
	}
	head = "{\"" + action + "\":" + head + "}"
	c.chDataChan <- &DataItem{Head: head, Body: body}
	return nil
}

// 单个提交
func (c *Client) Single(method, uri string, body []byte, timeout time.Duration) ([]byte, error) {
	var content []byte
	var err error
	for i := 0; i <= c.retryNum; i++ {
		content, err = c.Request(method, uri, body, timeout)

		if err != nil || content == nil {
			continue
		}

		break
	}
	return content, err
}

// Request 请求es，返回原始结果
// method: POST GET PUT ....
// body: 请求body
// timeout: 超时设置
// bulk: 是否批量操作(默认否)
func (c *Client) Request(method, uri string, body []byte, timeout time.Duration, bulk ...bool) ([]byte, error) {

	httpClient := pgo2.App().Component(c.httpId, phttp.New).(*phttp.Client)
	option := phttp.Option{}
	url := strings.TrimRight(c.esHost, "/") + "/" + strings.TrimLeft(uri, "/")

	option.SetTimeout(timeout)
	if len(bulk) == 1 && bulk[0] {
		option.SetHeader("Content-Type", "application/x-ndjson")
	} else {
		option.SetHeader("Content-Type", "application/json")
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Do(req, &option);
	if err != nil {
		return nil, err
	}
	if res != nil {
		content, _ := ioutil.ReadAll(res.Body)
		_ = res.Body.Close()
		return content, nil
	}

	return nil, nil
}
