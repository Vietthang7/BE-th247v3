package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	mdw "intern_247/middleware"
)

func AdminRoutes(app *fiber.App) {
	// Define a group for admin routes
	admin := app.Group("/api/admin", mdw.AdminAuthentication)
	auth := app.Group("/api/auth")

	// Call CenterRouter to register its routes
	app.Use(cors.New())
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
