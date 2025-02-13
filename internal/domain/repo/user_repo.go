package repo

import "github.com/cerfical/merchshop/internal/domain/model"

type UserRepo interface {
	GetUserByUsername(model.Username) (*model.User, error)
	PutUser(*model.User) (*model.User, error)
}
