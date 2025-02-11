package api

import (
	"net/http"
	"time"

	"github.com/cerfical/merchshop/internal/log"
	"github.com/cerfical/merchshop/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type authHandler struct {
	users  model.UserStore
	config *AuthConfig

	log *log.Logger
}

func (h *authHandler) authUser(w http.ResponseWriter, r *http.Request) {
	uc, err := decode[userAuthRequest](r.Body)
	if err != nil {
		encode(w, http.StatusBadRequest, errorResponse{
			Errors: "The request body is either not valid JSON or contains invalid fields",
		})
		return
	}

	u, err := h.users.GetUser(&model.UserCreds{
		Name:     uc.Username,
		Password: uc.Password,
	})
	if err != nil {
		internalError("Unable to create a new user due to database failure", err, h.log)(w, r)
		return
	}

	// TODO: Hash the password
	passwdHash := uc.Password
	if u.Password != passwdHash {
		encode(w, http.StatusUnauthorized, errorResponse{
			Errors: "The credentials provided are invalid",
		})
		return
	}

	token, err := h.makeToken(u.Name)
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
