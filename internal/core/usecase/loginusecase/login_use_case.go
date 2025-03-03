package loginusecase

import (
	"context"

	"chat-room-api/internal/core/domain"

	"github.com/sirupsen/logrus"
)

const ZeroToken = ""

type (
	UseCase struct {
		userStorage UserRepository
		jwtToken    []byte
	}

	UserRepository interface {
		Get(context.Context, string) (domain.User, error)
	}
)

func NewUseCase(
	userStorage UserRepository,
	jwtSecret []byte,
) *UseCase {
	return &UseCase{
		userStorage: userStorage,
		jwtToken:    jwtSecret,
	}
}

func (uc *UseCase) Login(ctx context.Context, user domain.User) (string, error) {
	savedUser, err := uc.userStorage.Get(ctx, user.Username)
	if err != nil {
		err = domain.WrapError(domain.ErrLogin, err)
		logrus.WithContext(ctx).WithError(err).WithField("user", user.Username).Error("Error login user")
		return ZeroToken, err
	}
	err = domain.ValidatePassword(savedUser.Password, user.Password)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField("user", user.Username).Error("Error validating password")
		return ZeroToken, domain.ErrInvalidPassword
	}
	token, err := domain.GenerateToken(uc.jwtToken, user.Username)
	if err != nil {
		logrus.WithContext(ctx).WithField("user", user.Username).Error("Error generating token")
		return ZeroToken, domain.ErrInvalidPassword
	}
	return token, nil
}
