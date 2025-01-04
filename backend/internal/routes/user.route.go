package routes

import (
	"backend/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app fiber.Router) {
	app.Get("/me", controllers.GetUser)
	app.Delete("/delete", controllers.DeleteUser)
	app.Post("/refresh", controllers.RefreshTokens)
}
