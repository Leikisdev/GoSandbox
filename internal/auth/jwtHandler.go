package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	signedTok, err := tok.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token, err: %v", err)
	}
	return signedTok, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	verifiedTok, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to verify token")
	}

	userId, err := verifiedTok.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to extract userId from JWT")
	}
	return uuid.Parse(userId)
}
