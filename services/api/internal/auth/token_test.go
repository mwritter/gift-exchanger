package auth

import (
	"regexp"
	"testing"
)

func TestHashTokenIsStableAndNotRaw(t *testing.T) {
	raw := "example-token"
	hash := HashToken(raw)
	if hash == raw {
		t.Fatal("hash must not equal the raw token")
	}
	if hash != HashToken(raw) {
		t.Fatal("hash must be deterministic")
	}
	if HashToken("other") == hash {
		t.Fatal("different tokens must not hash the same")
	}
}

func TestNewLoginCodeIsSixDigits(t *testing.T) {
	pattern := regexp.MustCompile(`^[0-9]{6}$`)
	for range 1000 {
		code, err := newLoginCode()
		if err != nil {
			t.Fatal(err)
		}
		if !pattern.MatchString(code) {
			t.Fatalf("code %q is not six digits", code)
		}
	}
}

func TestHashLoginCodeIsBoundToEmail(t *testing.T) {
	if HashLoginCode("a@example.com", "123456") == HashLoginCode("b@example.com", "123456") {
		t.Fatal("same code for different emails must not hash the same")
	}
}
