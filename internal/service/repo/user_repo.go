package repo

import (
	"context"

	"github.com/cerfical/merchshop/internal/service/model"
)

type UserRepo interface {
	GetUser(context.Context, model.Username) (*model.User, error)
	CreateUser(context.Context, model.Username, model.PasswordHash, model.NumCoins) (*model.User, error)

	TransferCoins(ctx context.Context, from model.UserID, to model.UserID, amount model.NumCoins) error
	PurchaseMerch(ctx context.Context, buyer model.UserID, m *model.MerchItem) error
}
