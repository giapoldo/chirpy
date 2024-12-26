package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

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
