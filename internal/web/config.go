package web

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/Leikisdev/GoSandbox/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
	SigningSecret  string
}

func (c *ApiConfig) MetricsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, c.FileserverHits.Load())))
}

func (c *ApiConfig) ResetHandler(w http.ResponseWriter, req *http.Request) {
	if c.Platform == "dev" {
		c.FileserverHits.Store(0)
		if err := c.DB.DeleteUsers(req.Context()); err != nil {
			respondWithError(w, http.StatusBadRequest, "Unable to delete users")
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0, users deleted"))
	} else {
		respondWithError(w, http.StatusForbidden, "Cannot reset outside of dev")
	}

}
