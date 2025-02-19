package coins

import (
	"context"
	"errors"
	"fmt"

	"github.com/cerfical/merchshop/internal/service/model"
	"github.com/cerfical/merchshop/internal/service/repo"
)

func NewCoinService(users repo.UserRepo) CoinService {
	return &coinService{users}
}

type CoinService interface {
	GetUser(context.Context, model.Username) (*model.User, error)
	SendCoins(ctx context.Context, from model.Username, to model.Username, amount model.NumCoins) error
	BuyItem(ctx context.Context, buyer model.Username, m *model.MerchItem) error
}

type coinService struct {
	users repo.UserRepo
}

func (s *coinService) GetUser(ctx context.Context, un model.Username) (*model.User, error) {
	u, err := s.users.GetUser(ctx, un)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

func (s *coinService) SendCoins(ctx context.Context, from model.Username, to model.Username, amount model.NumCoins) error {
	// Disallow coin transfers to the same user
	if from == to {
		return model.ErrSenderIsRecipient
	}

	fromUser, err := s.users.GetUser(ctx, from)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			return model.ErrSenderNotExist
		}
		return fmt.Errorf("identify sender: %w", err)
	}

	toUser, err := s.users.GetUser(ctx, to)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			return model.ErrRecipientNotExist
		}
		return fmt.Errorf("identify recipient: %w", err)
	}

	if err := s.users.TransferCoins(ctx, fromUser.ID, toUser.ID, amount); err != nil {
		return fmt.Errorf("transfer coins: %w", err)
	}

	return nil
}

func (s *coinService) BuyItem(ctx context.Context, buyer model.Username, m *model.MerchItem) error {
	buyerUser, err := s.users.GetUser(ctx, buyer)
	if err != nil {
		return fmt.Errorf("identify buyer: %w", err)
	}

	if err := s.users.PurchaseMerch(ctx, buyerUser.ID, m); err != nil {
		return fmt.Errorf("purchase merch: %w", err)
	}
	return nil
}
