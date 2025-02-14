package model

type UserRepo interface {
	GetUser(Username) (*User, error)
	CreateUser(Username, PasswordHash, NumCoins) (*User, error)

	TransferCoins(from UserID, to UserID, amount NumCoins) error
	PurchaseMerch(buyer UserID, m *MerchItem) error
}
