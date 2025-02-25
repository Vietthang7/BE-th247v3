package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func AdminRoutes(app *fiber.App) {
	// Define a group for admin routes
	admin := app.Group("/api/admin")
	auth := app.Group("/api/auth")

	// Call CenterRouter to register its routes
	app.Use(cors.New())
	CenterRouter(admin)
	AuthRoute(auth)
	StudentNeedsRouter(admin)
	StudentRouter(admin)
	UserRouter(admin)
	DocsCategoriesRouter(admin)
	NotificationRouter(admin)
	Sp_RequestRouter(admin)
	PermissionRoutes(admin)
	AddressRouter(admin)
	Branch(admin)
	Category(admin)
	SubjectRoutes(admin)
	ClassRoomRouter(admin)
}
