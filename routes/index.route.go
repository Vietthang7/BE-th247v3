package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	mdw "intern_247/middleware"
)

func AdminRoutes(app *fiber.App) {
	app.Use(cors.New())
	// Define a group for admin routes
	admin := app.Group("/api/admin")
	admin.Use(logger.New())
	admin.Use(mdw.AdminAuthentication)
	//admin := app.Group("/api/admin", mdw.AdminAuthentication)
	auth := app.Group("/api/auth")
	CenterRouter(admin)
	AuthRoute(auth)
	StudentNeedsRoute(admin)
	StudentRouter(admin)
	UserRouter(admin)
	DocsCategoriesRoute(admin)
	NotificationRoute(admin)
	Sp_RequestRouter(admin)
	PermissionRoutes(admin)
	AddressRouter(admin)
	BranchRoute(admin)
	CategoryRoute(admin)
	SubjectRoutes(admin)
	ClassRoomRouter(admin)
	LessonRoute(admin)
	WorkSessionRoute(admin)
	DocumentRoute(admin)
	ClassRoute(admin)
	ClassHolidayRoute(admin)
	TeachingScheduleRoute(admin)
	ScheduleClassRoute(admin)
}
