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

			next(w, requestWithUser(r, user))
		}
	}
}

func requestWithUser(r *http.Request, un model.Username) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), userContextKey{}, un))
}

func userFromRequest(r *http.Request) model.Username {
	return r.Context().Value(userContextKey{}).(model.Username)
}

type userContextKey struct{}
