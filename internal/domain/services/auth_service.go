package services

import (
	"fmt"

	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/cerfical/merchshop/internal/domain/repo"
)

func NewAuthService(users repo.UserRepo, hasher PasswordHasher, tokens TokenGen) AuthService {
	return &authService{users, hasher, tokens}
}

type PasswordHasher interface {
	HashPassword(passwd string) ([]byte, error)
	VerifyPassword(passwd string, passwdHash []byte) (bool, error)
}

type TokenGen interface {
	NewToken(u *model.User) (string, error)
}

type AuthService interface {
	AuthUser(model.UserCreds) (string, error)
}

type authService struct {
	users  repo.UserRepo
	hasher PasswordHasher
	tokens TokenGen
}

func (s *authService) AuthUser(uc model.UserCreds) (string, error) {
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
