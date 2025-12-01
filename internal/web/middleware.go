package web

import (
	"context"
	"net/http"

	"github.com/Leikisdev/GoSandbox/internal/auth"
)

type ctxUserIdKey struct{}

func (c *ApiConfig) CountingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *ApiConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userId, err := auth.ValidateJWT(token, c.SigningSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserIdKey{}, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
