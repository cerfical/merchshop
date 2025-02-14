package api

import (
	"net/http"

	"github.com/cerfical/merchshop/internal/log"
	"github.com/cerfical/merchshop/internal/service/auth"
	"github.com/cerfical/merchshop/internal/service/coins"
)

func NewHandler(auth auth.AuthService, coins coins.CoinService, log *log.Logger) http.Handler {
	mux := http.NewServeMux()

	a := authHandler{auth, log}
	mux.HandleFunc("POST /api/auth", a.authUser)

	c := coinsHandler{coins, log}
	mux.HandleFunc("GET /api/info", tokenAuth(auth)(c.info))
	mux.HandleFunc("POST /api/sendCoin", tokenAuth(auth)(c.sendCoin))

	return mux
}
