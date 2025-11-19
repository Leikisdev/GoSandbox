package auth

import (
	"testing"
)

// TestHashPassword tests the HashPassword function.
func TestHashPassword(t *testing.T) {
	pass := "mysecretpassword"
	if _, err := HashPassword(pass); err != nil {
		t.Errorf("failed to hash password: %v", err)
	}
}

// TestCompareHashedPass tests the CompareHashedPass function.
func TestCompareHashedPass(t *testing.T) {
	pass := "mysecretpassword"
	hash, err := HashPassword(pass)
	if err != nil {
		t.Errorf("failed to hash password: %v", err)
	}
	match, err := CompareHashedPass(pass, hash)
	if !match || err != nil {
		t.Errorf("failed to compare password and hash: %v", err)
	}
}
