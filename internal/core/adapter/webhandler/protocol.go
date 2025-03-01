package webhandler

import (
	"errors"

	"chat-room-api/internal/core/domain"
	"chat-room-api/shared/checker"
)

const InvalidUserCode = "invalid_user"

var ErrFailedValidations = errors.New("failed validations over the message")

type (
	User struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	Error struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}

	Token struct {
		Token string `json:"token"`
	}
)

func (e Error) Error() string {
	return e.Message
}

func BuildUser(user User) domain.User {
	return domain.User{
		Username: user.Username,
		Password: user.Password,
	}
}

func BuildErrorResponse(code string, err error) Error {
	return Error{
		Message: err.Error(),
		Code:    code,
	}
}

func ValidateAccountRequest(user User) error {
	if err := checker.Check(user); err != nil {
		msg := domain.WrapError(ErrFailedValidations, err)
		return BuildErrorResponse(InvalidUserCode, msg)
	}
	return nil
}
