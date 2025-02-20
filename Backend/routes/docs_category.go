package routes

import (
	"intern_247/controllers"
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
)

func DocsCategoriesRouter(docscategories fiber.Router) {

	docscategories.Use(mdw.AdminAuthentication)

	// Danh mục taif liệu
	docscategories.Post("/docs-categories", controllers.CreateDocsCategory)
	docscategories.Get("/docs-categories/:id", controllers.ReadDocsCategory)
	docscategories.Get("/docs-categories", controllers.ListDocsCategories)
	docscategories.Put("/docs-categories/:id", controllers.UpdateDocsCategory)
	docscategories.Delete("/docs-categories/:id", controllers.DeleteDocsCategory)

}
