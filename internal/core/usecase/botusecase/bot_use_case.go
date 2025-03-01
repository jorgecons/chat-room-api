package botusecase

import (
	"context"
	"fmt"

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

func (uc *UseCase) Bot(ctx context.Context, message domain.Message) error {
	fmt.Println(message)
	price, err := uc.stockRepo.GetPrice(ctx, message.Text)
	if err != nil {
		return err
	}
	fmt.Println("price ", price)
	message.Text = domain.CreateBotMessage(message.Text, price)

	return uc.publisher.Publish(ctx, message, domain.BotMessageType)
}
