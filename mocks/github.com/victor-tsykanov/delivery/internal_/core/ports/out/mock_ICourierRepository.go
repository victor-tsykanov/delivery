// Code generated by mockery v2.52.3. DO NOT EDIT.

package out

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	courier "github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
)

// MockICourierRepository is an autogenerated mock type for the ICourierRepository type
type MockICourierRepository struct {
	mock.Mock
}

type MockICourierRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockICourierRepository) EXPECT() *MockICourierRepository_Expecter {
	return &MockICourierRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *MockICourierRepository) Create(ctx context.Context, _a1 *courier.Courier) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *courier.Courier) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockICourierRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockICourierRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *courier.Courier
func (_e *MockICourierRepository_Expecter) Create(ctx interface{}, _a1 interface{}) *MockICourierRepository_Create_Call {
	return &MockICourierRepository_Create_Call{Call: _e.mock.On("Create", ctx, _a1)}
}

func (_c *MockICourierRepository_Create_Call) Run(run func(ctx context.Context, _a1 *courier.Courier)) *MockICourierRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*courier.Courier))
	})
	return _c
}

func (_c *MockICourierRepository_Create_Call) Return(_a0 error) *MockICourierRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockICourierRepository_Create_Call) RunAndReturn(run func(context.Context, *courier.Courier) error) *MockICourierRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// FindFree provides a mock function with given fields: ctx
func (_m *MockICourierRepository) FindFree(ctx context.Context) ([]*courier.Courier, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FindFree")
	}

	var r0 []*courier.Courier
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*courier.Courier, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*courier.Courier); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*courier.Courier)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockICourierRepository_FindFree_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindFree'
type MockICourierRepository_FindFree_Call struct {
	*mock.Call
}

// FindFree is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockICourierRepository_Expecter) FindFree(ctx interface{}) *MockICourierRepository_FindFree_Call {
	return &MockICourierRepository_FindFree_Call{Call: _e.mock.On("FindFree", ctx)}
}

func (_c *MockICourierRepository_FindFree_Call) Run(run func(ctx context.Context)) *MockICourierRepository_FindFree_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockICourierRepository_FindFree_Call) Return(_a0 []*courier.Courier, _a1 error) *MockICourierRepository_FindFree_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockICourierRepository_FindFree_Call) RunAndReturn(run func(context.Context) ([]*courier.Courier, error)) *MockICourierRepository_FindFree_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *MockICourierRepository) Get(ctx context.Context, id courier.ID) (*courier.Courier, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *courier.Courier
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, courier.ID) (*courier.Courier, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, courier.ID) *courier.Courier); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*courier.Courier)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, courier.ID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockICourierRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockICourierRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id courier.ID
func (_e *MockICourierRepository_Expecter) Get(ctx interface{}, id interface{}) *MockICourierRepository_Get_Call {
	return &MockICourierRepository_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *MockICourierRepository_Get_Call) Run(run func(ctx context.Context, id courier.ID)) *MockICourierRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(courier.ID))
	})
	return _c
}

func (_c *MockICourierRepository_Get_Call) Return(_a0 *courier.Courier, _a1 error) *MockICourierRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockICourierRepository_Get_Call) RunAndReturn(run func(context.Context, courier.ID) (*courier.Courier, error)) *MockICourierRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, _a1
func (_m *MockICourierRepository) Update(ctx context.Context, _a1 *courier.Courier) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *courier.Courier) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockICourierRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockICourierRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *courier.Courier
func (_e *MockICourierRepository_Expecter) Update(ctx interface{}, _a1 interface{}) *MockICourierRepository_Update_Call {
	return &MockICourierRepository_Update_Call{Call: _e.mock.On("Update", ctx, _a1)}
}

func (_c *MockICourierRepository_Update_Call) Run(run func(ctx context.Context, _a1 *courier.Courier)) *MockICourierRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*courier.Courier))
	})
	return _c
}

func (_c *MockICourierRepository_Update_Call) Return(_a0 error) *MockICourierRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockICourierRepository_Update_Call) RunAndReturn(run func(context.Context, *courier.Courier) error) *MockICourierRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockICourierRepository creates a new instance of MockICourierRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockICourierRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockICourierRepository {
	mock := &MockICourierRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
