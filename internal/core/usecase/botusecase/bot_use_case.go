package botusecase

import (
	"context"
	"fmt"
	"time"

	"chat-room-api/internal/core/domain"
)

type (
	UseCase struct {
		stockRepo         StockRepository
		messageRepository MessageRepository
		publisher         Publisher
	}

	StockRepository interface {
		GetPrice(context.Context, string) (float64, error)
	}

	MessageRepository interface {
		Save(context.Context, domain.Message) error
	}

	Publisher interface {
		Publish(context.Context, domain.Message, string) error
	}
)

func NewUseCase(
	sr StockRepository,
	mr MessageRepository,
	p Publisher,
) *UseCase {
	return &UseCase{
		stockRepo:         sr,
		messageRepository: mr,
		publisher:         p,
	}
}

func (uc *UseCase) Bot(ctx context.Context, msg domain.Message) error {
	var (
		message   = domain.NewMessage(msg.Room, domain.BotUsername, "", time.Now().UTC())
		stockName = domain.GetStockName(msg.Text)
	)
	price, err := uc.stockRepo.GetPrice(ctx, stockName)
	if err != nil {
		message.Text = fmt.Sprintf("%s: %s", stockName, err.Error())
	} else {
		message.Text = domain.CreateBotMessage(msg.Text, price)
		if err = uc.messageRepository.Save(ctx, message); err != nil {
			message.Text = fmt.Sprintf("%s: %s", stockName, err.Error())
		}
	}
	return uc.publisher.Publish(ctx, message, domain.BotMessageType)
}
