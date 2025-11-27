package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authString := headers.Get("Authorization")
	if authString == "" {
		return "", errors.New("missing authorization header")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authString, prefix) {
		return "", errors.New("invalid authorization header")
	}

	return strings.TrimPrefix(authString, prefix), nil
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("unable to generate refresh token")
	}

	return hex.EncodeToString(b), nil
}
