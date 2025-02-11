package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/cerfical/merchshop/internal/log"
	"github.com/cerfical/merchshop/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
	users  model.UserStore
	config *AuthConfig

	log *log.Logger
}

func (h *authHandler) authUser(w http.ResponseWriter, r *http.Request) {
	uc, err := decode[userAuthRequest](r.Body)
	if err != nil {
		encode(w, http.StatusBadRequest, &errorResponse{
			Errors: "The request body is either not valid JSON or contains invalid fields",
		})
		return
	}

	passwd := []byte(uc.Password)
	passwdHash, err := bcrypt.GenerateFromPassword(passwd, bcrypt.DefaultCost)
	if err != nil {
		internalError("Password hash generation failed", err, h.log)(w, r)
		return
	}

	u, err := h.users.GetUser(&model.UserCreds{
		Name:         uc.Username,
		PasswordHash: passwdHash,
	})
	if err != nil {
		internalError("Unable to create a new user due to database failure", err, h.log)(w, r)
		return
	}

	// Authenticate the user
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, passwd); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			encode(w, http.StatusUnauthorized, &errorResponse{
				Errors: "The credentials provided are invalid",
			})
		} else {
			internalError("Password validation failed", err, h.log)(w, r)
		}
		return
	}

	token, err := h.makeToken(uc.Username)
	if err != nil {
		internalError("Failed to generate a JWT token for the user", err, h.log)(w, r)
		return
	}

	encode(w, http.StatusOK, authResponse{
		Token: token,
	})
}

func (h *authHandler) makeToken(user string) (string, error) {
	claims := jwt.MapClaims{
		"sub": user,
		"iss": "merchshop",
		"iat": time.Now().Unix(),
	}

	if h.config.Token.Lifetime > 0 {
		claims["exp"] = time.Now().Add(h.config.Token.Lifetime).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString(h.config.Token.Secret)
}

type userAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Errors string `json:"errors"`
}
