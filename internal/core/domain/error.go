package domain

import (
	"errors"
	"fmt"
)

var (
	// Password
	ErrHashingPassword = errors.New("error hashing password")
	ErrInvalidPassword = errors.New("invalid password")

	// DB
	ErrSavingMessage   = errors.New("error saving message")
	ErrGettingMessages = errors.New("error getting messages")
	ErrSavingUser      = errors.New("error saving user")
	ErrGettingUser     = errors.New("error getting user")

	// Publish
	ErrPublishingMessage = errors.New("error publishing message")

	// User
	ErrCreatingAccount = errors.New("error creating account")
	ErrLogin           = errors.New("error login user")
	ErrInvalidUser     = errors.New("invalid user")
	ErrGeneratingToken = errors.New("error generating token")

	// Chat
	ErrConnectingChat = errors.New("error connecting chat")
	ErrChatting       = errors.New("error chatting")
	ErrUnknownCommand = errors.New("error unknown command")

	// Stock
	ErrGettingStockPrice  = errors.New("error getting stock price")
	ErrStockPriceNotFound = errors.New("stock price not found")
)

func WrapError(domainErr error, err error) error {
	return fmt.Errorf("%w. wrapped err=%w", domainErr, err)
}
