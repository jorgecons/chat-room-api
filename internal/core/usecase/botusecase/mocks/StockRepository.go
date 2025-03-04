// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// StockRepository is an autogenerated mock type for the StockRepository type
type StockRepository struct {
	mock.Mock
}

type StockRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *StockRepository) EXPECT() *StockRepository_Expecter {
	return &StockRepository_Expecter{mock: &_m.Mock}
}

// GetPrice provides a mock function with given fields: _a0, _a1
func (_m *StockRepository) GetPrice(_a0 context.Context, _a1 string) (float64, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetPrice")
	}

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (float64, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) float64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StockRepository_GetPrice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPrice'
type StockRepository_GetPrice_Call struct {
	*mock.Call
}

// GetPrice is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *StockRepository_Expecter) GetPrice(_a0 interface{}, _a1 interface{}) *StockRepository_GetPrice_Call {
	return &StockRepository_GetPrice_Call{Call: _e.mock.On("GetPrice", _a0, _a1)}
}

func (_c *StockRepository_GetPrice_Call) Run(run func(_a0 context.Context, _a1 string)) *StockRepository_GetPrice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *StockRepository_GetPrice_Call) Return(_a0 float64, _a1 error) *StockRepository_GetPrice_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StockRepository_GetPrice_Call) RunAndReturn(run func(context.Context, string) (float64, error)) *StockRepository_GetPrice_Call {
	_c.Call.Return(run)
	return _c
}

// NewStockRepository creates a new instance of StockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStockRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *StockRepository {
	mock := &StockRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
