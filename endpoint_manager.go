package main

import (
	"net/http"
)

const (
	readinessPath         = "GET /api/healthz"
	metricsPath           = "GET /admin/metrics"
	resetPath             = "POST /admin/reset"
	addchirpsPath         = "POST /api/chirps"
	getchirpsPath         = "GET /api/chirps"
	getSingletonChirpPath = "GET /api/chirps/{chirpID}"
	usersPath             = "POST /api/users"
	updateUsersPath       = "PUT /api/users"
	loginPath             = "POST /api/login"
	refreshPath           = "POST /api/refresh"
	revokePath            = "POST /api/revoke"
	deletePath            = "DELETE /api/chirps/{chirpID}"
	redWebhookPath        = "POST /api/polka/webhooks"
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
			Name:    "addChirps",
			Path:    addchirpsPath,
			Handler: apiCfg.handlerAddChirps,
		},
		endpoint{
			Name:    "getChirps",
			Path:    getchirpsPath,
			Handler: apiCfg.handlerGetChirps,
		},
		endpoint{
			Name:    "getSingletonChirp",
			Path:    getSingletonChirpPath,
			Handler: apiCfg.handlerGetSingletonChirp,
		},
		endpoint{
			Name:    "createUser",
			Path:    usersPath,
			Handler: apiCfg.handlerCreateUser,
		},
		endpoint{
			Name:    "loginUser",
			Path:    loginPath,
			Handler: apiCfg.handlerLoginUser,
		},
		endpoint{
			Name:    "refreshToken",
			Path:    refreshPath,
			Handler: apiCfg.handlerRefreshToken,
		},
		endpoint{
			Name:    "revokeToken",
			Path:    revokePath,
			Handler: apiCfg.handlerRevokeToken,
		},
		endpoint{
			Name:    "updateUser",
			Path:    updateUsersPath,
			Handler: apiCfg.handlerUpdateUser,
		},
		endpoint{
			Name:    "deleteChirp",
			Path:    deletePath,
			Handler: apiCfg.handlerDeleteSingletonChirp,
		},
		endpoint{
			Name:    "updateToRed",
			Path:    redWebhookPath,
			Handler: apiCfg.handlerUpgradeToRedWebhook,
		},
	}

	for _, endpoint := range endpoints {
		mux.HandleFunc(endpoint.Path, endpoint.Handler)
	}
}
