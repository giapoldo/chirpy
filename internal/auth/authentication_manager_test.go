package auth

import (
	"testing"
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
