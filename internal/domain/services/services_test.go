package services_test

import (
	"testing"

	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/cerfical/merchshop/internal/domain/services"
	"github.com/cerfical/merchshop/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var user = model.User{
	Username:     "testuser",
	PasswordHash: []byte("321"),
}

func TestServices(t *testing.T) {
	suite.Run(t, new(ServicesTest))
}

type ServicesTest struct {
	suite.Suite
}

func (t *ServicesTest) TestAuthUser() {
	tests := []struct {
		Name     string
		Username model.Username
		Password model.Password
		Token    string

		Hasher   func() *mocks.PasswordHasher
		TokenGen func() *mocks.TokenGen
		UserRepo func() *mocks.UserRepo

		Err assert.ErrorAssertionFunc
	}{
		{
			Name:     "auth_ok",
			Username: "testuser",
			Password: "123",
			Token:    "123321",

			Hasher: func() *mocks.PasswordHasher {
				hasher := mocks.NewPasswordHasher(t.T())

				e := hasher.EXPECT()
				e.HashPassword("123").
					Return([]byte("321"), nil)
				e.VerifyPassword("123", []byte("321")).
					Return(true, nil)

				return hasher
			},

			UserRepo: func() *mocks.UserRepo {
				users := mocks.NewUserRepo(t.T())
				users.EXPECT().
					PutUser(&user).
					Return(&user, nil)
				return users
			},

			TokenGen: func() *mocks.TokenGen {
				tokens := mocks.NewTokenGen(t.T())
				tokens.EXPECT().
					NewToken(&user).
					Return("123321", nil)
				return tokens
			},

			Err: assert.NoError,
		},

		{
			Name:     "auth_fail",
			Username: "testuser",
			Password: "124",
			Token:    "",

			Hasher: func() *mocks.PasswordHasher {
				hasher := mocks.NewPasswordHasher(t.T())
				hasher.EXPECT().
					HashPassword("124").
					Return([]byte("421"), nil)

				hasher.EXPECT().
					VerifyPassword("124", []byte("321")).
					Return(false, nil)

				return hasher
			},

			UserRepo: func() *mocks.UserRepo {
				users := mocks.NewUserRepo(t.T())
				users.EXPECT().
					PutUser(&model.User{
						Username:     "testuser",
						PasswordHash: []byte("421"),
					}).
					Return(&user, nil)
				return users
			},

			TokenGen: func() *mocks.TokenGen {
				return mocks.NewTokenGen(t.T())
			},

			Err: func(t assert.TestingT, err error, args ...any) bool {
				return assert.Error(t, model.ErrAuthFail, args...)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func() {
			service := services.NewAuthService(test.UserRepo(), test.Hasher(), test.TokenGen())
			token, err := service.AuthUser(model.UserCreds{
				Username: test.Username,
				Password: test.Password,
			})

			test.Err(t.T(), err)
			t.Equal(test.Token, token)
		})
	}
}
