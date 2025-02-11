package api

import (
	"net/http"

	"github.com/cerfical/merchshop/internal/log"
	"github.com/cerfical/merchshop/internal/model"
)

func New(auth *AuthConfig, users model.UserStore, log *log.Logger) http.Handler {
	mux := http.NewServeMux()

	a := authHandler{users, auth, log}
	mux.HandleFunc("POST /api/auth", a.authUser)

	return mux
}
