package auth_test

import (
	"testing"

	"github.com/cerfical/merchshop/internal/domain/auth"
	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/cerfical/merchshop/internal/mocks"
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
	user := model.User{
		Username:     "testuser",
		PasswordHash: model.PasswordHash("321"),
	}

	tests := []struct {
		Name     string
		Username model.Username
		Password model.Password
		Token    auth.Token

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
				e.HashPassword(model.Password("123")).
					Return(model.PasswordHash("321"), nil)
				e.VerifyPassword(model.Password("123"), model.PasswordHash("321")).
					Return(nil)

				t.users.EXPECT().
					PutUser(&user).
					Return(&user, nil)

				t.auth.EXPECT().
					IssueToken(user.Username).
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
				e.HashPassword(model.Password("124")).
					Return(model.PasswordHash("421"), nil)
				e.VerifyPassword(model.Password("124"), model.PasswordHash("321")).
					Return(model.ErrAuthFail)

				t.users.EXPECT().
					PutUser(&model.User{
						Username:     "testuser",
						PasswordHash: model.PasswordHash("421"),
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

			service := auth.NewAuthService(t.auth, t.users, t.hasher)
			token, err := service.AuthUser(test.Username, test.Password)

			test.Err(t.T(), err)
			t.Equal(test.Token, token)
		})
	}
}
