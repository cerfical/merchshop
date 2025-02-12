package repo

import "github.com/cerfical/merchshop/internal/domain/model"

type UserRepo interface {
	PutUser(*model.User) (*model.User, error)
}
