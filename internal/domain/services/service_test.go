package services_test

import (
	"testing"

	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/cerfical/merchshop/internal/domain/services"
	"github.com/cerfical/merchshop/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTest))
}

type ServiceTest struct {
	suite.Suite

	users    *mocks.UserRepo
	hasher   *mocks.PasswordHasher
	tokenGen *mocks.TokenGen
}

func (t *ServiceTest) SetupSubTest() {
	t.users = mocks.NewUserRepo(t.T())
	t.hasher = mocks.NewPasswordHasher(t.T())
	t.tokenGen = mocks.NewTokenGen(t.T())
}

func (t *ServiceTest) TestAuthUser() {
	user := model.User{
		Username:     "testuser",
		PasswordHash: []byte("321"),
	}

	tests := []struct {
		Name     string
		Username model.Username
		Password model.Password
		Token    string

		Setup func()
		Err   assert.ErrorAssertionFunc
	}{
		{
			Name:     "auth_ok",
			Username: "testuser",
			Password: "123",
			Token:    "123321",

			Setup: func() {
				e := t.hasher.EXPECT()
				e.HashPassword("123").
					Return([]byte("321"), nil)
				e.VerifyPassword("123", []byte("321")).
					Return(true, nil)

				t.users.EXPECT().
					PutUser(&user).
					Return(&user, nil)

				t.tokenGen.EXPECT().
					NewToken(&user).
					Return("123321", nil)
			},
			Err: assert.NoError,
		},

		{
			Name:     "auth_fail",
			Username: "testuser",
			Password: "124",
			Token:    "",

			Setup: func() {
				e := t.hasher.EXPECT()
				e.HashPassword("124").
					Return([]byte("421"), nil)
				e.VerifyPassword("124", []byte("321")).
					Return(false, nil)

				t.users.EXPECT().
					PutUser(&model.User{
						Username:     "testuser",
						PasswordHash: []byte("421"),
					}).
					Return(&user, nil)
			},
			Err: func(t assert.TestingT, err error, args ...any) bool {
				return assert.Error(t, model.ErrAuthFail, args...)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func() {
			test.Setup()

			service := services.NewAuthService(t.users, t.hasher, t.tokenGen)
			token, err := service.AuthUser(model.UserCreds{
				Username: test.Username,
				Password: test.Password,
			})

			test.Err(t.T(), err)
			t.Equal(test.Token, token)
		})
	}
}

func (t *ServiceTest) TestGetUserCoinBalance() {
	user := model.User{
		Username:     "testuser",
		PasswordHash: []byte("321"),
		Coins:        9,
	}

	tests := []struct {
		Name string

		Username model.Username
		Coins    model.NumCoins

		Setup func()
		Err   assert.ErrorAssertionFunc
	}{
		{
			Name:     "existing_user_ok",
			Username: "testuser",
			Coins:    9,

			Setup: func() {
				t.users.EXPECT().
					GetUserByUsername(model.Username("testuser")).
					Return(&user, nil)
			},
			Err: assert.NoError,
		},

		{
			Name:     "unknown_user_fail",
			Username: "badtestuser",

			Setup: func() {
				t.users.EXPECT().
					GetUserByUsername(model.Username("badtestuser")).
					Return(nil, model.ErrNotExist)
			},
			Err: func(t assert.TestingT, err error, args ...any) bool {
				return assert.Error(t, model.ErrNotExist, args...)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func() {
			test.Setup()

			service := services.NewCoinService(t.users)
			coins, err := service.GetCoinBalance(test.Username)

			test.Err(t.T(), err)
			t.Equal(test.Coins, coins)
		})
	}
}
