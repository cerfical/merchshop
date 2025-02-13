package api

import (
	"net/http"

	"github.com/cerfical/merchshop/internal/domain/services"
	"github.com/cerfical/merchshop/internal/log"
)

type coinsHandler struct {
	coinService services.CoinService
	log         *log.Logger
}

func (h *coinsHandler) info(w http.ResponseWriter, r *http.Request) {
}
