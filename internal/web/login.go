package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Leikisdev/GoSandbox/internal/auth"
	"github.com/Leikisdev/GoSandbox/internal/database"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

func (c *ApiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var params LoginRequest
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

	token, err := auth.MakeJWT(user.ID, c.SigningSecret, 3600*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	refresh_tk, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}

	if _, err := c.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresh_tk,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(1440 * time.Hour),
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to register refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh_tk,
	})
}

func (c *ApiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refresh_tk, err := c.DB.GetRefreshToken(r.Context(), tok)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	new_tok, err := auth.MakeJWT(refresh_tk.UserID, c.SigningSecret, 3600*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate new access token")
		return
	}

	respondWithJSON(w, http.StatusOK, RefreshTokenResponse{Token: new_tok})
}

func (c *ApiConfig) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if _, err := c.DB.RevokeToken(r.Context(), tok); err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to revoke refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
