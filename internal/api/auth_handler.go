package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cerfical/merchshop/internal/log"
	"github.com/cerfical/merchshop/internal/service/auth"
	"github.com/cerfical/merchshop/internal/service/model"
)

type authHandler struct {
	authService auth.AuthService
	log         *log.Logger
}

func (h *authHandler) authUser(w http.ResponseWriter, r *http.Request) {
	authReq, err := readRequest[authRequest](r.Body)
	if err != nil {
		// TODO: More descriptive error messages?
		badRequestHandler("The request body is malformed")(w, r)
		return
	}

	username, err := model.NewUsername(authReq.Username)
	if err != nil {
		badRequestHandler(fmt.Sprintf("The provided username is invalid: %v", err))(w, r)
		return
	}

	passwd, err := model.NewPassword(authReq.Password)
	if err != nil {
		badRequestHandler(fmt.Sprintf("The provided password is invalid: %v", err))(w, r)
		return
	}

	token, err := h.authService.AuthUser(r.Context(), username, passwd)
	if err != nil {
		if errors.Is(err, model.ErrAuthFail) {
			unauthorizedHandler("The provided credentials are invalid")(w, r)
		} else {
			internalErrorHandler(h.log, "User authentication failed", err)(w, r)
		}
		return
	}

	writeResponse(w, http.StatusOK, authResponse{
		Token: string(token),
	})
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}
