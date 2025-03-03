package loginusecase_test

import (
	"context"
	"errors"
	"testing"

	"chat-room-api/internal/core/domain"
	"chat-room-api/internal/core/usecase/loginusecase"
	"chat-room-api/internal/core/usecase/loginusecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type scenario struct {
	t        *testing.T
	userRepo *mocks.UserRepository
	user     domain.User
	secret   []byte
	context  context.Context
	errorGot error
	token    string
}

func TestUseCase_Login(t *testing.T) {
	user := domain.User{
		Username: "some",
		Password: "some",
	}
	secret := []byte("some_secret")

	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"login - success": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidUser(user)
			s.andAValidSecret(secret)
			s.andGetUserSuccess(user.Password)
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsOK()
		}},
		"login - error - invalid password": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidUser(user)
			s.andAValidSecret(secret)
			s.andGetUserSuccess("other_password")
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsNotOK(domain.ErrInvalidPassword)
		}},
		"login - error - get user": {func(t *testing.T) {
			s := startScenario(t)
			s.givenAValidUser(user)
			s.andAValidSecret(secret)
			s.andGetUserWithError()
			s.whenTheUseCaseIsExecuted()
			s.thenTheExecutionIsNotOK(domain.ErrLogin)
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

func (s *scenario) andAValidSecret(secret []byte) {
	s.secret = secret
}

func (s *scenario) andGetUserSuccess(pass string) {
	p, _ := domain.HashPassword(pass)
	res := domain.User{
		Username: s.user.Username,
		Password: p,
	}
	s.userRepo.EXPECT().
		Get(mock.Anything, s.user.Username).
		Return(res, nil).
		Once()
}

func (s *scenario) andGetUserWithError() {
	s.userRepo.EXPECT().
		Get(mock.Anything, mock.Anything).
		Return(domain.User{}, errors.New("some_error")).
		Once()
}

func (s *scenario) whenTheUseCaseIsExecuted() {
	uc := loginusecase.NewUseCase(s.userRepo, s.secret)
	s.token, s.errorGot = uc.Login(s.context, s.user)
}

func (s *scenario) thenTheExecutionIsOK() {
	s.t.Helper()
	assert.Nil(s.t, s.errorGot)
	assert.NotNil(s.t, s.token)
}

func (s *scenario) thenTheExecutionIsNotOK(err error) {
	s.t.Helper()
	assert.NotNil(s.t, s.errorGot)
	assert.ErrorIs(s.t, s.errorGot, err)
}
