package model

type UserRepo interface {
	GetUserByUsername(Username) (*User, error)
	PutUser(*User) (*User, error)
}
