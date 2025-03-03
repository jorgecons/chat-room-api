package chatusecase_test

import (
	"chat-room-api/internal/core/domain"
	"chat-room-api/internal/core/usecase/chatusecase"
	"chat-room-api/internal/core/usecase/chatusecase/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type scenario struct {
	t           *testing.T
	messageRepo *mocks.MessageRepository
	publisher   *mocks.StockPublisher
	msg         domain.Message
	context     context.Context
	errorGot    error
}

func TestUseCase_Chat(t *testing.T) {
	msg := domain.Message{
		Room:     "1",
		Username: "some",
		Text:     "some",
		Date:     time.Date(2025, time.February, 1, 12, 0, 0, 0, time.UTC),
	}

	commandMsg := domain.Message{
		Room:     "1",
		Username: "some",
		Text:     "/stock=some",
		Date:     time.Date(2025, time.February, 1, 12, 0, 0, 0, time.UTC),
	}

	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"chat - success": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidMessage(msg)
			s.andSaveMessageSuccess()
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsOK()
		}},
		"chat - command - success": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidMessage(commandMsg)
			s.andPublishMessageSuccess()
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsOK()
		}},
		"chat - error - save error": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidMessage(msg)
			s.andSaveMessageWithError()
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsNotOK()
		}},
		"chat - command - error - publish error": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidMessage(commandMsg)
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
		messageRepo: mocks.NewMessageRepository(t),
		publisher:   mocks.NewStockPublisher(t),
		context:     context.Background(),
	}
}

func (s *scenario) givenAValidMessage(message domain.Message) {
	s.msg = message
}

func (s *scenario) andSaveMessageSuccess() {
	s.messageRepo.EXPECT().
		Save(mock.Anything, s.msg).
		Return(nil).
		Once()
}

func (s *scenario) andSaveMessageWithError() {
	err := errors.New("some_error")
	s.messageRepo.EXPECT().
		Save(mock.Anything, s.msg).
		Return(err).
		Once()
}

func (s *scenario) andPublishMessageSuccess() {
	s.publisher.EXPECT().
		Publish(mock.Anything, s.msg, domain.StockMessageType).
		Return(nil).
		Once()
}

func (s *scenario) andPublishMessageWithError() {
	s.publisher.EXPECT().
		Publish(mock.Anything, s.msg, domain.StockMessageType).
		Return(domain.ErrPublishingMessage).
		Once()
}

func (s *scenario) whenTheUseCaseIsExecuted() {
	uc := chatusecase.NewUseCase(s.messageRepo, s.publisher)
	s.errorGot = uc.Chat(s.context, s.msg)
}

func (s *scenario) thenTheExecutionIsOK() {
	s.t.Helper()
	assert.Nil(s.t, s.errorGot)
}

func (s *scenario) thenTheExecutionIsNotOK() {
	s.t.Helper()
	assert.NotNil(s.t, s.errorGot)
	assert.ErrorIs(s.t, s.errorGot, domain.ErrChatting)
}
