package user_test

import (
	"context"
	"errors"
	"testing"

	"chat-room-api/internal/core/domain"
	"chat-room-api/internal/core/repository/user"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type scenario struct {
	t        *testing.T
	dbConn   pgxmock.PgxConnIface
	user     domain.User
	username string
	context  context.Context
	errorGot error
	userGot  domain.User
}

const (
	sqlSaveQuery = `INSERT INTO users \(username, password_hash\) VALUES \(\$1, \$2\)`
	sqlGetQuery  = `SELECT username, password_hash FROM users WHERE username=\$1`
)

func TestUser_Save(t *testing.T) {
	u := domain.User{
		Username: "some",
		Password: "some",
	}

	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"save - success": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidUser(u)
			s.andSaveUserSuccess()
			s.whenTheSaveIsExecuted()
			s.thenTheSaveExecutionIsOK()
		}},
		"save - error": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidUser(u)
			s.andSaveUserWithError()
			s.whenTheSaveIsExecuted()
			s.thenTheExecutionIsNotOK(domain.ErrSavingUser)
		}},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc.run)
	}
}

func TestUser_Get(t *testing.T) {
	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"get - success": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidUsername("some")
			s.andGetUserSuccess()
			s.whenTheGetIsExecuted()
			s.thenTheGetExecutionIsOK()
		}},
		"get - error - zero rows": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidUsername("some")
			s.andGetUserZeroRows()
			s.whenTheGetIsExecuted()
			s.thenTheExecutionIsNotOK(domain.ErrInvalidUser)
		}},
		"get - error": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidUsername("some")
			s.andGetUserWithError()
			s.whenTheGetIsExecuted()
			s.thenTheExecutionIsNotOK(domain.ErrGettingUser)
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

func (s *scenario) givenAValidUser(user domain.User) {
	s.user = user
}

func (s *scenario) givenAValidUsername(username string) {
	s.username = username
}

func (s *scenario) andSaveUserSuccess() {
	s.dbConn.ExpectExec(sqlSaveQuery).
		WithArgs(s.user.Username, s.user.Password).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
}

func (s *scenario) andSaveUserWithError() {
	s.dbConn.ExpectExec(sqlSaveQuery).
		WithArgs(s.user.Username, s.user.Password).
		WillReturnError(errors.New("database error"))
}

func (s *scenario) andGetUserSuccess() {
	s.dbConn.ExpectQuery(sqlGetQuery).
		WithArgs(s.username).
		WillReturnRows(pgxmock.NewRows([]string{"username", "password_hash"}).AddRow(s.username, "some"))
}

func (s *scenario) andGetUserZeroRows() {
	s.dbConn.ExpectQuery(sqlGetQuery).
		WithArgs(s.username).
		WillReturnError(pgx.ErrNoRows)
}

func (s *scenario) andGetUserWithError() {
	s.dbConn.ExpectQuery(sqlGetQuery).
		WithArgs(s.username).
		WillReturnError(errors.New("database error"))
}

func (s *scenario) whenTheSaveIsExecuted() {
	r := user.NewUserStorage(s.dbConn)
	s.errorGot = r.Save(s.context, s.user)
}

func (s *scenario) whenTheGetIsExecuted() {
	r := user.NewUserStorage(s.dbConn)
	s.userGot, s.errorGot = r.Get(s.context, s.username)
}

func (s *scenario) thenTheSaveExecutionIsOK() {
	s.t.Helper()
	assert.Nil(s.t, s.errorGot)
}

func (s *scenario) thenTheGetExecutionIsOK() {
	s.t.Helper()
	assert.Nil(s.t, s.errorGot)
	assert.NotNil(s.t, s.userGot)
	assert.Equal(s.t, s.username, s.userGot.Username)
}

func (s *scenario) thenTheExecutionIsNotOK(err error) {
	s.t.Helper()
	assert.NotNil(s.t, s.errorGot)
	assert.ErrorIs(s.t, s.errorGot, err)
}
