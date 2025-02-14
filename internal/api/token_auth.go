package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/cerfical/merchshop/internal/service/auth"
	"github.com/cerfical/merchshop/internal/service/model"
)

func tokenAuth(a auth.AuthService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Validate the Authorization header
			s := strings.Split(r.Header.Get("Authorization"), " ")
			if len(s) != 2 || s[0] != "Bearer" {
				unauthorizedHandler("The provided authentication method is invalid")(w, r)
				return
			}

			user, err := a.AuthToken(auth.Token(s[1]))
			if err != nil {
				unauthorizedHandler("The provided token is invalid")(w, r)
				return
			}

			rr := r.WithContext(contextWithUsername(r.Context(), user))
			next(w, rr)
		}
	}
}

func contextWithUsername(ctx context.Context, un model.Username) context.Context {
	return context.WithValue(ctx, userContextKey{}, un)
}

func usernameFromContext(ctx context.Context) model.Username {
	return ctx.Value(userContextKey{}).(model.Username)
}

type userContextKey struct{}
