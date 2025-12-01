package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Leikisdev/GoSandbox/internal/auth"
	"github.com/Leikisdev/GoSandbox/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func toUser(u database.User) User {
	return User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
}

func (c *ApiConfig) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	var params UserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to parse request")
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to hash password")
		return
	}

	user, err := c.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("unable to create user, ERR: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, toUser(user))
}

func (c *ApiConfig) UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var params UserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to parse request")
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to hash password")
		return
	}

	user, err := c.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             r.Context().Value(ctxUserIdKey{}).(uuid.UUID),
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("unable to update user, ERR: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, toUser(user))
}
