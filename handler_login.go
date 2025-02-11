package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Chance093/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	pl := payload{}
	if err := json.NewDecoder(r.Body).Decode(&pl); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while decoding payload", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), pl.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error while retrieving user", err)
		return
	}

	if err = auth.CheckPasswordHash(pl.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", nil)
    return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
