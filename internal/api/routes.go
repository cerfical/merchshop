package api

import (
	"net/http"

	"github.com/cerfical/merchshop/internal/domain/auth"
	"github.com/cerfical/merchshop/internal/log"
)

func New(authService auth.Service, log *log.Logger) http.Handler {
	mux := http.NewServeMux()

	auth := authHandler{authService, log}
	mux.HandleFunc("POST /api/auth", auth.authUser)

	return mux
}
