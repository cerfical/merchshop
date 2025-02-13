package auth

import (
	"errors"
	"fmt"

	"github.com/cerfical/merchshop/internal/domain/model"
)

func NewAuthService(auth TokenAuth, users model.UserRepo, hasher PasswordHasher) AuthService {
	return &authService{users, hasher, auth}
}

type Token string

type AuthService interface {
	AuthUser(model.UserCreds) (Token, error)
}

type PasswordHasher interface {
	HashPassword(model.Password) (model.Hash, error)
	VerifyPassword(model.Password, model.Hash) error
}

type TokenAuth interface {
	IssueToken(model.Username) (Token, error)
}

type authService struct {
	users  model.UserRepo
	hasher PasswordHasher
	tokens TokenAuth
}

func (s *authService) AuthUser(uc model.UserCreds) (Token, error) {
	passwdHash, err := s.hasher.HashPassword(uc.Password)
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

	if err := s.hasher.VerifyPassword(uc.Password, u.PasswordHash); err != nil {
		if errors.Is(err, model.ErrAuthFail) {
			return "", model.ErrAuthFail
		}
		return "", fmt.Errorf("password verification: %w", err)
	}

	token, err := s.tokens.IssueToken(u.Username)
	if err != nil {
		return "", fmt.Errorf("token generation: %w", err)
	}

	return token, nil
}
