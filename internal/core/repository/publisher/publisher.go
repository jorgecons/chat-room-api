package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"chat-room-api/internal/core/domain"

	amqp "github.com/rabbitmq/amqp091-go"
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
	fmt.Println("sending", body)
	err := p.ch.Publish(
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Type:        messageType,
		},
	)
	if err != nil {
		fmt.Println("error ", err)
	}
	return err
}
