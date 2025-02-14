package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cerfical/merchshop/internal/domain/auth"
	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/cerfical/merchshop/internal/log"
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

	uc, err := model.NewUserCreds(authReq.Username, authReq.Password)
	if err != nil {
		badRequestHandler(fmt.Sprintf("The provided credentials are invalid: %v", err))(w, r)
		return
	}

	token, err := h.authService.AuthUser(uc)
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
