package routes

import (
	"intern_247/controllers"

	"github.com/gofiber/fiber/v2"
)

func Document(document fiber.Router) {
	document.Post("/documents", controllers.CreateDocument)
	document.Get("/documents/:id", controllers.ReadDocument)
	document.Get("/documents", controllers.ListDocuments)
	// document.Put("/documents/:id", controllers.UpdateDocument)
	// document.Delete("/documents/:id", controllers.DeleteDocument)

}
