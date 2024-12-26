package main

import (
	"net/http"
)

const (
	readinessPath = "GET /api/healthz"
	metricsPath   = "GET /admin/metrics"
	resetPath     = "POST /admin/reset"
	validatePath  = "POST /api/validate_chirp"
	usersPath     = "POST /api/users"
)

type endpoint struct {
	Name    string
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

type endpoints []endpoint

func registerEndpoints(mux *http.ServeMux, apiCfg *apiConfig) {

	endpoints := endpoints{
		endpoint{
			Name:    "readiness",
			Path:    readinessPath,
			Handler: handlerReadiness,
		},
		endpoint{
			Name:    "hits",
			Path:    metricsPath,
			Handler: apiCfg.handlerHits,
		},
		endpoint{
			Name:    "reset",
			Path:    resetPath,
			Handler: apiCfg.handlerReset,
		},
		endpoint{
			Name:    "validate_chirp",
			Path:    validatePath,
			Handler: handlerValidateChirp,
		},
		endpoint{
			Name:    "create_user",
			Path:    usersPath,
			Handler: apiCfg.handlerCreateUser,
		},
	}

	for _, endpoint := range endpoints {
		mux.HandleFunc(endpoint.Path, endpoint.Handler)
	}
}
