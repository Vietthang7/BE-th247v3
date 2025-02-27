package routes

import (
	"github.com/gofiber/fiber/v2"
	mdw "intern_247/middleware"
)

func ClassRoute(class fiber.Router) {
	class.Use(mdw.AdminAuthentication)

}
