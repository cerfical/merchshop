package services

import (
	"errors"
	"fmt"

	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/cerfical/merchshop/internal/domain/repo"
)

func NewCoinService(users repo.UserRepo) CoinService {
	return &coinService{users}
}

type CoinService interface {
	GetCoinBalance(model.Username) (model.NumCoins, error)
}

type coinService struct {
	users repo.UserRepo
}

func (s *coinService) GetCoinBalance(un model.Username) (model.NumCoins, error) {
	u, err := s.users.GetUserByUsername(un)
	if err != nil {
		if errors.Is(err, model.ErrNotExist) {
			return 0, model.ErrNotExist
		}
		return 0, fmt.Errorf("coin balance checkout: %w", err)
	}
	return u.Coins, nil
}
