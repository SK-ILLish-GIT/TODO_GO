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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setCookie(c *fiber.Ctx, name, value string, expiry time.Duration) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		HTTPOnly: true, // Prevent access via JavaScript for security
		// Secure:   true, 					 // Use only over HTTPS
		// SameSite: "Strict",               // Enforce strict same-site policy
		Expires: time.Now().Add(expiry), // Set the expiration time
	})
	// log.Println("Cookie Set:", name, value)
}
func generateTokens(userID, email string) (string, string, error) {
	refreshToken, err := utils.GenerateRefreshToken(userID, email)
	if err != nil {
		return "", "", err
	}

	accessToken, err := utils.GenerateAccessToken(userID, email)
	if err != nil {
		return "", "", err
	}
	return refreshToken, accessToken, nil
}
func RegisterUser(c *fiber.Ctx) error {

	// get deatils from body
	// parse those deatils to Go struct format and validate
	// check username already exists or not
	// generate hashed password
	//generate refresh token
	//generate access token
	// add access token to cookies

	collection := config.GetCollection(constants.USER_COLLECTION)
	var user models.User

	// Parse and validate user input
	if err := c.BodyParser(&user); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, nil, "Invalid Input")
	}

	// Check for username, email, password
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return utils.JSONResponse(c, fiber.StatusBadRequest, nil, "Username/Email/Password is missing")
	}

	// Check if the user already exists
	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return utils.JSONResponse(c, fiber.StatusConflict, nil, "User already exists")
	}

	// Generate hashed password
	hashedPassword, err := utils.GeneratePassword(user.Password)
	if err != nil {
		fmt.Println(err)
		return utils.JSONResponse(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	// Generate refresh and access tokens
	userID := primitive.NewObjectID()
	refreshToken, accessToken, err := generateTokens(userID.Hex(), user.Email)
	if err != nil {
		utils.JSONResponse(c, fiber.StatusInternalServerError, nil, "Error generating Refresh/Access token")
	}

	// Prepare user data for insertion
	user.ID = userID
	user.Password = hashedPassword
	user.CreatedAt = time.Now().Unix()
	user.RefreshToken = refreshToken

	// Insert the user into the database
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, nil, "Error occured while creating user")
	}

	// Set cookies
	setCookie(c, "accessToken", accessToken, 1*time.Hour)
	setCookie(c, "refreshToken", refreshToken, 24*time.Hour)

	// Return success response
	data := fiber.Map{
		"userID":   userID.Hex(),
		"username": user.Username,
	}
	return utils.JSONResponse(c, fiber.StatusCreated, data, "User registered successfully")
}

func LoginUser(c *fiber.Ctx) error {
	// get deatils from body
	// parse those deatils to Go struct format and validate
	// check user exists or not
	// check hashed password
	// generate refresh token
	// generate access token
	// add access token to cookies
	// Get the MongoDB collection
	collection := config.GetCollection(constants.USER_COLLECTION)

	// Parse login request body
	var loginRequest types.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, nil, "Invalid syntax in request body")
	}

	// Validate required fields
	if loginRequest.Email == "" || loginRequest.Password == "" {
		return utils.JSONResponse(c, fiber.StatusBadRequest, nil, "Email or Password is missing")
	}

	// Find user by email
	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{"email": loginRequest.Email}).Decode(&existingUser)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, nil, "No user found with this email")
	}

	// Validate password
	isMatch, err := utils.ComparePassword(existingUser.Password, loginRequest.Password)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, nil, "Error while comparing passwords")
	}
	if !isMatch {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, nil, "Invalid email or password")
	}

	// Generate tokens
	refreshToken, accessToken, err := generateTokens(existingUser.ID.Hex(), existingUser.Email)
	if err != nil {
		utils.JSONResponse(c, fiber.StatusInternalServerError, nil, "Error generating Refresh/Access token")
	}

	// Update refresh token in the database
	update := bson.M{"$set": bson.M{"refreshToken": refreshToken}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": existingUser.ID}, update)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, nil, "Failed to update refresh token in database")
	}

	// Set cookies
	setCookie(c, "accessToken", accessToken, 1*time.Hour)
	setCookie(c, "refreshToken", refreshToken, 24*time.Hour)

	// Return success response
	data := fiber.Map{
		"userID":   existingUser.ID.Hex(),
		"username": existingUser.Username,
	}
	return utils.JSONResponse(c, fiber.StatusOK, data, "User login successful")
}

func LogoutUser(c *fiber.Ctx) error {
	setCookie(c, "accessToken", "", -time.Hour)
	setCookie(c, "refreshToken", "", -time.Hour)
	return utils.JSONResponse(c, fiber.StatusOK, nil, "Successfully logged out")
}
