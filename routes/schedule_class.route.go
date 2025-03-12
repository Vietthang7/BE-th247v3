package routes

import (
	"intern_247/controllers"

	"github.com/gofiber/fiber/v2"
)

func ScheduleClassRoute(scheduleclass fiber.Router) {
	scheduleclass.Post("/schedule-class", controllers.CreateScheduleClass)
	scheduleclass.Get("/schedule-class/:id", controllers.GetDetailScheduleClass)
	scheduleclass.Get("/schedule-classes/:id", controllers.GetListScheduleClass)
	scheduleclass.Get("/schedule-class/student", controllers.GetListScheduleClassForStudent)
}
