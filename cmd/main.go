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

func main() {
	godotenv.Load()
	const port = "8080"

	mux := http.NewServeMux()
	dbURL := os.Getenv("DB_URL")
	log.Print(dbURL)
	DB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to DB")
	}

	c := web.ApiConfig{
		FileserverHits: atomic.Int32{},
		DB:             database.New(DB),
		Platform:       os.Getenv("PLATFORM"),
		SigningSecret:  os.Getenv("SECRET"),
	}

	c.RegisterRoutes(mux)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving started on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
