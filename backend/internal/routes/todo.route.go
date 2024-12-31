package routes

import (
	"backend/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterTODORoutes(app fiber.Router) {
	app.Get("/", controllers.GetTodos)
	app.Post("/", controllers.AddTodo)
	app.Put("/:id", controllers.UpdateTodo)
	app.Delete("/:id", controllers.DeleteTodo)
}
