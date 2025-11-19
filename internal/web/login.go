package web

import (
	"encoding/json"
	"net/http"

	"github.com/Leikisdev/GoSandbox/internal/auth"
)

func (c *ApiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var params UserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to parse request")
		return
	}

	user, err := c.DB.LoginUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	match, err := auth.CompareHashedPass(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	respondWithJSON(w, http.StatusOK, toUser(user))
}
