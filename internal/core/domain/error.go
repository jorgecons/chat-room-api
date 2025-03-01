package domain

import (
	"errors"
	"fmt"
)

var (
	// User
	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
	ErrInvalidCardExpiration = errors.New("invalid card expiration format MM/YY")
	ErrInvalidTotalAmount    = errors.New("invalid total amount decimal")

	// Password
	ErrHashingPassword = errors.New("error hashing password")
	ErrInvalidPassword = errors.New("invalid password")

	// RefundCall
	ErrProcessRefundCall          = errors.New("error processing refund call")
	ErrAlreadyProcessedRefundCall = errors.New("error already processed refund call")

	// DB
	ErrSavingMessage = errors.New("error saving message")
	ErrSavingUser    = errors.New("error saving user")
	ErrGettingUser   = errors.New("error getting user")

	// User
	ErrCreatingAccount = errors.New("error creating account")
	ErrLogin           = errors.New("error login user")
	ErrInvalidUser     = errors.New("invalid user")
	ErrGeneratingToken = errors.New("error generating token")
)

func WrapError(domainErr error, err error) error {
	return fmt.Errorf("%w. wrapped err=%w", domainErr, err)
}
