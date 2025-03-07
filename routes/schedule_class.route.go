package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
	mdw "intern_247/middleware"
)

func ScheduleClassRoute(scheduleclass fiber.Router) {
	scheduleclass.Use(mdw.AdminAuthentication)
	scheduleclass.Post("/schedule-class", controllers.CreateScheduleClass)
}
