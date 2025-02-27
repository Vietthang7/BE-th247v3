package routes

import (
	"intern_247/controllers"

	"github.com/gofiber/fiber/v2"
)

func NotificationRoute(notification fiber.Router) {
	notification.Post("/notifications", controllers.CreateNotification)
	notification.Delete("/notifications", controllers.DeleteNotification)
	notification.Patch("/notifications", controllers.MarkNotificationIsRead)
	notification.Get("/notifications", controllers.ListNotification)
	notification.Put("/notifications/:id", controllers.UpdateNotification)

}
