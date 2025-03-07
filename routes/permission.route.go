package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
	mdw "intern_247/middleware"
)

func PermissionRoutes(permission fiber.Router) {
	// Quyền
	permission.Post("/permissions", controllers.CreatePermission)
	// Tag quyền
	permission.Post("/permission-tags", controllers.CreatePermissionTag)
	permission.Get("/permission-tags", controllers.ListPermissionTags)
	permission.Get("/permission-tags/:id", controllers.ReadPermissionTag)
	// Nhóm quyền
	permission.Post("/permission-grp", mdw.Gate("permission_grp", "create"), controllers.CreatePermissionGrp)
	permission.Get("/permission-grp", mdw.Gate("permission_grp", "list"), controllers.ListPermissionGrp)
	permission.Get("/permission-grp/:id", mdw.Gate("permission_grp", "read"), controllers.ReadPermissionGrp)
	permission.Put("/permission-grp/:id", mdw.Gate("permission_grp", "update"), controllers.UpdatePermissionGrp)
	permission.Delete("/permission-grp", mdw.Gate("permission_grp", "delete"), controllers.DeletePermissionGroup)
	//permission.Get("/permission-grp-excel", mdw.Gate("permission_grp", "export"), controllers.ExportListGrps)
}
