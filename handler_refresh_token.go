package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/giapoldo/chirpy/internal/auth"
	"github.com/giapoldo/chirpy/internal/database"
)

type tokenResponse struct {
	Token string `json:"token"`
}

// POST /api/refresh
func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("Refresh_token handler:", err)
	}

	db_refreshToken, err := cfg.db.GetSingleRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not in database")
		return
	} else if db_refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token expired")
	}

	access_token, err := auth.MakeJWT(db_refreshToken.UserID, cfg.jwtSecret, time.Duration(1*time.Hour))
	if err != nil {
		log.Println("Refresh_token handler, make access token:", err)
		return
	}

	tokenResponse := tokenResponse{
		Token: access_token,
	}

	respondWithJSON(w, http.StatusOK, tokenResponse)
	return

}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("Revoke token handler, get token:", err)
		return
	}

	revoked_time := time.Now()
	err = cfg.db.RevokeSingleRefreshToken(r.Context(), database.RevokeSingleRefreshTokenParams{
		Token: refreshToken,
		RevokedAt: sql.NullTime{
			Time:  revoked_time,
			Valid: true,
		},
		UpdatedAt: revoked_time,
	})
	if err != nil {
		log.Println("Revoke token handler, update db:", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
	return

}
