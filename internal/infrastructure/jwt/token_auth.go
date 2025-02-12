package jwt

import (
	"time"

	"github.com/cerfical/merchshop/internal/domain/auth"
	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/golang-jwt/jwt/v5"
)

func NewTokenAuth(config TokenConfig) *TokenAuth {
	return &TokenAuth{config}
}

type TokenAuth struct {
	config TokenConfig
}

func (a *TokenAuth) NewToken(u *model.User) (auth.Token, error) {
	claims := jwt.MapClaims{
		"sub": u.Username,
		"iss": "merchshop",
		"iat": time.Now().Unix(),
	}

	if lifetime := a.config.Lifetime; lifetime > 0 {
		claims["exp"] = time.Now().Add(lifetime).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	signedToken, err := token.SignedString(a.config.Secret)
	if err != nil {
		return "", err
	}

	return auth.Token(signedToken), nil
}
