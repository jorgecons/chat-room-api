package publisher

import (
	"context"
	"encoding/json"

	"chat-room-api/internal/core/domain"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type (
	Publisher struct {
		ch        *amqp.Channel
		queueName string
	}
)

func NewPublisher(ch *amqp.Channel, queueName string) *Publisher {
	return &Publisher{
		ch:        ch,
		queueName: queueName,
	}
}

func (p *Publisher) Publish(ctx context.Context, msg domain.Message, messageType string) error {
	body, _ := json.Marshal(msg)
	err := p.ch.Publish(
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: jsonContentTypeHeader,
			Body:        body,
			Type:        messageType,
		},
	)
	if err != nil {
		logrus.WithContext(ctx).
			WithError(err).
			WithField("room", msg.Room).
			WithField("stock", msg.Text).
			Error("Error publishing message")
		return domain.WrapError(domain.ErrPublishingMessage, err)
	}
	return nil
}
