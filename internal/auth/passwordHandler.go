package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(pass string) (string, error) {
	hashedPass, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("unable to hash password, ERR: %w", err)
	}
	return hashedPass, nil
}

func CompareHashedPass(pass, hashedPass string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(pass, hashedPass)
	if err != nil {
		return false, fmt.Errorf("unable to compare pass with hash, ERR: %w", err)
	}
	return match, nil
}
