package routes

import (
	"backend/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	authRoutes := api.Group("/auth")
	RegisterAuthRoutes(authRoutes)

	todoRoutes := api.Group("/todo")
	RegisterTODORoutes(todoRoutes)

	userRoutes := api.Group("/user", middlewares.Authenticate())
	RegisterUserRoutes(userRoutes)

}
