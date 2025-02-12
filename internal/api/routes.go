package api

import (
	"net/http"

	"github.com/cerfical/merchshop/internal/domain/services"
	"github.com/cerfical/merchshop/internal/log"
)

func NewHandler(auth services.AuthService, log *log.Logger) http.Handler {
	mux := http.NewServeMux()

	a := authHandler{auth, log}
	mux.HandleFunc("POST /api/auth", a.authUser)

	return mux
}
