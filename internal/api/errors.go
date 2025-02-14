package api

import (
	"net/http"

	"github.com/cerfical/merchshop/internal/log"
)

func badRequestHandler(errors string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeErrorResponse(w, http.StatusBadRequest, errors)
	}
}

func internalErrorHandler(log *log.Logger, msg string, err error) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Error(msg, err)
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server failure")
	}
}

func unauthorizedHandler(errors string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("WWW-Authenticate", "Bearer")
		writeErrorResponse(w, http.StatusUnauthorized, errors)
	}
}

func writeErrorResponse(w http.ResponseWriter, status int, errors string) {
	writeResponse(w, status, errorResponse{
		Errors: errors,
	})
}

type errorResponse struct {
	Errors string `json:"errors"`
}
