package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/giapoldo/chirpy/internal/auth"
	"github.com/giapoldo/chirpy/internal/database"
	"github.com/google/uuid"
)

type user struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAtAt time.Time `json:"updated_at"`
	Email       string    `json:"email"`
}

type userParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hashed_pwd, err := auth.HashPassword(params.Password)

	createdUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_pwd,
	})
	if err != nil {
		log.Printf("Couldn't create user : %s", err)
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

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	foundUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Couldn't find user : %s", err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user")
		return
	}

	if err != nil {
		log.Printf("Hashing error %s", err)
	}
	err = auth.CheckPasswordHash(params.Password, foundUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong password")
	}
	user := user{
		ID:          foundUser.ID,
		CreatedAt:   foundUser.CreatedAt,
		UpdatedAtAt: foundUser.UpdatedAt,
		Email:       foundUser.Email,
	}

	respondWithJSON(w, http.StatusOK, user)
	return

}
