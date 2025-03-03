package webhandler

import (
	"errors"

	"chat-room-api/internal/core/domain"
	"chat-room-api/shared/checker"
)

const (
	InvalidCredentialsErrorCode = "invalid_credentials_error"
	CreateUserErrorCode         = "create_user_error"
	BadRequestErrorCode         = "bad_request"
)

var (
	ErrFailedValidations = errors.New("failed validations over the message")
	ErrLoginError        = errors.New("error login")
)

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

func BuildErrorResponse(err error, code string) Error {
	return Error{
		Message: err.Error(),
		Code:    code,
	}
}

func ValidateAccountRequest(user User) error {
	if err := checker.Check(user); err != nil {
		return domain.WrapError(ErrFailedValidations, err)
	}
	return nil
}
