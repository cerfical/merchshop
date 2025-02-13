package bcrypt

import (
	"errors"

	"github.com/cerfical/merchshop/internal/domain/model"
	"golang.org/x/crypto/bcrypt"
)

func NewHasher() *Hasher {
	return &Hasher{}
}

type Hasher struct{}

func (c *Hasher) HashPassword(passwd model.Password) (model.Hash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (c *Hasher) VerifyPassword(passwd model.Password, passwdHash model.Hash) error {
	if err := bcrypt.CompareHashAndPassword(passwdHash, []byte(passwd)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return model.ErrAuthFail
		}
		return err
	}
	return nil
}
