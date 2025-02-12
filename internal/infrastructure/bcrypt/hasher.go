package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func NewHasher() *Hasher {
	return &Hasher{}
}

type Hasher struct{}

func (c *Hasher) HashPassword(passwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (c *Hasher) VerifyPassword(passwd string, passwdHash []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(passwdHash, []byte(passwd)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
	}
	return true, nil
}
