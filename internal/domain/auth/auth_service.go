package auth

import (
	"fmt"

	"github.com/cerfical/merchshop/internal/domain/model"
)

type Token string

type AuthService interface {
	AuthUser(model.Username, model.Password) (Token, error)
	AuthToken(Token) (model.Username, error)
}

type PasswordHasher interface {
	HashPassword(model.Password) (model.PasswordHash, error)
	VerifyPassword(model.Password, model.PasswordHash) error
}

type TokenAuth interface {
	IssueToken(model.Username) (Token, error)
	AuthToken(Token) (model.Username, error)
}

type authService struct {
	users  model.UserRepo
	hasher PasswordHasher
	auth   TokenAuth
}

func (s *authService) AuthUser(username model.Username, passwd model.Password) (Token, error) {
	passwdHash, err := s.hasher.HashPassword(passwd)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	u, err := s.users.PutUser(&model.User{
		Username:     username,
		PasswordHash: passwdHash,
	})

	if err != nil {
		return "", fmt.Errorf("read user data: %w", err)
	}

	if err := s.hasher.VerifyPassword(passwd, u.PasswordHash); err != nil {
		return "", fmt.Errorf("verify password: %w", err)
	}

	token, err := s.auth.IssueToken(u.Username)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (s *authService) AuthToken(token Token) (model.Username, error) {
	return s.auth.AuthToken(token)
}

func NewAuthService(auth TokenAuth, users model.UserRepo, hasher PasswordHasher) AuthService {
	return &authService{users, hasher, auth}
}
