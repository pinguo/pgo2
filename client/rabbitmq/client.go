package rabbitmq

import (
    "bytes"
    "encoding/gob"
    "errors"
    "time"

    "github.com/pinguo/pgo2/core"
    "github.com/streadway/amqp"
)

// RabbitMq client component,
// support Publisher-Consumer  configuration:
// components:
//      rabbitMq:
//          tlsRootCAs:""
//          tlsCert: ""
//          tlsCertKey: ""
//          user: "guest"
//          pass: "guest"
//          exchangeName: ""
//          exchangeType: ""
//          maxChannelNum: 2000
//          maxIdleChannel: "200"
//          maxIdleChannelTime:"10s"
//          probeInterval: "0s"
//          maxWaitTime: "200ms"
//          serverName: "pgo-xxx"
//          servers:
//              - "127.0.0.1:6379"
//              - "127.0.0.1:6380"
func New(config map[string]interface{}) (interface{}, error) {
    c := &Client{}
    c.connList = make(map[string]*ConnBox)
    c.servers = make(map[string]*serverInfo)
    c.maxChannelNum = dftMaxChannelNum
    c.maxIdleChannel = dftMaxIdleChannel
    c.maxIdleChannelTime = dftMaxIdleChannelTime
    c.exchangeType = dftExchangeType
    c.exchangeName = dftExchangeName
    c.maxWaitTime = dftMaxWaitTime
    c.probeInterval = dftProbeInterval

    if err := core.ClientConfigure(c, config); err != nil {
        return nil, err
    }

    if err := c.Init(); err != nil {
        return nil, err
    }

    return c, nil

}

type Client struct {
    Pool
}

func (c *Client) DecodeBody(d amqp.Delivery, ret interface{}) error {
    var network *bytes.Buffer
    network = bytes.NewBuffer(d.Body)
    dec := gob.NewDecoder(network)
    err := dec.Decode(ret)

    return err
}

func (c *Client) DecodeHeaders(d amqp.Delivery) *RabbitHeaders {
    ret := &RabbitHeaders{
        Exchange:  d.Exchange,
        RouteKey:  d.RoutingKey,
        Timestamp: d.Timestamp,
        MessageId: d.MessageId,
    }

    for k, iV := range d.Headers {
        v, _ := iV.(string)
        switch k {
        case "logId":
            ret.LogId = v
        case "service":
            ret.Service = v
        case "opUid":
            ret.OpUid = v
        }
    }

    return ret
}

func (c *Client) SetExchangeDeclare() error {
    ch, err := c.getFreeChannel()
    if err != nil {
        return err
    }
    defer ch.Close(false)

    return c.exchangeDeclare(ch)

}

func (c *Client) Publish(parameter *PublishData, logId string) (bool, error) {
    if parameter.OpCode == "" || parameter.Data == nil {
        return false, errors.New("Rabbit OpCode and LogId cannot be empty")
    }

    ch, err := c.getFreeChannel()
    if err != nil {
        return false, err
    }
    defer ch.Close(false)

    // 增加速度，在消费端定义交换机 或者单独定义交换机
    // c.exchangeDeclare(ch)

    var goBytes bytes.Buffer
    myGob := gob.NewEncoder(&goBytes)
    err = myGob.Encode(parameter.Data)
    if err != nil {
        return false, c.failOnError(err, "Encode err")
    }

    err = ch.channel.Publish(
        c.getExchangeName(),             // exchange
        c.getRouteKey(parameter.OpCode), // routing key
        false,                           // mandatory
        false,                           // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        goBytes.Bytes(),
            Headers:     amqp.Table{"logId": logId, "service": c.ServiceName, "opUid": parameter.OpUid},
            Timestamp:   time.Now(),
        })
    if err != nil {
        return false, c.failOnError(err, "Failed to publish a message")
    }

    return true, nil
}

// 定义交换机
func (c *Client) exchangeDeclare(ch *ChannelBox) error {
    err := ch.channel.ExchangeDeclare(
        c.getExchangeName(), // name
        c.exchangeType,      // type
        true,                // durable
        false,               // auto-deleted
        false,               // internal
        false,               // no-wait
        nil,                 // arguments
    )
    if err != nil {
        return c.failOnError(err, "Failed to declare an exchange")
    }

    return nil
}

// 定义交换机
func (c *Client) bindQueue(ch *ChannelBox, queueName string, opCodes []string) error {
    for _, opCode := range opCodes {
        err := ch.channel.QueueBind(
            queueName,             // queue name
            c.getRouteKey(opCode), // routing key
            c.getExchangeName(),   // exchange
            false,
            nil)
        if err != nil {
            return c.failOnError(err, "Failed to bind a queue")
        }

    }

    return nil
}

func (c *Client) queueDeclare(ch *ChannelBox, queueName string) (amqp.Queue, error) {
    q, err := ch.channel.QueueDeclare(
        queueName, // name
        true,      // durable
        false,     // delete when usused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )

    if err != nil {
        return q, c.failOnError(err, "Failed to declare a queue")
    }

    return q, nil
}

func (c *Client) GetConsumeChannelBox(queueName string, opCodes []string) (*ChannelBox, error) {
    ch, err := c.getFreeChannel()
    if err != nil {
        return nil, err
    }
    // 定义交换机
    if err := c.exchangeDeclare(ch); err != nil {
        return nil, err
    }
    // 定义queue
    if _, err := c.queueDeclare(ch, queueName); err != nil {
        return nil, err
    }
    // 绑定queue
    if err := c.bindQueue(ch, queueName, opCodes); err != nil {
        return nil, err
    }

    return ch, nil

}

func (c *Client) Consume(parameter *ConsumeData) (<-chan amqp.Delivery, error) {
    ch, err := c.GetConsumeChannelBox(parameter.QueueName, parameter.OpCodes)
    if err != nil {
        return nil, err
    }
    // defer ch.Close(false)

    if err := ch.channel.Qos(parameter.Limit, 0, false); err != nil {
        return nil, c.failOnError(err, "set Qos err")
    }

    messages, err := ch.channel.Consume(
        parameter.QueueName, // queue
        "",                  // consumer
        parameter.AutoAck,   // auto ack
        parameter.Exclusive, // exclusive
        false,               // no local
        parameter.NoWait,    // no wait
        nil,                 // args
    )
    if err != nil {
        return nil, c.failOnError(err, "get msg err")
    }

    return messages, nil
}

func (c *Client) failOnError(err error, msg string) error {
    return errors.New("Rabbit:" + msg + ",err:" + err.Error())

}
