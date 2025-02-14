package coins_test

import (
	"testing"

	"github.com/cerfical/merchshop/internal/mocks"
	"github.com/cerfical/merchshop/internal/service/coins"
	"github.com/cerfical/merchshop/internal/service/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		Username: "test_user",
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
			Name:     "existing_user",
			Username: "test_user",
			Coins:    9,

			Setup: func() {
				t.users.EXPECT().
					GetUserByUsername(model.Username("test_user")).
					Return(&user, nil)
			},
			Err: assert.NoError,
		},

		{
			Name:     "user_not_found",
			Username: "bad_test_user",

			Setup: func() {
				t.users.EXPECT().
					GetUserByUsername(model.Username("bad_test_user")).
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

func (t *CoinServiceTest) TestSendCoins() {
	users := []model.User{{
		ID:       0,
		Username: "test_sender",
		Coins:    9,
	}, {
		ID:       1,
		Username: "test_recipient",
		Coins:    9,
	}}

	tests := []struct {
		Name string

		Sender    model.Username
		Recipient model.Username
		Amount    model.NumCoins

		Setup func()
		Err   assert.ErrorAssertionFunc
	}{
		{
			Name:      "ok",
			Sender:    "test_sender",
			Recipient: "test_recipient",
			Amount:    9,

			Setup: func() {
				e := t.users.EXPECT()
				e.GetUserByUsername(model.Username("test_sender")).
					Return(&users[0], nil)
				e.GetUserByUsername(model.Username("test_recipient")).
					Return(&users[1], nil)
				e.TransferCoins(model.UserID(0), model.UserID(1), model.NumCoins(9)).
					Return(nil)
			},
			Err: assert.NoError,
		},

		{
			Name:      "invalid_sender",
			Sender:    "bad_test_sender",
			Recipient: "test_recipient",
			Amount:    9,

			Setup: func() {
				e := t.users.EXPECT()
				e.GetUserByUsername(model.Username("bad_test_sender")).
					Return(nil, model.ErrUserNotExist)
				e.GetUserByUsername(mock.Anything).
					Return(&users[1], nil)
			},
			Err: func(t assert.TestingT, err error, args ...any) bool {
				return assert.ErrorIs(t, err, model.ErrSenderNotExist, args...)
			},
		},

		{
			Name:      "invalid_recipient",
			Sender:    "test_sender",
			Recipient: "bad_test_recipient",
			Amount:    9,

			Setup: func() {
				e := t.users.EXPECT()
				e.GetUserByUsername(model.Username("bad_test_recipient")).
					Return(nil, model.ErrUserNotExist)
				e.GetUserByUsername(mock.Anything).
					Return(&users[1], nil)
			},
			Err: func(t assert.TestingT, err error, args ...any) bool {
				return assert.ErrorIs(t, err, model.ErrRecipientNotExist, args...)
			},
		},

		{
			Name:      "insufficient_coins",
			Sender:    "test_sender",
			Recipient: "test_recipient",
			Amount:    10,

			Setup: func() {
				e := t.users.EXPECT()
				e.GetUserByUsername(model.Username("test_sender")).
					Return(&users[0], nil)
				e.GetUserByUsername(model.Username("test_recipient")).
					Return(&users[1], nil)
				e.TransferCoins(model.UserID(0), model.UserID(1), model.NumCoins(10)).
					Return(model.ErrNotEnoughCoins)
			},
			Err: func(t assert.TestingT, err error, args ...any) bool {
				return assert.ErrorIs(t, err, model.ErrNotEnoughCoins, args...)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func() {
			test.Setup()

			service := coins.NewCoinService(t.users)
			err := service.SendCoins(test.Sender, test.Recipient, test.Amount)
			test.Err(t.T(), err)
		})
	}
}
