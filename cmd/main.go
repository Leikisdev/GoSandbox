package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Leikisdev/GoSandbox/internal/database"
	"github.com/Leikisdev/GoSandbox/internal/web"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	godotenv.Load()
	const port = "8080"

	mux := http.NewServeMux()
	dbURL := os.Getenv("DB_URL")
	DB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to DB")
	}

	c := web.ApiConfig{
		FileserverHits: atomic.Int32{},
		DB:             database.New(DB),
		Platform:       os.Getenv("PLATFORM"),
	}

	fileHandler := http.StripPrefix("/app/", http.FileServer(web.Assets()))
	mux.Handle("/app/", c.CountingMiddleware(fileHandler))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", c.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", c.ResetHandler)
	mux.HandleFunc("POST /api/chirps", c.ChirpPostHandler)
	mux.HandleFunc("GET /api/chirps", c.ChirpGetAllHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", c.ChirpGetSingleHandler)
	mux.HandleFunc("POST /api/users", c.UserCreateHandler)
	mux.HandleFunc("POST /api/login", c.LoginHandler)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving started on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
