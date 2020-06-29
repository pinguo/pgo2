package rabbitmq

import (
	"time"

	"github.com/streadway/amqp"
)

const (
	dftExchangeType       = "direct"
	dftExchangeName       = "direct_pgo2_dft"
	dftMaxChannelNum      = 2000
	dftMaxIdleChannel     = 200
	dftMaxIdleChannelTime = 60 * time.Second
	dftMaxWaitTime        = 100 * time.Microsecond
	dftProbeInterval      = 0
	dftProtocol           = "amqp"
	defaultTimeout        = 1 * time.Second
	errSetProp            = "rabbitMq: failed to set %s, %s"
)

type RabbitHeaders struct {
	LogId     string
	Exchange  string
	RouteKey  string
	Service   string
	OpUid     string
	Timestamp time.Time
	MessageId string
}

type ExchangeData struct {
	Name string // 交换机名
	Type string // 交换机类型
	Durable bool
	AutoDelete bool
	Internal bool
	NoWait bool
	Args amqp.Table
}

// rabbit 发布结构
type PublishData struct {
	ServiceName string // 服务名
	ExChange *ExchangeData
	OpCode string      // 操作code 和queue绑定相关
	OpUid  string      // 操作用户id 可以为空
	ContentType string // 内容类型 默认为："text/plain"
	Data   interface{} // 发送数据
}

type ConsumeData struct {
	ExChange *ExchangeData
	QueueName string
	OpCodes   []string
	AutoAck   bool
	NoWait    bool
	Exclusive bool
	Limit     int
}
