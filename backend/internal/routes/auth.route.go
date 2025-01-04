package routes

import (
	"backend/internal/controllers"
	"backend/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app fiber.Router) {
	app.Post("/register", controllers.RegisterUser)
	app.Post("/login", controllers.LoginUser)

	// Logout route (with authentication)
	app.Post("/logout", middlewares.Authenticate(), controllers.LogoutUser)
}
