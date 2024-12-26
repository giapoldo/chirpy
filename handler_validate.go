package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if l := len(params.Body); l > 140 {

		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Chirp longer than 140 characters (%v)", l))
		return
	}

	retVals := returnVals{
		Error:       "",
		CleanedBody: filterBadWords(params.Body),
	}

	respondWithJSON(w, http.StatusOK, retVals)
	return

}
