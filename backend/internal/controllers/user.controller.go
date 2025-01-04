package controllers

import (
	"backend/config"
	"backend/constants"
	"backend/internal/models"
	"backend/internal/types"
	"backend/pkg/utils"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUser(c *fiber.Ctx) error {

	// get user email,id from context
	// make the ID HEX format
	// get collection
	// get user details from db
	// return hiding sensetive data
	// log.Println(c.Locals("email"))
	// log.Println(c.Locals("userID"))
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, nil, "Invalid or missing user ID")
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, nil, "Invalid user ID format")
	}
	collection := config.GetCollection(constants.USER_COLLECTION)

	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusNotFound, nil, "User not found")
	}

	data := fiber.Map{
		"userID":    user.ID.Hex(),
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	}
	return utils.JSONResponse(c, fiber.StatusOK, data, "User details retrived successfully")
}

func DeleteUser(c *fiber.Ctx) error {
	// get userID from context locals
	// make it premitive from hex
	// getCollection
	// delete the user from
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, nil, "Invalid or missing user ID")
	}
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, nil, "Invalid user ID format")
	}

	collection := config.GetCollection(constants.USER_COLLECTION)

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil || result.DeletedCount == 0 {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, nil, "User not found")
	}

	setCookie(c, "accessToken", "", -time.Hour)
	setCookie(c, "refreshToken", "", -time.Hour)

	return utils.JSONResponse(c, fiber.StatusOK, nil, "User deleted successfully")
}

func RefreshTokens(c *fiber.Ctx) error {
	//get refreshtoken
	//parse and validate refresh token
	//generate new accesstoken and refresh token
	// save it in db
	//save it in
	//
	// Retrieve the refresh token from cookies
	refreshToken := c.Cookies("refreshToken")

	// If no refresh token is found, return unauthorized
	if refreshToken == "" {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, nil, "Refresh token is missing")
	}

	// Parse and validate the refresh token
	refreshSecretKey := config.GetEnv("REFRESH_TOKEN_SECRET")
	claims := &types.CustomClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(refreshSecretKey), nil
	})

	// Handle token validation errors
	if err != nil || !token.Valid {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, nil, "Invalid or expired refresh token")
	}

	// Generate a new access token
	accessToken, err := utils.GenerateAccessToken(claims.ID, claims.Email)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, nil, "Failed to generate access token")
	}

	// Set the new access token in cookies
	setCookie(c, "accessToken", accessToken, 1*time.Hour)

	// Return success response with the new access token
	data := fiber.Map{
		"accessToken": accessToken,
	}
	return utils.JSONResponse(c, fiber.StatusOK, data, "Access token regenerated successfully")
}
