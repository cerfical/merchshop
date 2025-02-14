package coins

import (
	"errors"
	"fmt"

	"github.com/cerfical/merchshop/internal/service/model"
)

func NewCoinService(users model.UserRepo) CoinService {
	return &coinService{users}
}

type CoinService interface {
	GetCoinBalance(model.Username) (model.NumCoins, error)
	SendCoins(from model.Username, to model.Username, amount model.NumCoins) error
}

type coinService struct {
	users model.UserRepo
}

func (s *coinService) GetCoinBalance(un model.Username) (model.NumCoins, error) {
	u, err := s.users.GetUserByUsername(un)
	if err != nil {
		return 0, fmt.Errorf("get user data: %w", err)
	}
	return u.Coins, nil
}

func (s *coinService) SendCoins(from model.Username, to model.Username, amount model.NumCoins) error {
	fromUser, err := s.users.GetUserByUsername(from)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			return model.ErrSenderNotExist
		}
		return fmt.Errorf("identify sender: %w", err)
	}

	toUser, err := s.users.GetUserByUsername(to)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			return model.ErrRecipientNotExist
		}
		return fmt.Errorf("identify recipient: %w", err)
	}

	if err := s.users.TransferCoins(fromUser.ID, toUser.ID, amount); err != nil {
		return fmt.Errorf("transfer coins: %w", err)
	}

	return nil
}
