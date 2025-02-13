package api

import (
	"net/http"

	"github.com/cerfical/merchshop/internal/domain/services"
	"github.com/cerfical/merchshop/internal/log"
)

func NewHandler(auth services.AuthService, coins services.CoinService, log *log.Logger) http.Handler {
	mux := http.NewServeMux()

	a := authHandler{auth, log}
	mux.HandleFunc("POST /api/auth", a.authUser)

	c := coinsHandler{coins, log}
	mux.HandleFunc("GET /api/info", c.info)

	return mux
}
