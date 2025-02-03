package routes

import (
	"intern_247/controllers"

	"github.com/gofiber/fiber/v2"
)

func CenterRouter(center fiber.Router) {
	center.Post("/center", controllers.CreateCenter)
	//admin.Get("/center", controllers.ReadCenter)
}
