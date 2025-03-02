package stockusecase

import (
	"context"
	"time"

	"chat-room-api/internal/core/domain"
)

type (
	UseCase struct {
		stockPublisher stockPublisher
	}

	stockPublisher interface {
		Publish(context.Context, domain.Message, string) error
	}
)

func NewUseCase(
	publisher stockPublisher,
) *UseCase {
	return &UseCase{
		stockPublisher: publisher,
	}
}

func (uc *UseCase) Stock(ctx context.Context, room, stock string) error {
	msg := domain.Message{
		Room:     room,
		Username: "bot",
		Text:     stock,
		Date:     time.Now(),
	}
	return uc.stockPublisher.Publish(ctx, msg, domain.StockMessageType)
}
