package rabbitmq

import (
    "errors"
    "time"

    "github.com/pinguo/pgo2"
    "github.com/streadway/amqp"
)

func newChannelBox(connBox *ConnBox, pool *Pool) (*ChannelBox, error) {
    connBox.newChannelLock.Lock()
    defer connBox.newChannelLock.Unlock()
    channel, err := connBox.connection.Channel()

    if err != nil {
        return nil, errors.New("Rabbit newChannelBox err:" + err.Error())
    }
    connBox.useConnCount++
    return &ChannelBox{connBoxId: connBox.id, pool: pool, channel: channel, connStartTime: connBox.startTime, lastActive: time.Now()}, nil
}

type ChannelBox struct {
    connBoxId     string
    pool          *Pool
    channel       *amqp.Channel
    connStartTime time.Time
    lastActive    time.Time
}

func (c *ChannelBox) Close(force bool) error {
    conn, err := c.pool.getConnBox(c.connBoxId)
    if err != nil {
        return err
    }

    if force || c.connStartTime != conn.startTime {
        return c.channelClose()
    }

    ret, err := c.pool.putFreeChannel(c)
    if err != nil {
        return err
    }

    if !ret {
        return c.channelClose()
    } else {
        c.lastActive = time.Now()
    }

    return nil
}

func (c *ChannelBox) channelClose() error {
    connBox, err := c.pool.getConnBox(c.connBoxId)
    if err != nil {
        return err
    }

    connBox.useConnCount--
    if connBox.isClosed() == false {
        err := c.channel.Close()
        if err != nil {
            pgo2.GLogger().Warn("Rabbit ChannelBox.channelClose.channel.Close() err:" + err.Error())
        }
    }

    return nil
}

func (c *ChannelBox) GetChannel() *amqp.Channel {
    return c.channel
}
