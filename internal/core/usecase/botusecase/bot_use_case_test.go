package botusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"chat-room-api/internal/core/domain"
	"chat-room-api/internal/core/usecase/botusecase"
	"chat-room-api/internal/core/usecase/botusecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type scenario struct {
	t           *testing.T
	stockRepo   *mocks.StockRepository
	messageRepo *mocks.MessageRepository
	publisher   *mocks.Publisher
	msg         domain.Message
	botMsg      domain.Message
	context     context.Context
	errorGot    error
}

func TestUseCase_Bot(t *testing.T) {
	msg := domain.Message{
		Room:     "1",
		Username: "some",
		Text:     "/stock=some",
		Date:     time.Date(2025, time.February, 1, 12, 0, 0, 0, time.UTC),
	}

	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"bot - success": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidMessage(msg)
			s.andGetPriceSuccess(100)
			s.andSaveMessageSuccess()
			s.andPublishMessageSuccess()
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsOK()
		}},
		"bot - error - price error": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidMessage(msg)
			s.andGetPriceError()
			s.andPublishMessageSuccess()
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsOK()
		}},
		"bot - error - save error": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidMessage(msg)
			s.andGetPriceSuccess(100)
			s.andSaveMessageWithError()
			s.andPublishMessageSuccess()
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsOK()
		}},
		"bot - error - publish error": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidMessage(msg)
			s.andGetPriceSuccess(100)
			s.andSaveMessageSuccess()
			s.andPublishMessageWithError()
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsNotOK()
		}},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc.run)
	}
}

func startScenario(t *testing.T) *scenario {
	domain.Now = func() time.Time { return time.Date(2025, time.February, 1, 12, 10, 0, 0, time.UTC) }
	return &scenario{
		t:           t,
		stockRepo:   mocks.NewStockRepository(t),
		messageRepo: mocks.NewMessageRepository(t),
		publisher:   mocks.NewPublisher(t),
		context:     context.Background(),
	}
}

func (s *scenario) givenAValidMessage(message domain.Message) {
	s.msg = message
}

func (s *scenario) andGetPriceSuccess(price float64) {
	msg := domain.NewBotMessage(s.msg.Room)
	msg.Text = domain.CreateBotText(domain.GetStockName(s.msg.Text), price)
	s.botMsg = msg
	s.stockRepo.EXPECT().
		GetPrice(mock.Anything, domain.GetStockName(s.msg.Text)).
		Return(price, nil).
		Once()
}

func (s *scenario) andGetPriceError() {
	err := errors.New("some_error")
	msg := domain.NewBotMessage(s.msg.Room)
	msg.Text = domain.CreateBotErrorText(domain.GetStockName(s.msg.Text), err)
	s.botMsg = msg
	s.stockRepo.EXPECT().
		GetPrice(mock.Anything, domain.GetStockName(s.msg.Text)).
		Return(0, err).
		Once()
}

func (s *scenario) andSaveMessageSuccess() {
	s.messageRepo.EXPECT().
		Save(mock.Anything, s.botMsg).
		Return(nil).
		Once()
}

func (s *scenario) andSaveMessageWithError() {
	err := errors.New("some_error")
	s.messageRepo.EXPECT().
		Save(mock.Anything, s.botMsg).
		Return(err).
		Once()

	s.botMsg.Text = domain.CreateBotErrorText(domain.GetStockName(s.msg.Text), err)
}

func (s *scenario) andPublishMessageSuccess() {
	s.publisher.EXPECT().
		Publish(mock.Anything, s.botMsg, domain.BotMessageType).
		Return(nil).
		Once()
}

func (s *scenario) andPublishMessageWithError() {
	s.publisher.EXPECT().
		Publish(mock.Anything, s.botMsg, domain.BotMessageType).
		Return(domain.ErrPublishingMessage).
		Once()
}

func (s *scenario) whenTheUseCaseIsExecuted() {
	uc := botusecase.NewUseCase(s.stockRepo, s.messageRepo, s.publisher)
	s.errorGot = uc.Bot(s.context, s.msg)
}

func (s *scenario) thenTheExecutionIsOK() {
	s.t.Helper()
	assert.Nil(s.t, s.errorGot)
}

func (s *scenario) thenTheExecutionIsNotOK() {
	s.t.Helper()
	assert.NotNil(s.t, s.errorGot)
	assert.ErrorIs(s.t, s.errorGot, domain.ErrPublishingMessage)
}
