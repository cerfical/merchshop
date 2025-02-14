package repo

import "github.com/cerfical/merchshop/internal/service/model"

type UserRepo interface {
	GetUser(model.Username) (*model.User, error)
	CreateUser(model.Username, model.PasswordHash, model.NumCoins) (*model.User, error)

	TransferCoins(from model.UserID, to model.UserID, amount model.NumCoins) error
	PurchaseMerch(buyer model.UserID, m *model.MerchItem) error
}
