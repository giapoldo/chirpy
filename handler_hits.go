package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerHits(resWriter http.ResponseWriter, req *http.Request) {

	response := fmt.Sprintf("<html><body>\n<h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())

	// response := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())

	resWriter.Header().Set("Content-Type", "text/html")
	resWriter.WriteHeader(http.StatusOK)
	resWriter.Write([]byte(response))

}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if cfg.platform == "dev" {
		cfg.fileserverHits.Store(0)

		err := cfg.db.DeleteAllUsers(r.Context())
		if err != nil {
			log.Fatal("Unable to reset the database: %w", err)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server hits and database reset"))
	} else {
		// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		// w.Write([]byte("Server hits reset"))
	}

}
