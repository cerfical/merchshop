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

	// TODO: Think of a better solution? Limit the number of attempts?
	var u *model.User
	for {
		u, err = s.users.GetUser(username)
		if errors.Is(err, model.ErrUserNotExist) {
			// Proceed to creating a new user
			u, err = s.users.CreateUser(username, passwdHash, defaultCoinBalance)
			if errors.Is(err, model.ErrUserExist) {
				// Some other transaction created a user with the given username before us,
				// try to retrieve the user yet another time
				// In the absence of user deletions, that should be OK,
				// but otherwise there is a chance, that the user will be deleted before we access it,
				// so here is the loop
				continue
			} else if err != nil {
				return "", fmt.Errorf("create new user: %w", err)
			}
		} else if err != nil {
			return "", fmt.Errorf("get user: %w", err)
		}

		// Success: the user does exist and there were no errors in retrieving it
		break
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
