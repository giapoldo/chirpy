package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/giapoldo/chirpy/internal/auth"
	"github.com/giapoldo/chirpy/internal/database"
	"github.com/google/uuid"
)

type user struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAtAt  time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type userParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type redWebhookParams struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hashed_pwd, err := auth.HashPassword(params.Password)

	createdUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_pwd,
	})
	if err != nil {
		log.Printf("Couldn't create user : %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	user := user{
		ID:          createdUser.ID,
		CreatedAt:   createdUser.CreatedAt,
		UpdatedAtAt: createdUser.UpdatedAt,
		Email:       createdUser.Email,
		IsChirpyRed: createdUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, user)
	return

}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Println(err)
	}

	foundUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Couldn't find user : %s", err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user")
		return
	}

	db_refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Duration(60 * 24 * time.Hour)),
		UserID:    foundUser.ID,
	})
	if err != nil {
		log.Println("Insert refresh token:", err)
	}

	err = auth.CheckPasswordHash(params.Password, foundUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong password")
	}

	expiration := time.Duration(1 * time.Hour)

	JWTString, err := auth.MakeJWT(foundUser.ID, cfg.jwtSecret, expiration)
	if err != nil {
		log.Println(err)
	}

	user := user{
		ID:           foundUser.ID,
		CreatedAt:    foundUser.CreatedAt,
		UpdatedAtAt:  foundUser.UpdatedAt,
		Email:        foundUser.Email,
		Token:        JWTString,
		RefreshToken: db_refreshToken.Token,
		IsChirpyRed:  foundUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, user)
	return

}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("UpdateUser, accestoken:", err)
		respondWithError(w, http.StatusUnauthorized, "No token in request")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err = decoder.Decode(&params)

	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusUnauthorized, "No token in request")
		return
	}

	jwt_user, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Println("UpdateUser, validatejwt:", err)
		respondWithError(w, http.StatusUnauthorized, "Malformed token")
		return
	}

	hashed_pwd, err := auth.HashPassword(params.Password)

	updated_user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_pwd,
		ID:             jwt_user,
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		log.Printf("Couldn't update user : %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	user := user{
		ID:          updated_user.ID,
		CreatedAt:   updated_user.CreatedAt,
		UpdatedAtAt: updated_user.UpdatedAt,
		Email:       updated_user.Email,
		IsChirpyRed: updated_user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, user)
	return

}

func (cfg *apiConfig) handlerUpgradeToRedWebhook(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No apiKey")
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Wrong apiKey")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := redWebhookParams{}
	err = decoder.Decode(&params)

	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "")
		return
	}

	_, err = cfg.db.UpgradeUserToRed(r.Context(), uuid.MustParse(params.Data.UserID))
	if err != nil {
		log.Printf("Couldn't find user : %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't find user")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
	return

}
