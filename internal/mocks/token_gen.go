// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	model "github.com/cerfical/merchshop/internal/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// TokenGen is an autogenerated mock type for the TokenGen type
type TokenGen struct {
	mock.Mock
}

type TokenGen_Expecter struct {
	mock *mock.Mock
}

func (_m *TokenGen) EXPECT() *TokenGen_Expecter {
	return &TokenGen_Expecter{mock: &_m.Mock}
}

// NewToken provides a mock function with given fields: u
func (_m *TokenGen) NewToken(u *model.User) (string, error) {
	ret := _m.Called(u)

	if len(ret) == 0 {
		panic("no return value specified for NewToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.User) (string, error)); ok {
		return rf(u)
	}
	if rf, ok := ret.Get(0).(func(*model.User) string); ok {
		r0 = rf(u)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*model.User) error); ok {
		r1 = rf(u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenGen_NewToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewToken'
type TokenGen_NewToken_Call struct {
	*mock.Call
}

// NewToken is a helper method to define mock.On call
//   - u *model.User
func (_e *TokenGen_Expecter) NewToken(u interface{}) *TokenGen_NewToken_Call {
	return &TokenGen_NewToken_Call{Call: _e.mock.On("NewToken", u)}
}

func (_c *TokenGen_NewToken_Call) Run(run func(u *model.User)) *TokenGen_NewToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.User))
	})
	return _c
}

func (_c *TokenGen_NewToken_Call) Return(_a0 string, _a1 error) *TokenGen_NewToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TokenGen_NewToken_Call) RunAndReturn(run func(*model.User) (string, error)) *TokenGen_NewToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewTokenGen creates a new instance of TokenGen. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenGen(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenGen {
	mock := &TokenGen{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
