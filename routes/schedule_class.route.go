package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
)

func ScheduleClassRoute(scheduleclass fiber.Router) {
	scheduleclass.Post("/schedule-class", controllers.CreateScheduleClass)
}
