package createaccountusecase_test

import (
	"context"
	"errors"
	"testing"

	"chat-room-api/internal/core/domain"
	"chat-room-api/internal/core/usecase/createaccountusecase"
	"chat-room-api/internal/core/usecase/createaccountusecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type scenario struct {
	t        *testing.T
	userRepo *mocks.UserRepository
	user     domain.User
	context  context.Context
	errorGot error
}

func TestUseCase_CreateAccount(t *testing.T) {
	user := domain.User{
		Username: "some",
		Password: "some",
	}
	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"create account - success": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidUser(user)
			s.andSaveUserSuccess()
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsOK()
		}},
		"connect chat - error - save error": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidUser(user)
			s.andSaveUserWithError()
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
		t:        t,
		userRepo: mocks.NewUserRepository(t),
		context:  context.Background(),
	}
}

func (s *scenario) givenAValidUser(user domain.User) {
	s.user = user
}

func (s *scenario) andSaveUserSuccess() {
	s.userRepo.EXPECT().
		Save(mock.Anything, mock.Anything).
		Return(nil).
		Once()
}

func (s *scenario) andSaveUserWithError() {
	s.userRepo.EXPECT().
		Save(mock.Anything, mock.Anything).
		Return(errors.New("some_error")).
		Once()
}

func (s *scenario) whenTheUseCaseIsExecuted() {
	uc := createaccountusecase.NewUseCase(s.userRepo)
	s.errorGot = uc.CreateAccount(s.context, s.user)
}

func (s *scenario) thenTheExecutionIsOK() {
	s.t.Helper()
	assert.Nil(s.t, s.errorGot)
}

func (s *scenario) thenTheExecutionIsNotOK() {
	s.t.Helper()
	assert.NotNil(s.t, s.errorGot)
	assert.ErrorIs(s.t, s.errorGot, domain.ErrCreatingAccount)
}
