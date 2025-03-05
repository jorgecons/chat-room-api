package eventhandler

import (
	"context"
	"encoding/json"

	"bot/internal/core/domain"

	"github.com/sirupsen/logrus"
)

type (
	BotFeature interface {
		Bot(context.Context, domain.Message) error
	}
	bot struct {
		feature BotFeature
	}
)

func NewBotHandler(f BotFeature) func(context.Context, []byte) error {
	return bot{feature: f}.Handle
}

func (b bot) Handle(ctx context.Context, ev []byte) error {
	req := Message{}
	if err := json.Unmarshal(ev, &req); err != nil {
		logrus.WithContext(ctx).WithError(err).Error("Error unmarshalling event")
		return err
	}
	msg := BuildMessage(req)
	return b.feature.Bot(ctx, msg)
}
