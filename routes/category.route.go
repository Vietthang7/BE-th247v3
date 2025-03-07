package routes

import (
	"intern_247/controllers"
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
)

func CategoryRoute(category fiber.Router) {

	category.Post("/category", mdw.Gate("category", "create"), controllers.CreateCategory)
	category.Get("/category/:id", mdw.Gate("category", "read"), controllers.ReadCategory)
	category.Get("/category", controllers.ReadListCategory)
	category.Patch("/category/:id", mdw.Gate("category", "update"), controllers.UpdateCategory)
	category.Delete("/category/:id", mdw.Gate("category", "delete"), controllers.DeleteCategory)
}
