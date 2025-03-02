package botusecase

import (
	"context"
	"time"

	"chat-room-api/internal/core/domain"
)

type (
	UseCase struct {
		stockRepo stockRepository
		publisher Publisher
	}

	stockRepository interface {
		GetPrice(context.Context, string) (float64, error)
	}

	Publisher interface {
		Publish(context.Context, domain.Message, string) error
	}
)

func NewUseCase(
	sr stockRepository,
	p Publisher,
) *UseCase {
	return &UseCase{
		stockRepo: sr,
		publisher: p,
	}
}

func (uc *UseCase) Bot(ctx context.Context, msg domain.Message) error {
	var (
		text      string
		stockName = domain.GetStockName(msg.Text)
	)
	price, err := uc.stockRepo.GetPrice(ctx, stockName)
	if err != nil {
		text = err.Error()
	} else {
		text = domain.CreateBotMessage(msg.Text, price)
	}
	message := domain.NewMessage(msg.Room, domain.BotUsername, text, time.Now().UTC())
	return uc.publisher.Publish(ctx, message, domain.BotMessageType)
}
