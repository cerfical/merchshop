package auth_test

import (
	"testing"

	"github.com/cerfical/merchshop/internal/service/auth"
	"github.com/cerfical/merchshop/internal/service/model"
	"github.com/cerfical/merchshop/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTest))
}

type AuthServiceTest struct {
	suite.Suite

	users  *mocks.UserRepo
	hasher *mocks.PasswordHasher
	auth   *mocks.TokenAuth
}

func (t *AuthServiceTest) SetupSubTest() {
	t.users = mocks.NewUserRepo(t.T())
	t.hasher = mocks.NewPasswordHasher(t.T())
	t.auth = mocks.NewTokenAuth(t.T())
}

func (t *AuthServiceTest) TestAuthUser() {
	users := []model.User{{
		Username:     "test_user",
		PasswordHash: model.PasswordHash("321"),
	}, {
		Username:     "new_test_user",
		PasswordHash: model.PasswordHash("421"),
	}}

	tests := []struct {
		Name     string
		Username model.Username
		Password model.Password
		Token    auth.Token

		Setup func()
		Err   assert.ErrorAssertionFunc
	}{
		{
			Name:     "ok",
			Username: "test_user",
			Password: "123",
			Token:    "123321",

			Setup: func() {
				e := t.hasher.EXPECT()
				e.HashPassword(model.Password("123")).
					Return(model.PasswordHash("321"), nil)
				e.VerifyPassword(model.Password("123"), model.PasswordHash("321")).
					Return(nil)

				t.users.EXPECT().
					GetUser(model.Username("test_user")).
					Return(&users[0], nil)

				t.auth.EXPECT().
					IssueToken(model.Username("test_user")).
					Return("123321", nil)
			},
			Err: assert.NoError,
		},

		{
			Name:     "new_user",
			Username: "new_test_user",
			Password: "124",
			Token:    "124421",

			Setup: func() {
				he := t.hasher.EXPECT()
				he.HashPassword(model.Password("124")).
					Return(model.PasswordHash("421"), nil)
				he.VerifyPassword(model.Password("124"), model.PasswordHash("421")).
					Return(nil)

				ue := t.users.EXPECT()
				ue.GetUser(model.Username("new_test_user")).
					Return(nil, model.ErrUserNotExist)
				ue.CreateUser(model.Username("new_test_user"), model.PasswordHash("421"), model.NumCoins(1000)).
					Return(&users[1], nil)

				t.auth.EXPECT().
					IssueToken(model.Username("new_test_user")).
					Return("124421", nil)
			},
			Err: assert.NoError,
		},

		{
			Name:     "bad_password",
			Username: "test_user",
			Password: "bad_123",
			Token:    "",

			Setup: func() {
				e := t.hasher.EXPECT()
				e.HashPassword(model.Password("bad_123")).
					Return(model.PasswordHash("bad_321"), nil)
				e.VerifyPassword(model.Password("bad_123"), model.PasswordHash("321")).
					Return(model.ErrAuthFail)

				t.users.EXPECT().
					GetUser(model.Username("test_user")).
					Return(&users[0], nil)
			},
			Err: func(t assert.TestingT, err error, args ...any) bool {
				return assert.Error(t, model.ErrAuthFail, args...)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func() {
			test.Setup()

			service := auth.NewAuthService(t.auth, t.users, t.hasher)
			token, err := service.AuthUser(test.Username, test.Password)

			test.Err(t.T(), err)
			t.Equal(test.Token, token)
		})
	}
}
