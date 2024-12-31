package main

import (
	"backend/config"
	"backend/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	//Load env var
	config.LoadEnv()

	//Initialize Fiber App
	app := fiber.New()

	// Conect MongoDB
	config.ConnectDB()

	// Setup routes
	routes.SetupRoutes(app)

	// Get env variable and Start the server
	port := config.GetEnv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("Starting server on port", port)
	log.Fatal(app.Listen(":" + port))
}
