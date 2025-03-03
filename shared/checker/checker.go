package checker

import (
	"errors"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	vi   *validator.Validate
	lock sync.Mutex
)

// getValidator validator singleton
func getValidator() *validator.Validate {
	if vi == nil {
		lock.Lock()
		defer lock.Unlock()
		if vi == nil {
			vi = validator.New()
		}
	}

	return vi
}

// Check validate request by using validation tags
func Check(req interface{}) error {
	if err := getValidator().Struct(req); err != nil {
		var fieldErrList validator.ValidationErrors
		if ok := errors.As(err, &fieldErrList); !ok {
			return err
		}

		var valuesError string
		for i, v := range fieldErrList {
			if i == 0 {
				valuesError = v.Field()
			} else {
				valuesError = valuesError + "," + v.Field()
			}
		}

		return errors.New(valuesError)
	}

	return nil
}
