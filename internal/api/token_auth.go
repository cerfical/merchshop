package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/cerfical/merchshop/internal/domain/auth"
	"github.com/cerfical/merchshop/internal/domain/model"
)

func tokenAuth(a auth.AuthService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Validate the Authorization header
			s := strings.Split(r.Header.Get("Authorization"), " ")
			if len(s) != 2 || s[0] != "Bearer" {
				unauthorized(w, "The provided authentication method is invalid")
				return
			}

			user, err := a.AuthToken(auth.Token(s[1]))
			if err != nil {
				unauthorized(w, "The provided token is invalid")
				return
			}

			rr := r.WithContext(contextWithUsername(r.Context(), user))
			next(w, rr)
		}
	}
}

func unauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	writeErrorResponse(w, http.StatusUnauthorized, msg)
}

func contextWithUsername(ctx context.Context, un model.Username) context.Context {
	return context.WithValue(ctx, userContextKey{}, un)
}

func usernameFromContext(ctx context.Context) model.Username {
	return ctx.Value(userContextKey{}).(model.Username)
}

type userContextKey struct{}
