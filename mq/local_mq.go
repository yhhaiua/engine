package mq

import (
	"github.com/streadway/amqp"
	"github.com/yhhaiua/engine/log"
	"strings"
)

var logger = log.GetLogger()

// NewConnection 建立mq新连接
func NewConnection(config *RabbitConfig) *amqp.Connection {
	b := &strings.Builder{}
	b.WriteString("amqp://")
	b.WriteString(config.UserName)
	b.WriteString(":")
	b.WriteString(config.PassWord)
	b.WriteString("@")
	b.WriteString(config.Path)
	b.WriteString(config.VirtualHost)
	conn, err := amqp.Dial(b.String())
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		return nil
	}
	return conn
}

// CreateChannel 在连接基础上创建一个通道
func CreateChannel(connect *amqp.Connection) *amqp.Channel {
	if connect == nil {
		return nil
	}
	ch, err := connect.Channel()
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		return nil
	}
	return ch
}

// ExchangeDeclare 创建和连接交换机
func ExchangeDeclare(channel *amqp.Channel, name string) {
	if channel == nil {
		return
	}
	err := channel.ExchangeDeclare(name, "topic", true, false, false, false, nil)
	if err != nil {
		logger.Errorf("error:%s", err.Error())
	}
}

func QueueDeclare(channel *amqp.Channel, queueName string) {
	if channel == nil {
		return
	}
	_, err := channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		logger.Errorf("error:%s", err.Error())
	}
}
func QueueBind(channel *amqp.Channel, queueName, key string, exchangeName string) {
	if channel == nil {
		return
	}
	err := channel.QueueBind(
		queueName,    // name of the queue
		key,          // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		logger.Errorf("error:%s", err.Error())
	}
}

func Consume(channel *amqp.Channel, queueName string) <-chan amqp.Delivery {
	if channel == nil {
		return nil
	}
	deliveries, err := channel.Consume(
		queueName, // name
		"",        // consumerTag,
		false,     // noAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		logger.Errorf("error:%s", err.Error())
	}
	return deliveries
}
