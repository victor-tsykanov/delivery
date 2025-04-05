// Code generated by mockery v2.52.3. DO NOT EDIT.

package eventdispatcher

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockIEventDispatcher is an autogenerated mock type for the IEventDispatcher type
type MockIEventDispatcher struct {
	mock.Mock
}

type MockIEventDispatcher_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIEventDispatcher) EXPECT() *MockIEventDispatcher_Expecter {
	return &MockIEventDispatcher_Expecter{mock: &_m.Mock}
}

// Dispatch provides a mock function with given fields: ctx, event
func (_m *MockIEventDispatcher) Dispatch(ctx context.Context, event interface{}) error {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for Dispatch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIEventDispatcher_Dispatch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Dispatch'
type MockIEventDispatcher_Dispatch_Call struct {
	*mock.Call
}

// Dispatch is a helper method to define mock.On call
//   - ctx context.Context
//   - event interface{}
func (_e *MockIEventDispatcher_Expecter) Dispatch(ctx interface{}, event interface{}) *MockIEventDispatcher_Dispatch_Call {
	return &MockIEventDispatcher_Dispatch_Call{Call: _e.mock.On("Dispatch", ctx, event)}
}

func (_c *MockIEventDispatcher_Dispatch_Call) Run(run func(ctx context.Context, event interface{})) *MockIEventDispatcher_Dispatch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(interface{}))
	})
	return _c
}

func (_c *MockIEventDispatcher_Dispatch_Call) Return(_a0 error) *MockIEventDispatcher_Dispatch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIEventDispatcher_Dispatch_Call) RunAndReturn(run func(context.Context, interface{}) error) *MockIEventDispatcher_Dispatch_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIEventDispatcher creates a new instance of MockIEventDispatcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIEventDispatcher(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIEventDispatcher {
	mock := &MockIEventDispatcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
