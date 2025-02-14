package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cerfical/merchshop/internal/log"
	"github.com/cerfical/merchshop/internal/service/coins"
	"github.com/cerfical/merchshop/internal/service/model"
)

type coinsHandler struct {
	coinService coins.CoinService
	log         *log.Logger
}

func (h *coinsHandler) info(w http.ResponseWriter, r *http.Request) {
	// TODO: The assumption is that the username provided refers to an existing user
	user := userFromRequest(r)
	coins, err := h.coinService.GetCoinBalance(user)
	if err != nil {
		return
	}

	writeResponse(w, http.StatusOK, &infoResponse{
		Coins: int(coins),
	})
}

func (h *coinsHandler) sendCoin(w http.ResponseWriter, r *http.Request) {
	sendCoinReq, err := readRequest[sendCoinRequest](r.Body)
	if err != nil {
		// TODO: More descriptive error messages?
		badRequestHandler("The request body is malformed")(w, r)
		return
	}

	toUser, err := model.NewUsername(sendCoinReq.ToUser)
	if err != nil {
		badRequestHandler(fmt.Sprintf("The recipient is invalid: %v", err))(w, r)
		return
	}

	amount, err := model.NewNumCoins(sendCoinReq.Amount)
	if err != nil {
		badRequestHandler(fmt.Sprintf("Invalid amount of coins to transfer: %v", err))(w, r)
		return
	}

	fromUser := userFromRequest(r)
	if err := h.coinService.SendCoins(fromUser, toUser, amount); err != nil {
		if modelErr := model.Error(""); errors.As(err, &modelErr) {
			badRequestHandler(fmt.Sprintf("Couldn't complete the coin transfer: %v", modelErr))(w, r)
		} else {
			internalErrorHandler(h.log, "Failed to perform coin transfer", err)(w, r)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *coinsHandler) buyItem(w http.ResponseWriter, r *http.Request) {
	merch, err := model.NewMerchItem(r.PathValue("item"))
	if err != nil {
		badRequestHandler(fmt.Sprintf("The requested merch is unavailable: %v", err))(w, r)
		return
	}

	buyer := userFromRequest(r)
	if err := h.coinService.BuyItem(buyer, merch); err != nil {
		if modelErr := model.Error(""); errors.As(err, &modelErr) {
			badRequestHandler(fmt.Sprintf("Couldn't complete the merch purchase: %v", modelErr))(w, r)
		} else {
			internalErrorHandler(h.log, "Failed to perform merch purchase", err)(w, r)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

type sendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type infoResponse struct {
	Coins int `json:"coins"`
}
