package api

import (
	"net/http"
	"strconv"

	"github.com/cerfical/merchshop/internal/log"
	"github.com/cerfical/merchshop/internal/model"
)

type authHandler struct {
	users model.UserStore
	log   *log.Logger
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

	// TODO: Implement a more secure authentication mechanism
	token := strconv.Itoa(u.ID)
	encode(w, http.StatusOK, authResponse{
		Token: token,
	})
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
