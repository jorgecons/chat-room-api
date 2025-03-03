package createaccountusecase

import (
	"context"

	"chat-room-api/internal/core/domain"

	"github.com/sirupsen/logrus"
)

type (
	UseCase struct {
		userStorage UserRepository
	}

	UserRepository interface {
		Save(context.Context, domain.User) error
	}
)

func NewUseCase(
	userStorage UserRepository,
) *UseCase {
	return &UseCase{
		userStorage: userStorage,
	}
}

func (uc *UseCase) CreateAccount(ctx context.Context, user domain.User) error {
	hashedPassword, err := domain.HashPassword(user.Password)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField("user", user).Error("Error creating user")
		return err
	}
	user.Password = hashedPassword
	if err = uc.userStorage.Save(ctx, user); err != nil {
		err = domain.WrapError(domain.ErrCreatingAccount, err)
		logrus.WithContext(ctx).WithError(err).WithField("user", user.Username).Error("Error creating user")
		return err
	}
	return nil
}
