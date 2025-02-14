package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
	mdw "intern_247/middleware"
)

func PermissionRoutes(permission fiber.Router) {
	permission.Use(mdw.AdminAuthentication)
	permission.Post("/permissions", controllers.CreatePermission)
	// Tag quyền
	permission.Post("/permission-tags", controllers.CreatePermissionTag)
	permission.Get("/permission-tags", controllers.ListPermissionTags)
	permission.Get("/permission-tags/:id", controllers.ReadPermissionTag)
}
