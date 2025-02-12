package jwt

import (
	"time"

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

func (a *TokenAuth) NewToken(u *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": u.Username,
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

	return signedToken, nil
}
