package jwt

import (
	"errors"
	"time"

	"github.com/cerfical/merchshop/internal/domain/auth"
	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/golang-jwt/jwt/v5"
)

func NewTokenAuth(config *TokenConfig) *TokenAuth {
	return &TokenAuth{config.Secret, config.Lifetime}
}

type TokenAuth struct {
	secret   []byte
	lifetime time.Duration
}

func (a *TokenAuth) IssueToken(u model.Username) (auth.Token, error) {
	claims := jwt.MapClaims{
		"sub": u,
		"iss": "merchshop",
		"iat": time.Now().Unix(),
	}

	if a.lifetime > 0 {
		claims["exp"] = time.Now().Add(a.lifetime).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	signedToken, err := token.SignedString(a.secret)
	if err != nil {
		return "", err
	}

	return auth.Token(signedToken), nil
}

func (a *TokenAuth) AuthToken(token auth.Token) (model.Username, error) {
	t, err := jwt.Parse(string(token), func(*jwt.Token) (any, error) {
		return a.secret, nil
	})
	if err != nil {
		return "", err
	}

	if !t.Valid {
		return "", errors.New("token is invalid")
	}

	u, err := t.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return model.Username(u), nil
}
