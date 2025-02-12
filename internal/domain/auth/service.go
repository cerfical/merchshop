package auth

import (
	"fmt"

	"github.com/cerfical/merchshop/internal/domain/model"
)

func NewService(users UserRepo, hasher PasswordHasher, tokens TokenGen) Service {
	return &service{users, hasher, tokens}
}

type Service interface {
	AuthUser(model.UserCreds) (Token, error)
}

type UserRepo interface {
	PutUser(*model.User) (*model.User, error)
}

type PasswordHasher interface {
	HashPassword(passwd string) ([]byte, error)
	VerifyPassword(passwd string, passwdHash []byte) (bool, error)
}

type TokenGen interface {
	NewToken(u *model.User) (Token, error)
}

type Token string

type service struct {
	users  UserRepo
	hasher PasswordHasher
	tokens TokenGen
}

func (s *service) AuthUser(uc model.UserCreds) (Token, error) {
	passwdHash, err := s.hasher.HashPassword(string(uc.Password))
	if err != nil {
		return "", fmt.Errorf("password hashing: %w", err)
	}

	u, err := s.users.PutUser(&model.User{
		Username:     uc.Username,
		PasswordHash: passwdHash,
	})

	if err != nil {
		return "", fmt.Errorf("accessing user storage: %w", err)
	}

	ok, err := s.hasher.VerifyPassword(string(uc.Password), u.PasswordHash)
	if err != nil {
		return "", fmt.Errorf("password verification: %w", err)
	}

	if !ok {
		return "", model.ErrAuthFail
	}

	token, err := s.tokens.NewToken(u)
	if err != nil {
		return "", fmt.Errorf("token generation: %w", err)
	}

	return token, nil
}
