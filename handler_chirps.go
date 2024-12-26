package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

	if l := len(params.Body); l > 140 {

		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Chirp longer than 140 characters (%v)", l))
		return
	}

	cleanedBody := filterBadWords(params.Body)
	// retVals := returnVals{
	// 	Error:       "",
	// 	CleanedBody: filterBadWords(params.Body),
	// }

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: uuid.MustParse(params.UserID),
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
