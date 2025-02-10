package api

import (
	"net/http"

	"github.com/cerfical/merchshop/internal/log"
	"github.com/cerfical/merchshop/internal/model"
)

func New(users model.UserStore, log *log.Logger) http.Handler {
	mux := http.NewServeMux()

	auth := authHandler{users, log}
	mux.HandleFunc("POST /api/auth", auth.authUser)

	return mux
}
