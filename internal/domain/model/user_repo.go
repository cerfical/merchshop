package model

type UserRepo interface {
	GetUserByUsername(Username) (*User, error)
	PutUser(*User) (*User, error)
	TransferCoins(from UserID, to UserID, amount NumCoins) error
}
