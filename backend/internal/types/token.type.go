package types

import "github.com/golang-jwt/jwt/v5"

type TokenPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type CustomClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}
