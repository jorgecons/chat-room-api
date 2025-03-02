package chatusecase

import (
	"context"
	"strings"

	"chat-room-api/internal/core/domain"

	"github.com/sirupsen/logrus"
)

var zeroMessage domain.Message

type (
	UseCase struct {
		messageStorage messageRepository
		stockPublisher stockPublisher
	}

	stockPublisher interface {
		Publish(context.Context, domain.Message, string) error
	}

	messageRepository interface {
		Save(context.Context, domain.Message) error
	}
)

func NewUseCase(
	messageStorage messageRepository,
	publisher stockPublisher,
) *UseCase {
	return &UseCase{
		messageStorage: messageStorage,
		stockPublisher: publisher,
	}
}

func (uc *UseCase) Chat(ctx context.Context, msg domain.Message) error {
	if strings.HasPrefix(msg.Text, "/") {
		if !strings.HasPrefix(msg.Text, domain.StockPrefix) {
			logrus.WithContext(ctx).
				WithField("command", msg.Text).
				WithField("room", msg.Room).
				WithField("username", msg.Username).
				Error("Error validating command")
			return domain.ErrUnknownCommand
		}
		if err := uc.stockPublisher.Publish(ctx, msg, domain.StockMessageType); err != nil {
			logrus.WithContext(ctx).
				WithError(err).
				WithField("command", msg.Text).
				WithField("room", msg.Room).
				WithField("username", msg.Username).
				Error("Error publishing command")
			return domain.WrapError(domain.ErrChatting, err)
		}
		return nil

	}
	if err := uc.messageStorage.Save(ctx, msg); err != nil {
		logrus.WithContext(ctx).
			WithError(err).
			WithField("room", msg.Room).
			WithField("username", msg.Username).
			Error("Error saving message")
		return domain.WrapError(domain.ErrChatting, err)
	}
	return nil
}
