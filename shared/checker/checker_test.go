package checker

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type (
	scenario struct {
		testStruct interface{}
		validator  *validator.Validate
		error      error
	}
	request struct {
		Foo string `json:"foo" validate:"required,uuid4"`
		Var string `json:"var" validate:"required,uuid4"`
	}
)

func TestChecker(t *testing.T) {
	var (
		rightRequest = &request{
			Foo: "f3e35fb3-dd99-4269-8f50-b98e1c79fd69",
			Var: "f3e35fb3-dd99-4269-8f50-b98e1c79fd69",
		}
		wrongRequest = &request{
			Foo: "f3e35fb3",
		}
		emptyRequest          = &request{}
		nilRequest   *request = nil
	)

	testCases := map[string]func(t *testing.T){
		"right struct": func(t *testing.T) {
			s := startScenario()
			s.givenStruct(rightRequest)
			s.whenCheckerExecuted(func() {
				s.error = Check(s.testStruct)
			})
			s.thenTheExecutionIsOK(t)
		},
		"wrong struct": func(t *testing.T) {
			s := startScenario()
			s.givenStruct(wrongRequest)
			s.whenCheckerExecuted(func() {
				s.error = Check(s.testStruct)
			})
			s.thenTheExecutionIsNotOK(t)
		},
		"empty struct": func(t *testing.T) {
			s := startScenario()
			s.givenStruct(emptyRequest)
			s.whenCheckerExecuted(func() {
				s.error = Check(s.testStruct)
			})
			s.thenTheExecutionIsNotOK(t)
		},
		"nil struct": func(t *testing.T) {
			s := startScenario()
			s.givenStruct(nilRequest)
			s.whenCheckerExecuted(func() {
				s.error = Check(s.testStruct)
			})
			s.thenTheExecutionIsNotOK(t)
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc)
	}
}

func startScenario() *scenario {
	return &scenario{}
}

func (s *scenario) givenStruct(t interface{}) {
	s.testStruct = t
}

func (s *scenario) whenCheckerExecuted(function func()) {
	function()
}

func (s *scenario) thenTheExecutionIsOK(t *testing.T) {
	assert.Nil(t, s.error)
}

func (s *scenario) thenTheExecutionIsNotOK(t *testing.T) {
	assert.NotNil(t, s.error)
}
