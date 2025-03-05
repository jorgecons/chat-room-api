package message_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"bot/internal/core/domain"
	"bot/internal/core/repository/message"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type scenario struct {
	t           *testing.T
	dbConn      pgxmock.PgxConnIface
	message     domain.Message
	room        string
	context     context.Context
	errorGot    error
	messagesGot []domain.Message
}

const (
	sqlSaveQuery = `INSERT INTO messages \(room, username, text, date\) VALUES \(\$1, \$2, \$3, \$4\)`
)

func TestMessage_Save(t *testing.T) {
	domain.Now = func() time.Time { return time.Date(2025, time.February, 1, 12, 10, 0, 0, time.UTC) }
	m := domain.Message{
		Room:     "1",
		Username: "some",
		Text:     "some",
		Date:     domain.Now(),
	}
	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"save - success": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidMessage(m)
			s.andSaveMessageSuccess()
			s.whenTheSaveIsExecuted()
			s.thenTheSaveExecutionIsOK()
		}},
		"save - error": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidMessage(m)
			s.andSaveMessageWithError()
			s.whenTheSaveIsExecuted()
			s.thenTheExecutionIsNotOK(domain.ErrSavingMessage)
		}},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc.run)
	}
}

func startScenario(t *testing.T) *scenario {
	m, err := pgxmock.NewConn()
	require.NoError(t, err)
	return &scenario{
		t:       t,
		dbConn:  m,
		context: context.Background(),
	}
}

func (s *scenario) givenAValidMessage(msg domain.Message) {
	s.message = msg
}

func (s *scenario) andSaveMessageSuccess() {
	s.dbConn.ExpectExec(sqlSaveQuery).
		WithArgs(s.message.Room, s.message.Username, s.message.Text, s.message.Date).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
}

func (s *scenario) andSaveMessageWithError() {
	s.dbConn.ExpectExec(sqlSaveQuery).
		WithArgs(s.message.Room, s.message.Username, s.message.Text, s.message.Date).
		WillReturnError(errors.New("database error"))
}

func (s *scenario) whenTheSaveIsExecuted() {
	r := message.NewMessageStorage(s.dbConn)
	s.errorGot = r.Save(s.context, s.message)
}

func (s *scenario) thenTheSaveExecutionIsOK() {
	s.t.Helper()
	assert.Nil(s.t, s.errorGot)
}

func (s *scenario) thenTheGetLastMessagesExecutionIsOK() {
	s.t.Helper()
	assert.Nil(s.t, s.errorGot)
	assert.NotNil(s.t, s.messagesGot)
	assert.Equal(s.t, s.room, s.messagesGot[0].Room)
}

func (s *scenario) thenTheExecutionIsNotOK(err error) {
	s.t.Helper()
	assert.NotNil(s.t, s.errorGot)
	assert.ErrorIs(s.t, s.errorGot, err)
}
