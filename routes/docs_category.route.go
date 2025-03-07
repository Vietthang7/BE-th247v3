package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
)

func DocsCategoriesRoute(docscategories fiber.Router) {
	// Danh mục taif liệu
	docscategories.Post("/docs-categories", controllers.CreateDocsCategory)
	docscategories.Get("/docs-categories/:id", controllers.ReadDocsCategory)
	docscategories.Get("/docs-categories", controllers.ListDocsCategories)
	docscategories.Put("/docs-categories/:id", controllers.UpdateDocsCategory)
	docscategories.Delete("/docs-categories/:id", controllers.DeleteDocsCategory)

}
