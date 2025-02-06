package routes

import "github.com/gofiber/fiber/v2"

func AdminRoutes(app *fiber.App) {
	// Define a group for admin routes
	admin := app.Group("/api/admin")
	auth := app.Group("/api/auth")
	CenterRouter(admin) // Call CenterRouter to register its routes
	AuthRoute(auth)
	StudentNeedsRouter(admin)
	StudentRouter(admin)
	UserRouter(admin)
}
