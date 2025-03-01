package chatusecase

import (
	"context"
	"fmt"

	"chat-room-api/internal/core/domain"
)

var zeroMessage domain.Message

type (
	UseCase struct {
		messageStorage messageRepository
	}

	messageRepository interface {
		Save(context.Context, domain.Message) error
	}
)

func NewUseCase(
	messageStorage messageRepository,
) *UseCase {
	return &UseCase{
		messageStorage: messageStorage,
	}
}

func (uc *UseCase) Chat(ctx context.Context, msg domain.Message) error {
	fmt.Println(msg)
	return uc.messageStorage.Save(ctx, msg)
}
