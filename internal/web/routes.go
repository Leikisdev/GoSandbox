package web

import (
	"net/http"
)

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (c *ApiConfig) RegisterRoutes(mux *http.ServeMux) {
	// Public routes (no auth)
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("POST /api/login", c.LoginHandler)
	mux.HandleFunc("GET /api/chirps", c.ChirpGetAllHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", c.ChirpGetSingleHandler)
	mux.HandleFunc("POST /api/refresh", c.RefreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", c.RevokeTokenHandler)
	mux.HandleFunc("POST /api/users", c.UserCreateHandler)

	// Protected routes (require auth)
	protectedChirpPost := c.AuthMiddleware(http.HandlerFunc(c.ChirpPostHandler))
	mux.Handle("POST /api/chirps", protectedChirpPost)
	protectedUserUpdate := c.AuthMiddleware(http.HandlerFunc(c.UserUpdateHandler))
	mux.Handle("PUT /api/users", protectedUserUpdate)

	// Admin routes
	mux.HandleFunc("/admin/metrics", c.MetricsHandler)
	mux.HandleFunc("/admin/reset", c.ResetHandler)

	// Static files
	fileHandler := http.StripPrefix("/app/", http.FileServer(Assets()))
	mux.Handle("/app/", c.CountingMiddleware(fileHandler))
}
