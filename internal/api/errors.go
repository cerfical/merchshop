package api

import (
	"net/http"

	"github.com/cerfical/merchshop/internal/log"
)

func internalError(msg string, err error, log *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Error(msg, err)

		encode(w, http.StatusInternalServerError, errorResponse{
			Errors: "Internal server error",
		})
	}
}
