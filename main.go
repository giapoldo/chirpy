package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/giapoldo/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	fileServerPath = "/app/"
	rootPath       = "."
	port           = "8080"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	cfgPlatform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening database: %w", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{}
	apiCfg.db = dbQueries
	apiCfg.platform = cfgPlatform
	apiCfg.jwtSecret = jwtSecret
	apiCfg.polkaKey = polkaKey

	serveMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(rootPath))
	fS := http.StripPrefix("/app", fileServer)
	serveMux.Handle(fileServerPath, apiCfg.middlewareMetricsInc(fS))

	registerEndpoints(serveMux, &apiCfg)

	server := http.Server{}
	server.Addr = ":" + port
	server.Handler = serveMux

	server.ListenAndServe()
}
