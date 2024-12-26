package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type parameters struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}
type returnVals struct {
	Error string `json:"error"`
}

type chirpData struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	retVals := returnVals{
		Error: msg,
	}
	data, err := json.Marshal(retVals)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}
