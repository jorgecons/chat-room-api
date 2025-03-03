package connectchatusecase

import (
	"context"

	"github.com/sirupsen/logrus"

	"chat-room-api/internal/core/domain"
)

var zeroMessages []domain.Message

type (
	UseCase struct {
		messageStorage GetLastMessages
	}

	GetLastMessages interface {
		GetLastMessages(context.Context, string) ([]domain.Message, error)
	}
)

func NewUseCase(
	messageStorage GetLastMessages,
) *UseCase {
	return &UseCase{
		messageStorage: messageStorage,
	}
}

func (uc *UseCase) ConnectChat(ctx context.Context, room string) ([]domain.Message, error) {
	lm, err := uc.messageStorage.GetLastMessages(ctx, room)
	if err != nil {
		err = domain.WrapError(domain.ErrConnectingChat, err)
		logrus.WithContext(ctx).WithError(err).WithField("room", room).Error("Error getting messages")
		return zeroMessages, err
	}
	res := reverseMessages(lm)
	return res, nil
}

func reverseMessages(lm []domain.Message) []domain.Message {
	res := make([]domain.Message, len(lm))
	for i, j := len(lm)-1, 0; i >= 0 && j < len(lm); i, j = i-1, j+1 {
		res[j] = lm[i]
	}
	return res
}
