package api

import (
	"net/http"

	"github.com/cerfical/merchshop/internal/domain/coins"
	"github.com/cerfical/merchshop/internal/log"
)

type coinsHandler struct {
	coinService coins.CoinService
	log         *log.Logger
}

func (h *coinsHandler) info(w http.ResponseWriter, r *http.Request) {
	user := usernameFromContext(r.Context())
	coins, err := h.coinService.GetCoinBalance(user)
	if err != nil {
		return
	}

	writeResponse(w, http.StatusOK, &infoResponse{
		Coins: int(coins),
	})
}

type infoResponse struct {
	Coins int `json:"coins"`
}
