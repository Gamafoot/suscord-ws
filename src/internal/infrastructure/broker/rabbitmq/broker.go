package rabbitmq

import (
	"context"
	"encoding/json"
	"suscord_ws/internal/domain/eventbus"
	"time"

	pkgErrors "github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	eventbus eventbus.Bus
	opts     ConsumerOptions
}

func NewBroker(url string) (*Broker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	if err = exchangeDeclare(ch); err != nil {
		return nil, err
	}

	if err = queueDeclare(ch); err != nil {
		return nil, err
	}

	return &Broker{
		conn:    conn,
		channel: ch,
	}, nil
}

func (b *Broker) SetEventbus(eventbus eventbus.Bus) {
	b.eventbus = eventbus
}

func (b *Broker) Close() error {
	if err := b.channel.Close(); err != nil {
		return pkgErrors.WithStack(err)
	}

	if err := b.conn.Close(); err != nil {
		return pkgErrors.WithStack(err)
	}

	return nil
}

func (c *Broker) Consume(queue string) error {
	if c.eventbus == nil {
		return pkgErrors.New("eventbus не установлен")
	}

	msgs, err := c.channel.Consume(queue, c.opts.ConsumerTag, false, false, false, false, nil)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	for msg := range msgs {
		ctx, cansel := context.WithTimeout(context.Background(), 7*time.Second)
		if err := c.handleDelivery(ctx, msg); err != nil {
			cansel()
			return err
		}
		cansel()
	}

	return nil
}

func (c *Broker) handleDelivery(ctx context.Context, msg amqp.Delivery) error {
	if err := ctx.Err(); err != nil {
		_ = msg.Nack(false, true)
		return pkgErrors.WithStack(err)
	}

	body := new(MessageBody)
	if err := json.Unmarshal(msg.Body, body); err != nil {
		_ = msg.Nack(false, false)
		return nil
	}

	if err := c.eventbus.Publish(ctx, body.Type, eventbus.Payload(body.Payload)); err != nil {
		requeue := c.opts.DefaultRequeue
		if action, ok := NackActionFromError(err); ok {
			requeue = action == NackRequeue
		}
		_ = msg.Nack(false, requeue)
		return nil
	}

	if err := msg.Ack(false); err != nil {
		return pkgErrors.WithStack(err)
	}

	return nil
}

func exchangeDeclare(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"chat.events",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	return nil
}

func queueDeclare(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(
		"chat.ws",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	err = ch.QueueBind(
		q.Name,
		"chat.ws",
		"chat.events",
		false,
		nil,
	)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	return nil
}
