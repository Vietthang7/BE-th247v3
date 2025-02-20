package routes

import (
	"intern_247/controllers"

	"github.com/gofiber/fiber/v2"
)

func Sp_RequestRouter(sp_requests fiber.Router) {
	sp_requests.Post("/sp-requests", controllers.CreateSpRequest)
	sp_requests.Get("/sp-requests", controllers.ListSpRequests)
	sp_requests.Get("/sp-requests/:id", controllers.ReadSpRequests)
	sp_requests.Put("/sp-requests/:id", controllers.UpdateSpRequests)
	sp_requests.Put("/respond-sp-requests", controllers.RespondSpRequests)
	sp_requests.Delete("/sp-requests/:id", controllers.DeleteSpRequests)

}
