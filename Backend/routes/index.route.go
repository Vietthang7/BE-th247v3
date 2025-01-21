package routes

import "github.com/gofiber/fiber/v2"

func AdminRoutes(app *fiber.App) {
	// Define a group for admin routes
	admin := app.Group("/api/admin")

	// Call CenterRouter to register its routes
	CenterRouter(admin)
}
