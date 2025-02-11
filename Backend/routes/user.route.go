package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
	mdw "intern_247/middleware"
)

func UserRouter(user fiber.Router) {
	user.Use(mdw.AdminAuthentication)
	user.Post("/users", mdw.Gate("hr", "create"), controllers.CreateUser)
	user.Get("/users", mdw.Gate("hr", "list"), controllers.ListUsers)
	user.Get("/users/:id", mdw.Gate("hr", "read"), controllers.ReadUser)
	user.Put("/users/:id", mdw.Gate("hr", "update"), controllers.UpdateUser)
}
