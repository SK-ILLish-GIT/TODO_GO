package utils

import (
	"backend/config"
	"backend/internal/types"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken is a helper function to create JWT tokens
func GenerateToken(id, email, secretKey string, expiry time.Duration) (string, error) {
	if secretKey == "" {
		return "", fmt.Errorf("secret key is not set")
	}

	claims := types.CustomClaims{
		ID:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(id, email string) (string, error) {
	refreshSecretKey := config.GetEnv("REFRESH_TOKEN_SECRET")
	return GenerateToken(id, email, refreshSecretKey, 24*time.Hour)
}

// GenerateAccessToken generates an access token
func GenerateAccessToken(id, email string) (string, error) {
	accessSecretKey := config.GetEnv("ACCESS_TOKEN_SECRET")
	return GenerateToken(id, email, accessSecretKey, 1*time.Hour)
}
