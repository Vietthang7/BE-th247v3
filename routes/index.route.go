package routes

import (
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func AdminRoutes(app *fiber.App) {
	app.Use(cors.New())
	admin := app.Group("/api/admin")
	admin.Use(logger.New())
	admin.Use(mdw.AdminAuthentication)
	auth := app.Group("/api/auth")
	auth.Use(logger.New())
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
