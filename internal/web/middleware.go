package web

import (
	"net/http"

	"github.com/Leikisdev/GoSandbox/internal/auth"
)

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

		if _, err := auth.ValidateJWT(token, c.SigningSecret); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
