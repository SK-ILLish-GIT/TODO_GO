package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(pass string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password hash")
	}
	return string(passwordHash), nil
}

func ComparePassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, fmt.Errorf("failed to match passwords")
	}
	return true, nil
}
