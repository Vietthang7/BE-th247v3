package routes

import (
	"intern_247/controllers"
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
)

func CenterRouter(center fiber.Router) {
	center.Use(mdw.AdminAuthentication)
	center.Post("/center", controllers.CreateCenter)
	//admin.Get("/center", controllers.ReadCenter)
}
