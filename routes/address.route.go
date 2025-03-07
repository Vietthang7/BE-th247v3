package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
)

func AddressRouter(address fiber.Router) {
	address.Get("/provinces", controllers.ListProvinces)
	address.Get("/districts", controllers.ListDistricts)
	address.Get("/wards", controllers.ListWards)
}
