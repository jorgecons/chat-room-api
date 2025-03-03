package application

import (
	"context"
	"encoding/json"

	"chat-room-api/internal/core/adapter/sockethandler"
	"chat-room-api/internal/core/domain"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func (a *App) MapEventRoutes() *App {
	a.queueConsumer.consumerFn = func(ctx context.Context, forever chan bool, msgs <-chan amqp.Delivery) {
		for {
			select {
			case <-ctx.Done():
				logrus.Warnln("Timeout reached, stopping consumer.")
				forever <- false
				return
			case d, ok := <-msgs:
				if !ok {
					logrus.Warnln("Channel closed, stopping consumer.")
					forever <- false
					return
				}
				switch d.Type {
				case domain.StockMessageType:
					err := a.handlers.ConsumeBotMessageHandler(ctx, d.Body)
					if err != nil {
						_ = d.Nack(false, true)
						return
					}
					_ = d.Ack(false)
				case domain.BotMessageType:
					msg := sockethandler.Message{}
					_ = json.Unmarshal(d.Body, &msg)
					a.socket.broadcast <- msg
					_ = d.Ack(false)
				default:
					_ = d.Reject(false)
					return
				}

			}
		}
	}
	return a
}
