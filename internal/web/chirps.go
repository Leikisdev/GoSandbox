package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Leikisdev/GoSandbox/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (c *ApiConfig) ChirpPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload Chirp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Printf("decode error: %v", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	} else if len(payload.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := c.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   payload.Body,
		UserID: r.Context().Value(ctxUserIdKey{}).(uuid.UUID),
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Unable to create user, ERR: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (c *ApiConfig) ChirpGetAllHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := c.DB.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Unable to fetch chirps, ERR: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (c *ApiConfig) ChirpGetSingleHandler(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
	}

	chirp, err := c.DB.GetChirpById(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Unable to fetch chirp, ERR: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (c *ApiConfig) ChirpDeleteHandler(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID")
	}

	chirp, err := c.DB.GetChirpById(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to find chirp")
		return
	}

	userId := r.Context().Value(ctxUserIdKey{}).(uuid.UUID)
	if userId != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "cannot delete other chirps of other users")
		return
	}

	if err := c.DB.DeleteChirpById(r.Context(), chirpId); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
