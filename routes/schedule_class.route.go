package routes

import (
	"intern_247/controllers"

	"github.com/gofiber/fiber/v2"
)

func ScheduleClassRoute(scheduleclass fiber.Router) {
	scheduleclass.Post("/schedule-class", controllers.CreateScheduleClass)
}
