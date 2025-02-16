package auth

import (
	"errors"
	"fmt"

	"github.com/cerfical/merchshop/internal/service/model"
	"github.com/cerfical/merchshop/internal/service/repo"
)

const (
	defaultCoinBalance = 1000
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
	users  repo.UserRepo
	hasher PasswordHasher
	auth   TokenAuth
}

func (s *authService) AuthUser(username model.Username, passwd model.Password) (Token, error) {
	passwdHash, err := s.hasher.HashPassword(passwd)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	// TODO: Think of a better solution?
	u, err := s.users.GetUser(username)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			u, err = s.users.CreateUser(username, passwdHash, defaultCoinBalance)
			if err != nil {
				return "", fmt.Errorf("create user: %w", err)
			}
		} else {
			return "", fmt.Errorf("find user by username: %w", err)
		}
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

func NewAuthService(auth TokenAuth, users repo.UserRepo, hasher PasswordHasher) AuthService {
	return &authService{users, hasher, auth}
}
