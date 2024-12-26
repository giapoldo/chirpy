package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type user struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAtAt time.Time `json:"updated_at"`
	Email       string    `json:"email"`
}

type email struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := email{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createdUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Fatal("Couldn't create user : %w", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	user := user{
		ID:          createdUser.ID,
		CreatedAt:   createdUser.CreatedAt,
		UpdatedAtAt: createdUser.UpdatedAt,
		Email:       createdUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)
	return

}
