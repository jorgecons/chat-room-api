package connectchatusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"chat-room-api/internal/core/domain"
	"chat-room-api/internal/core/usecase/connectchatusecase"
	"chat-room-api/internal/core/usecase/connectchatusecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type scenario struct {
	t           *testing.T
	messageRepo *mocks.GetLastMessages
	room        string
	context     context.Context
	errorGot    error
	messages    []domain.Message
}

func TestUseCase_ConnectChat(t *testing.T) {
	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"connect chat - success": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidRoom("1")
			s.andGetLastMessagesSuccess()
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsOK()
		}},
		"connect chat - error - get last error": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidRoom("1")
			s.andGetLastMessagesWithError()
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
	return &scenario{
		t:           t,
		messageRepo: mocks.NewGetLastMessages(t),
		context:     context.Background(),
	}
}

func (s *scenario) givenAValidRoom(room string) {
	s.room = room
}

func (s *scenario) andGetLastMessagesSuccess() {
	msgs := []domain.Message{
		{
			Room:     s.room,
			Username: "some 1",
			Text:     "some 1",
			Date:     time.Date(2025, time.February, 1, 11, 0, 0, 0, time.UTC),
		},
		{
			Room:     s.room,
			Username: "some 2",
			Text:     "some 2",
			Date:     time.Date(2025, time.February, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			Room:     s.room,
			Username: "some 3",
			Text:     "some 3",
			Date:     time.Date(2025, time.February, 1, 13, 0, 0, 0, time.UTC),
		},
	}
	s.messageRepo.EXPECT().
		GetLastMessages(mock.Anything, s.room).
		Return(msgs, nil).
		Once()
}

func (s *scenario) andGetLastMessagesWithError() {
	err := errors.New("some_error")
	s.messageRepo.EXPECT().
		GetLastMessages(mock.Anything, s.room).
		Return(nil, err).
		Once()
}

func (s *scenario) whenTheUseCaseIsExecuted() {
	uc := connectchatusecase.NewUseCase(s.messageRepo)
	s.messages, s.errorGot = uc.ConnectChat(s.context, s.room)
}

func (s *scenario) thenTheExecutionIsOK() {
	s.t.Helper()
	assert.Nil(s.t, s.errorGot)
	assert.NotNil(s.t, s.messages)
	assert.Equal(s.t, "some 3", s.messages[0].Username)
}

func (s *scenario) thenTheExecutionIsNotOK() {
	s.t.Helper()
	assert.NotNil(s.t, s.errorGot)
	assert.ErrorIs(s.t, s.errorGot, domain.ErrConnectingChat)
}
