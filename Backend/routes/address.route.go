package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
	mdw "intern_247/middleware"
)

func AddressRouter(address fiber.Router) {
	address.Use(mdw.AdminAuthentication)
	address.Get("/provinces", controllers.ListProvinces)
	address.Get("/districts", controllers.ListDistricts)
	address.Get("/wards", controllers.ListWards)
}
