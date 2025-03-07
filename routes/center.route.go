package routes

import (
	"intern_247/controllers"
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
)

func CenterRouter(center fiber.Router) {
	center.Get("/center", mdw.Gate("center", "read"), controllers.ReadCenter)
	center.Put("/center", mdw.Gate("center", "update"), controllers.UpdateCenter)
	center.Post("/center", controllers.CreateCenter)
}
