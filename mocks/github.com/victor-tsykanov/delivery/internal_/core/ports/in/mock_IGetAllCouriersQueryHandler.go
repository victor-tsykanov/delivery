// Code generated by mockery v2.52.3. DO NOT EDIT.

package in

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	in "github.com/victor-tsykanov/delivery/internal/core/ports/in"
)

// MockIGetAllCouriersQueryHandler is an autogenerated mock type for the IGetAllCouriersQueryHandler type
type MockIGetAllCouriersQueryHandler struct {
	mock.Mock
}

type MockIGetAllCouriersQueryHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIGetAllCouriersQueryHandler) EXPECT() *MockIGetAllCouriersQueryHandler_Expecter {
	return &MockIGetAllCouriersQueryHandler_Expecter{mock: &_m.Mock}
}

// Handle provides a mock function with given fields: ctx
func (_m *MockIGetAllCouriersQueryHandler) Handle(ctx context.Context) ([]*in.Courier, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Handle")
	}

	var r0 []*in.Courier
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*in.Courier, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*in.Courier); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*in.Courier)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIGetAllCouriersQueryHandler_Handle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Handle'
type MockIGetAllCouriersQueryHandler_Handle_Call struct {
	*mock.Call
}

// Handle is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockIGetAllCouriersQueryHandler_Expecter) Handle(ctx interface{}) *MockIGetAllCouriersQueryHandler_Handle_Call {
	return &MockIGetAllCouriersQueryHandler_Handle_Call{Call: _e.mock.On("Handle", ctx)}
}

func (_c *MockIGetAllCouriersQueryHandler_Handle_Call) Run(run func(ctx context.Context)) *MockIGetAllCouriersQueryHandler_Handle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockIGetAllCouriersQueryHandler_Handle_Call) Return(_a0 []*in.Courier, _a1 error) *MockIGetAllCouriersQueryHandler_Handle_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIGetAllCouriersQueryHandler_Handle_Call) RunAndReturn(run func(context.Context) ([]*in.Courier, error)) *MockIGetAllCouriersQueryHandler_Handle_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIGetAllCouriersQueryHandler creates a new instance of MockIGetAllCouriersQueryHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIGetAllCouriersQueryHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIGetAllCouriersQueryHandler {
	mock := &MockIGetAllCouriersQueryHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
