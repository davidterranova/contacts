package xhttp

import (
	"net/http"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/auth"
)

type AuthFn func(r *http.Request) (*domain.User, error)

func AuthMiddleware(authFn AuthFn) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, err := authFn(r)
			if err != nil {
				WriteError(ctx, w, http.StatusUnauthorized, "unauthorized", err)
				return
			}

			reqWithCtx := r.WithContext(auth.ContextWithUser(ctx, *user))
			next.ServeHTTP(w, reqWithCtx)
		})
	}
}

func BasicAuthFn(username string, password string) AuthFn {
	return func(r *http.Request) (*domain.User, error) {
		user, err := auth.BasicAuth(username, password)(r.Header.Get("Authorization"))
		if err != nil {
			return nil, err
		}

		return &user, nil
	}
}
