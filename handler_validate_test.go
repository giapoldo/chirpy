package main

import "testing"

func TestFilterBadWords(t *testing.T) {
	// let's write our first test case
	input := "This is a kerfuffle! string"
	want := "This is a kerfuffle! string"
	got := filterBadWords(input)
	if got != want {
		t.Errorf("cleanChirp(%q) = %q, want %q", input, got, want)
	}
}
