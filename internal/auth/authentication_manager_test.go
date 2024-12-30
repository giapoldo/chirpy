package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	// let's write our first test case
	input := "A_VeryStrongPa$$word!"
	// want := nil
	hash, _ := HashPassword(input)
	got := CheckPasswordHash(input, hash)

	if got != nil {
		t.Errorf("Password: %q, Hash check: %q", input, got)
	}
}

func TestJWT(t *testing.T) {

	userID := uuid.MustParse("dc57d778-b265-479f-94a4-6bd893d458a4")
	secret := "12345"
	expiration := time.Duration(time.Second)

	JWTString, err := MakeJWT(userID, secret, expiration)

	if err != nil {
		t.Errorf("JWT: %q,", JWTString)
	}
	// t.Log(JWTString)

	user, err := ValidateJWT(JWTString, secret)

	if err != nil {
		t.Errorf("JWT validation failed before expiration: %q,", user)
	}
	if user != userID {
		t.Errorf("wrong user from token")
	}
	time.Sleep(2 * time.Second)

	user, err = ValidateJWT(JWTString, secret)
	if err != nil {
		t.Logf("JWT validation failed after expiration: %q,", user)
	}

}

func TestGetBearer(t *testing.T) {
	// let's write our first test case
	h := http.Header{}
	h.Add("Authorization", "Bearer askgsdnfksdglkfgm293450823lklkdnldnsdgsdgm,mndgs")
	want := "askgsdnfksdglkfgm293450823lklkdnldnsdgsdgm,mndgs"
	got, _ := GetBearerToken(h)

	if got != want {
		t.Errorf("input: %q, want: %q, got: %q", h, want, got)
	}
}
