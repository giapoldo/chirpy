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
	// readinessPath  = "GET /api/healthz"
	// metricsPath    = "GET /admin/metrics"
	// resetPath      = "POST /admin/reset"
	// validatePath   = "POST /api/validate_chirp"
	// usersPath      = "POST /api/users"

	rootPath = "."
	port     = "8080"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
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

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening database: %w", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{}
	apiCfg.db = dbQueries
	apiCfg.platform = cfgPlatform

	serveMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(rootPath))
	fS := http.StripPrefix("/app", fileServer)
	serveMux.Handle(fileServerPath, apiCfg.middlewareMetricsInc(fS))

	registerEndpoints(serveMux, &apiCfg)
	// serveMux.HandleFunc(readinessPath, handlerReadiness)
	// serveMux.HandleFunc(metricsPath, apiCfg.handlerHits)
	// serveMux.HandleFunc(resetPath, apiCfg.handlerReset)
	// serveMux.HandleFunc(validatePath, handlerValidateChirp)
	// serveMux.HandleFunc(usersPath, apiCfg.handlerCreateUser)

	server := http.Server{}
	server.Addr = ":" + port
	server.Handler = serveMux

	server.ListenAndServe()
}
