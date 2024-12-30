package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/giapoldo/chirpy/internal/auth"
	"github.com/giapoldo/chirpy/internal/database"
	"github.com/google/uuid"
)

var badWords []string = []string{"kerfuffle", "sharbert", "fornax"}

func filterBadWords(chirp string) (cleanChirp string) {
	split_chirp := strings.Split(chirp, " ")

	for i, word := range split_chirp {
		for _, bword := range badWords {
			if strings.ToLower(word) == strings.ToLower(bword) {
				split_chirp[i] = "****"
			}
		}
	}
	cleanChirp = strings.Join(split_chirp, " ")
	return
}

// POST /api/chirps
func (cfg *apiConfig) handlerAddChirps(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("No token in request 1", err)
		respondWithError(w, http.StatusUnauthorized, "No token in request 1")
		return
	}

	JWTuserID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Println("No token in request 2", err)
		respondWithError(w, http.StatusUnauthorized, "No token in request 2")
		return
	}

	if l := len(params.Body); l > 140 {

		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Chirp longer than 140 characters (%v)", l))
		return
	}

	cleanedBody := filterBadWords(params.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: JWTuserID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s\n", err)
		return
	}

	new_chirp := chirpData{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}

	respondWithJSON(w, http.StatusCreated, new_chirp)
	return

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	db_chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve Chirps")

	}

	retrieved_chirps := []chirpData{}
	for _, chirp := range db_chirps {
		retrieved_chirps = append(retrieved_chirps, chirpData{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.String(),
			UpdatedAt: chirp.UpdatedAt.String(),
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		})
	}

	respondWithJSON(w, http.StatusOK, retrieved_chirps)
	return
}

func (cfg *apiConfig) handlerGetSingletonChirp(w http.ResponseWriter, r *http.Request) {

	chirpID := uuid.MustParse(r.PathValue("chirpID"))

	chirp, err := cfg.db.GetSingleChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find Chirp")

	}

	retrieved_chirp := chirpData{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}

	respondWithJSON(w, http.StatusOK, retrieved_chirp)
	return
}

func (cfg *apiConfig) handlerDeleteSingletonChirp(w http.ResponseWriter, r *http.Request) {

	chirpID := uuid.MustParse(r.PathValue("chirpID"))

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("UpdateUser, accestoken:", err)
		respondWithError(w, http.StatusUnauthorized, "No token in request")
		return
	}

	jwt_user, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Println("UpdateUser, validatejwt:", err)
		respondWithError(w, http.StatusUnauthorized, "Malformed token")
		return
	}

	chirp, err := cfg.db.GetSingleChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find Chirp")
		return
	}

	if jwt_user != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "User mismatch")
		return
	}

	cfg.db.DeleteSingletonChirps(r.Context(), chirp.ID)

	respondWithJSON(w, http.StatusNoContent, "")
	return
}
