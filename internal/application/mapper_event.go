package application

import (
	"context"
	"encoding/json"
	"time"

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
				if d.Type == domain.StockMessageType {
					err := a.handlers.ConsumeBotMessageHandler(ctx, d.Body)
					if err != nil {
						logrus.Warnln("Channel closed, stopping consumer.")
					}
				}
				if d.Type == domain.BotMessageType {
					msg := sockethandler.Message{}
					_ = json.Unmarshal(d.Body, &msg)
					msg.Date = time.Now()
					a.socket.broadcast <- msg
				}

			}
		}
	}
	return a
}
