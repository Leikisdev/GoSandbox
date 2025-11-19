package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTCreateAndVal(t *testing.T) {
	secret := "mysecret"
	userId, _ := uuid.NewRandom()
	jwt, err := MakeJWT(userId, secret, 5*time.Second)
	if err != nil {
		t.Errorf("failed to generate JWT: %v", err)
	}

	userIdJWT, err := ValidateJWT(jwt, secret)
	if err != nil || userIdJWT != userId {
		t.Errorf("failed to validate JWT: %v", err)
	}
}

func TestJWTValidationError(t *testing.T) {
	secret := "mysecret"
	wrongSecret := "mywrongsecret"
	userId, _ := uuid.NewRandom()
	jwt, err := MakeJWT(userId, wrongSecret, 5*time.Second)
	if err != nil {
		t.Errorf("failed to generate JWT: %v", err)
	}

	_, err2 := ValidateJWT(jwt, secret)
	if err2 == nil {
		t.Error("expected validation with wrong secret to fail")
	}
}

func TestExpiredJWT(t *testing.T) {
	secret := "mysecret"
	userId, _ := uuid.NewRandom()
	jwt, err := MakeJWT(userId, secret, 1*time.Second)
	if err != nil {
		t.Errorf("failed to generate JWT: %v", err)
	}

	time.Sleep(2 * time.Second)

	_, err2 := ValidateJWT(jwt, secret)
	if err2 == nil {
		t.Error("expected validation of expired token to fail")
	}
}
