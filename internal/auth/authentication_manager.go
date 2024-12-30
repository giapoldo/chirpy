package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var regClaim *jwt.RegisteredClaims

func HashPassword(password string) (string, error) {
	hashed_pwd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Printf("Couldn't hash password: %s\n", err)
		return "", err
	}
	return string(hashed_pwd), nil

}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("Wrong password: %s\n", err)
		return err
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	regClaim = &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, regClaim)

	tokenStr, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Println("Couldn't sign the JWT token", err)
		return "", err
	}
	return tokenStr, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	log.Println("DEBUG:", tokenString)
	log.Println()
	log.Println("DEBUG:", regClaim)

	token, err := jwt.ParseWithClaims(tokenString, regClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		log.Println("Parse with claims error:", err)
		return uuid.UUID{}, err
	}
	uuidStr, err := token.Claims.GetSubject()
	if err != nil {
		log.Println("Couldn't parse subject to uuid:", err)
		return uuid.UUID{}, err
	}
	return uuid.MustParse(uuidStr), nil

}

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return "", fmt.Errorf("Couldn't retrieve token from header")
	}
	return strings.Fields(bearer)[1], nil
}

func MakeRefreshToken() (string, error) {

	c := 32
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	return hex.EncodeToString(b), nil

}

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("Couldn't retrieve api key from header")
	}
	return strings.Fields(apiKey)[1], nil
}
