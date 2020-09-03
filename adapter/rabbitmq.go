package adapter

import (
	"time"

	"github.com/pinguo/pgo2"
	"github.com/pinguo/pgo2/client/rabbitmq"
	"github.com/pinguo/pgo2/iface"
	"github.com/pinguo/pgo2/util"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

var RabbitMqClass string

func init() {
	container := pgo2.App().Container()
	RabbitMqClass = container.Bind(&RabbitMq{})
}

type RabbitMq struct {
	pgo2.Object
	client       *rabbitmq.Client
	panicRecover bool
}

// NewRabbitMq of RabbitMq Client, add context support.
// usage: rabbitMq := this.GetObj(adapter.NewRabbitMq()).(adapter.IRabbitMq)/(*adapter.RabbitMq)
func NewRabbitMq(dftConfig ...string) *RabbitMq {

	id := DefaultRabbitId
	l := len(dftConfig)
	if l > 0 {
		id = dftConfig[0]
	}

	op:= map[string]interface{}{"logger": pgo2.GLogger()}
	if l >= 3 {
		op["user"] = dftConfig[1]
		op["pass"] = dftConfig[2]
	}

	r := &RabbitMq{}
	r.client = pgo2.App().Component(id, rabbitmq.New, op).(*rabbitmq.Client)
	r.panicRecover = true

	return r
}

// NewRabbitMqPool of RabbitMq Client from pool, add context support.
// usage: rabbitMq := this.GetObjPool(adapter.RabbitMqClass,adapter.NewRabbitMqPool).(adapter.IRabbitMq)/(*adapter.RabbitMq)
func NewRabbitMqPool(iObj iface.IObject, dftConfig ...interface{}) iface.IObject {

	id := DefaultRabbitId
	l := len(dftConfig)
	if l > 0 {
		id = dftConfig[0].(string)
	}

	op:= map[string]interface{}{"logger": pgo2.GLogger()}
	if l >= 3 {
		op["user"] = dftConfig[1]
		op["pass"] = dftConfig[2]
	}

	r := iObj.(*RabbitMq)
	r.client = pgo2.App().Component(id, rabbitmq.New, op).(*rabbitmq.Client)
	r.panicRecover = true

	return r
}

func (r *RabbitMq) SetPanicRecover(v bool) {
	r.panicRecover = v
}

func (r *RabbitMq) GetClient() *rabbitmq.Client {
	return r.client
}

func (r *RabbitMq) handlePanic() {
	if r.panicRecover {
		if v := recover(); v != nil {
			r.Context().Error(util.ToString(v))
		}
	}
}

func (r *RabbitMq) ExchangeDeclare(dftExchange ...*rabbitmq.ExchangeData) {
	profile := "rabbit.ExchangeDeclare"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()
	var exchange  *rabbitmq.ExchangeData
	if len(dftExchange) > 0 {
		exchange = dftExchange[0]
	}

	err := r.client.SetExchangeDeclare(exchange)
	panicErr(err)
}

func (r *RabbitMq) Publish(opCode string, data interface{}, dftOpUid ...string) bool {
	profile := "rabbit.Publish"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	opUid := ""
	if len(dftOpUid) > 0 {
		opUid = dftOpUid[0]
	}

	res, err := r.client.Publish(&rabbitmq.PublishData{ OpCode: opCode, Data: data, OpUid: opUid, DeliveryMode:amqp.Persistent}, r.Context().LogId())
	panicErr(err)

	return res
}

func (r *RabbitMq) PublishExchange(serviceName, exchangeName, exchangeType, opCode string, data interface{}, dftOpUid ...string) bool {
	profile := "rabbit.PublishExchange"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	opUid := ""
	if len(dftOpUid) > 0 {
		opUid = dftOpUid[0]
	}

	res, err := r.client.Publish(&rabbitmq.PublishData{ServiceName: serviceName,DeliveryMode:amqp.Persistent,
		ExChange: &rabbitmq.ExchangeData{Name: exchangeName, Type: exchangeType, Durable:true},
		OpCode:   opCode, Data: data, OpUid: opUid}, r.Context().LogId())
	panicErr(err)

	return res
}

func (r *RabbitMq) PublishExchangeDelivery(serviceName, exchangeName, exchangeType, opCode string, data interface{},deliveryMode uint8, dftOpUid ...string) bool {
	profile := "rabbit.PublishExchangeDelivery"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	opUid := ""
	if len(dftOpUid) > 0 {
		opUid = dftOpUid[0]
	}

	if deliveryMode > amqp.Persistent  || deliveryMode < 0 {
		panicErr(errors.New("Invalid deliveryMode"))
	}

	res, err := r.client.Publish(&rabbitmq.PublishData{ServiceName: serviceName,DeliveryMode:deliveryMode,
		ExChange: &rabbitmq.ExchangeData{Name: exchangeName, Type: exchangeType, Durable:true},
		OpCode:   opCode, Data: data, OpUid: opUid}, r.Context().LogId())
	panicErr(err)

	return res
}

func (r *RabbitMq) ChannelBox() *rabbitmq.ChannelBox {
	profile := "rabbit.ChannelBox"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	res, err := r.client.FreeChannel()
	panicErr(err)

	return res
}

func (r *RabbitMq) GetConsumeChannelBox(queueName string, opCodes []string, dftExchange ...*rabbitmq.ExchangeData) *rabbitmq.ChannelBox {
	profile := "rabbit.GetConsumeChannelBox"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	var exchange  *rabbitmq.ExchangeData
	if len(dftExchange) > 0 {
		exchange = dftExchange[0]
	}
	res, err := r.client.GetConsumeChannelBox(queueName, opCodes, exchange)
	panicErr(err)

	return res
}

// 消费，返回chan 可以不停取数据
// queueName 队列名字
// opCodes 绑定队列的code
// limit 每次接收多少条
// autoAck 是否自动答复 如果为false 需要手动调用Delivery.ack(false)
// noWait 是否一直等待
// exclusive 是否独占队列
func (r *RabbitMq) Consume(queueName string, opCodes []string, limit int, autoAck, noWait, exclusive bool) <-chan amqp.Delivery {
	profile := "rabbit.Consume"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	ConsumeName := queueName + "-" + time.Now().Format("20060102150405")
	res, err := r.client.Consume(&rabbitmq.ConsumeData{QueueName: queueName, Name:ConsumeName, OpCodes: opCodes, Limit: limit, AutoAck: autoAck, NoWait: noWait, Exclusive: exclusive})
	panicErr(err)

	return res
}

// 消费，返回chan 可以不停取数据
// exchangeName 交换机名
// exchangeType 交换机type
// queueName 队列名字
// opCodes 绑定队列的code
// limit 每次接收多少条
// autoAck 是否自动答复 如果为false 需要手动调用Delivery.ack(false)
// noWait 是否一直等待
// exclusive 是否独占队列
func (r *RabbitMq) ConsumeExchange(exchangeName, exchangeType, queueName string, opCodes []string, limit int, autoAck, noWait, exclusive bool) <-chan amqp.Delivery {
	profile := "rabbit.Consume"
	r.Context().ProfileStart(profile)
	defer r.Context().ProfileStop(profile)
	defer r.handlePanic()

	ConsumeName := queueName + "-" + time.Now().Format("20060102150405")
	res, err := r.client.Consume(&rabbitmq.ConsumeData{Name:ConsumeName, ExChange: &rabbitmq.ExchangeData{Name: exchangeName, Type: exchangeType, Durable:true},
		QueueName: queueName, OpCodes: opCodes, Limit: limit, AutoAck: autoAck, NoWait: noWait, Exclusive: exclusive})
	panicErr(err)

	return res
}

func (r *RabbitMq) DecodeBody(d amqp.Delivery, ret interface{}) error {
	return r.client.DecodeBody(d, ret)
}

func (r *RabbitMq) DecodeHeaders(d amqp.Delivery) *rabbitmq.RabbitHeaders {
	return r.client.DecodeHeaders(d)
}
