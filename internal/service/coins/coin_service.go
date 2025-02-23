package coins

import (
	"errors"
	"fmt"

	"github.com/cerfical/merchshop/internal/service/model"
	"github.com/cerfical/merchshop/internal/service/repo"
)

func NewCoinService(users repo.UserRepo) CoinService {
	return &coinService{users}
}

type CoinService interface {
	GetUser(model.Username) (*model.User, error)
	SendCoins(from model.Username, to model.Username, amount model.NumCoins) error
	BuyItem(buyer model.Username, m *model.MerchItem) error
}

type coinService struct {
	users repo.UserRepo
}

func (s *coinService) GetUser(un model.Username) (*model.User, error) {
	u, err := s.users.GetUser(un)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

func (s *coinService) SendCoins(from model.Username, to model.Username, amount model.NumCoins) error {
	// Disallow coin transfers between the same user
	if from == to {
		return model.ErrSenderIsRecipient
	}

	fromUser, err := s.users.GetUser(from)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			return model.ErrSenderNotExist
		}
		return fmt.Errorf("identify sender: %w", err)
	}

	toUser, err := s.users.GetUser(to)
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

func (s *coinService) BuyItem(buyer model.Username, m *model.MerchItem) error {
	buyerUser, err := s.users.GetUser(buyer)
	if err != nil {
		return fmt.Errorf("identify buyer: %w", err)
	}

	if err := s.users.PurchaseMerch(buyerUser.ID, m); err != nil {
		return fmt.Errorf("purchase merch: %w", err)
	}
	return nil
}
