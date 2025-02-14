package coins_test

import (
	"testing"

	"github.com/cerfical/merchshop/internal/domain/coins"
	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/cerfical/merchshop/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestCoinService(t *testing.T) {
	suite.Run(t, new(CoinServiceTest))
}

type CoinServiceTest struct {
	suite.Suite

	users *mocks.UserRepo
}

func (t *CoinServiceTest) SetupSubTest() {
	t.users = mocks.NewUserRepo(t.T())
}

func (t *CoinServiceTest) TestGetUserCoinBalance() {
	user := model.User{
		Username: "testuser",
		Coins:    9,
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
					Return(nil, model.ErrUserNotExist)
			},
			Err: func(t assert.TestingT, err error, args ...any) bool {
				return assert.ErrorIs(t, err, model.ErrUserNotExist, args...)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func() {
			test.Setup()

			service := coins.NewCoinService(t.users)
			coins, err := service.GetCoinBalance(test.Username)

			test.Err(t.T(), err)
			t.Equal(test.Coins, coins)
		})
	}
}
