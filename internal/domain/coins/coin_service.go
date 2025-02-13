package coins

import (
	"errors"
	"fmt"

	"github.com/cerfical/merchshop/internal/domain/model"
)

func NewCoinService(users model.UserRepo) CoinService {
	return &coinService{users}
}

type CoinService interface {
	GetCoinBalance(model.Username) (model.NumCoins, error)
}

type coinService struct {
	users model.UserRepo
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
