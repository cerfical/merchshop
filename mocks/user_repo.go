// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/cerfical/merchshop/internal/service/model"
	mock "github.com/stretchr/testify/mock"
)

// UserRepo is an autogenerated mock type for the UserRepo type
type UserRepo struct {
	mock.Mock
}

type UserRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *UserRepo) EXPECT() *UserRepo_Expecter {
	return &UserRepo_Expecter{mock: &_m.Mock}
}

// CreateUser provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *UserRepo) CreateUser(_a0 context.Context, _a1 model.Username, _a2 model.PasswordHash, _a3 model.NumCoins) (*model.User, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Username, model.PasswordHash, model.NumCoins) (*model.User, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.Username, model.PasswordHash, model.NumCoins) *model.User); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.Username, model.PasswordHash, model.NumCoins) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepo_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type UserRepo_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 model.Username
//   - _a2 model.PasswordHash
//   - _a3 model.NumCoins
func (_e *UserRepo_Expecter) CreateUser(_a0 interface{}, _a1 interface{}, _a2 interface{}, _a3 interface{}) *UserRepo_CreateUser_Call {
	return &UserRepo_CreateUser_Call{Call: _e.mock.On("CreateUser", _a0, _a1, _a2, _a3)}
}

func (_c *UserRepo_CreateUser_Call) Run(run func(_a0 context.Context, _a1 model.Username, _a2 model.PasswordHash, _a3 model.NumCoins)) *UserRepo_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.Username), args[2].(model.PasswordHash), args[3].(model.NumCoins))
	})
	return _c
}

func (_c *UserRepo_CreateUser_Call) Return(_a0 *model.User, _a1 error) *UserRepo_CreateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepo_CreateUser_Call) RunAndReturn(run func(context.Context, model.Username, model.PasswordHash, model.NumCoins) (*model.User, error)) *UserRepo_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// GetUser provides a mock function with given fields: _a0, _a1
func (_m *UserRepo) GetUser(_a0 context.Context, _a1 model.Username) (*model.User, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Username) (*model.User, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.Username) *model.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.Username) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepo_GetUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUser'
type UserRepo_GetUser_Call struct {
	*mock.Call
}

// GetUser is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 model.Username
func (_e *UserRepo_Expecter) GetUser(_a0 interface{}, _a1 interface{}) *UserRepo_GetUser_Call {
	return &UserRepo_GetUser_Call{Call: _e.mock.On("GetUser", _a0, _a1)}
}

func (_c *UserRepo_GetUser_Call) Run(run func(_a0 context.Context, _a1 model.Username)) *UserRepo_GetUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.Username))
	})
	return _c
}

func (_c *UserRepo_GetUser_Call) Return(_a0 *model.User, _a1 error) *UserRepo_GetUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepo_GetUser_Call) RunAndReturn(run func(context.Context, model.Username) (*model.User, error)) *UserRepo_GetUser_Call {
	_c.Call.Return(run)
	return _c
}

// PurchaseMerch provides a mock function with given fields: ctx, buyer, m
func (_m *UserRepo) PurchaseMerch(ctx context.Context, buyer model.UserID, m *model.MerchItem) error {
	ret := _m.Called(ctx, buyer, m)

	if len(ret) == 0 {
		panic("no return value specified for PurchaseMerch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.UserID, *model.MerchItem) error); ok {
		r0 = rf(ctx, buyer, m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepo_PurchaseMerch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PurchaseMerch'
type UserRepo_PurchaseMerch_Call struct {
	*mock.Call
}

// PurchaseMerch is a helper method to define mock.On call
//   - ctx context.Context
//   - buyer model.UserID
//   - m *model.MerchItem
func (_e *UserRepo_Expecter) PurchaseMerch(ctx interface{}, buyer interface{}, m interface{}) *UserRepo_PurchaseMerch_Call {
	return &UserRepo_PurchaseMerch_Call{Call: _e.mock.On("PurchaseMerch", ctx, buyer, m)}
}

func (_c *UserRepo_PurchaseMerch_Call) Run(run func(ctx context.Context, buyer model.UserID, m *model.MerchItem)) *UserRepo_PurchaseMerch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.UserID), args[2].(*model.MerchItem))
	})
	return _c
}

func (_c *UserRepo_PurchaseMerch_Call) Return(_a0 error) *UserRepo_PurchaseMerch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepo_PurchaseMerch_Call) RunAndReturn(run func(context.Context, model.UserID, *model.MerchItem) error) *UserRepo_PurchaseMerch_Call {
	_c.Call.Return(run)
	return _c
}

// TransferCoins provides a mock function with given fields: ctx, from, to, amount
func (_m *UserRepo) TransferCoins(ctx context.Context, from model.UserID, to model.UserID, amount model.NumCoins) error {
	ret := _m.Called(ctx, from, to, amount)

	if len(ret) == 0 {
		panic("no return value specified for TransferCoins")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.UserID, model.UserID, model.NumCoins) error); ok {
		r0 = rf(ctx, from, to, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepo_TransferCoins_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TransferCoins'
type UserRepo_TransferCoins_Call struct {
	*mock.Call
}

// TransferCoins is a helper method to define mock.On call
//   - ctx context.Context
//   - from model.UserID
//   - to model.UserID
//   - amount model.NumCoins
func (_e *UserRepo_Expecter) TransferCoins(ctx interface{}, from interface{}, to interface{}, amount interface{}) *UserRepo_TransferCoins_Call {
	return &UserRepo_TransferCoins_Call{Call: _e.mock.On("TransferCoins", ctx, from, to, amount)}
}

func (_c *UserRepo_TransferCoins_Call) Run(run func(ctx context.Context, from model.UserID, to model.UserID, amount model.NumCoins)) *UserRepo_TransferCoins_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.UserID), args[2].(model.UserID), args[3].(model.NumCoins))
	})
	return _c
}

func (_c *UserRepo_TransferCoins_Call) Return(_a0 error) *UserRepo_TransferCoins_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepo_TransferCoins_Call) RunAndReturn(run func(context.Context, model.UserID, model.UserID, model.NumCoins) error) *UserRepo_TransferCoins_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserRepo creates a new instance of UserRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepo {
	mock := &UserRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
