// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	domain "chat-room-api/internal/core/domain"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Publisher is an autogenerated mock type for the Publisher type
type Publisher struct {
	mock.Mock
}

type Publisher_Expecter struct {
	mock *mock.Mock
}

func (_m *Publisher) EXPECT() *Publisher_Expecter {
	return &Publisher_Expecter{mock: &_m.Mock}
}

// Publish provides a mock function with given fields: _a0, _a1, _a2
func (_m *Publisher) Publish(_a0 context.Context, _a1 domain.Message, _a2 string) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for Publish")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Message, string) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Publisher_Publish_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Publish'
type Publisher_Publish_Call struct {
	*mock.Call
}

// Publish is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 domain.Message
//   - _a2 string
func (_e *Publisher_Expecter) Publish(_a0 interface{}, _a1 interface{}, _a2 interface{}) *Publisher_Publish_Call {
	return &Publisher_Publish_Call{Call: _e.mock.On("Publish", _a0, _a1, _a2)}
}

func (_c *Publisher_Publish_Call) Run(run func(_a0 context.Context, _a1 domain.Message, _a2 string)) *Publisher_Publish_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Message), args[2].(string))
	})
	return _c
}

func (_c *Publisher_Publish_Call) Return(_a0 error) *Publisher_Publish_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Publisher_Publish_Call) RunAndReturn(run func(context.Context, domain.Message, string) error) *Publisher_Publish_Call {
	_c.Call.Return(run)
	return _c
}

// NewPublisher creates a new instance of Publisher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPublisher(t interface {
	mock.TestingT
	Cleanup(func())
}) *Publisher {
	mock := &Publisher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
