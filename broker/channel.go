// date: 2019-03-18
package broker

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/streadway/amqp"
)

type rabbitMQChannel struct {
	uuid       string
	connection *amqp.Connection
	channel    *amqp.Channel
}

//实例化连接
func newRabbitChannel(conn *amqp.Connection, prefetchCount int, prefetchGlobal bool) (*rabbitMQChannel, error) {
	id, err := uuid.NewV1()
	if err != nil {
		return nil, err
	}

	rabbitCh := &rabbitMQChannel{uuid: id.String(), connection: conn}

	if err := rabbitCh.Connect(prefetchCount, prefetchGlobal); err != nil {
		return nil, err
	}
	return rabbitCh, nil
}

//连接
func (r *rabbitMQChannel) Connect(prefetchCount int, prefetchGlobal bool) error {
	var err error
	r.channel, err = r.connection.Channel()
	if err != nil {
		return err
	}

	err = r.channel.Qos(prefetchCount, 0, prefetchGlobal)
	if err != nil {
		return err
	}
	return nil
}

//关闭channel
func (r *rabbitMQChannel) Close() error {
	if r.channel == nil {
		return errors.New("channel is nil")
	}
	return r.channel.Close()
}

//发布消息
func (r *rabbitMQChannel) Publish(exchange, key string, message amqp.Publishing) error {
	if r.channel == nil {
		return errors.New("channel is nil")
	}
	return r.channel.Publish(exchange, key, false, false, message)
}

//声明交换机
func (r *rabbitMQChannel) DeclareExchange(exchange string) error {
	return r.channel.ExchangeDeclare(exchange,
		"topic",
		false,
		false,
		false,
		false,
		nil)
}

//声明持久化交换机
func (r *rabbitMQChannel) DeclareDurableExchange(exchange string) error {
	return r.channel.ExchangeDeclare(exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil)
}

//声明队列
func (r *rabbitMQChannel) DeclareQueue(queue string, args amqp.Table) error {
	_, err := r.channel.QueueDeclare(queue,
		false,
		true,
		false,
		false,
		args)
	return err
}

//声明持久化队列
func (r *rabbitMQChannel) DeclareDurableQueue(queue string, args amqp.Table) error {
	_, err := r.channel.QueueDeclare(queue,
		true,
		true,
		false,
		false,
		args)
	return err
}

//声明重试队列
func (r *rabbitMQChannel) DeclareReplyQueue(queue string) error {
	_, err := r.channel.QueueDeclare(queue,
		false,
		true,
		true,
		false,
		nil)
	return err
}

//消费消息队列
func (r *rabbitMQChannel) ConsumeQueue(queue string, autoAck bool) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(queue,
		r.uuid,
		autoAck,
		false,
		false,
		false,
		nil)
}

//绑定消息队列
func (r *rabbitMQChannel) BindQueue(queue, key, exchange string, args amqp.Table) error {
	return r.channel.QueueBind(queue, key, exchange, false, args)
}
