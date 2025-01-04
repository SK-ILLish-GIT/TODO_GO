package middlewares

import (
	"backend/config"
	"backend/internal/types"
	"backend/pkg/utils"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("accessToken")
		log.Println("AccessToken : ", accessToken)
		if accessToken == "" {
			authHeader := c.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				accessToken = authHeader[7:]
			}
		}
		if accessToken == "" {
			return utils.JSONResponse(c, fiber.StatusUnauthorized, nil, "Accesstoken is missing")
		}
		accessSecretKey := config.GetEnv("ACCESS_TOKEN_SECRET")
		claims := &types.CustomClaims{}
		token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(accessSecretKey), nil
		})

		// Handle parsing or validation errors
		if err != nil || !token.Valid {
			return utils.JSONResponse(c, fiber.StatusUnauthorized, nil, "Invalid or expired access token")
		}

		// Token is valid; attach claims to the context
		c.Locals("userID", claims.ID)
		c.Locals("email", claims.Email)

		// Call the next middleware or handler
		return c.Next()
	}
}
