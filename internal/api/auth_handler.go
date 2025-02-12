package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/cerfical/merchshop/internal/domain/services"
	"github.com/cerfical/merchshop/internal/log"
)

type authHandler struct {
	auth services.AuthService
	log  *log.Logger
}

func (h *authHandler) authUser(w http.ResponseWriter, r *http.Request) {
	authReq, err := readRequest[authRequest](r.Body)
	if err != nil {
		// TODO: More specific error messages?
		writeErrorResponse(w, http.StatusBadRequest, "The request body is malformed")
		return
	}

	uc, err := model.NewUserCreds(authReq.Username, authReq.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("The provided credentials are invalid: %v", err))
		return
	}

	token, err := h.auth.AuthUser(uc)
	if err != nil {
		if errors.Is(err, model.ErrAuthFail) {
			writeErrorResponse(w, http.StatusUnauthorized, "The provided credentials are invalid")
		} else {
			h.log.Error("User authentication failed", err)
			writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	writeResponse(w, http.StatusOK, authResponse{
		Token: string(token),
	})
}

func writeErrorResponse(w http.ResponseWriter, status int, msg string) {
	writeResponse(w, status, errorResponse{
		Errors: msg,
	})
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Errors string `json:"errors"`
}
