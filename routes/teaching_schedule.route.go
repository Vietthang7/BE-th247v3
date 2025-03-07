package routes

import (
	"intern_247/controllers"
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
)

func TeachingScheduleRoute(TeachingSchedule fiber.Router) {

	TeachingSchedule.Post("/teach-schedule", mdw.Gate("TeachingSchedule", "create"), controllers.CreateTeachSchedule)
	TeachingSchedule.Delete("/teach-schedule/:id", mdw.Gate("TeachingSchedule", "create"), controllers.DeleteTeachSchedule)
	TeachingSchedule.Get("/teach-schedule/:id", mdw.Gate("TeachingSchedule", "get"), controllers.GetTeachSchedule)
	TeachingSchedule.Get("/teach-schedule", mdw.Gate("TeachingSchedule", "get"), controllers.GetListTeachSchedule)
	TeachingSchedule.Put("/teach-schedule/:id", mdw.Gate("TeachingSchedule", "update"), controllers.UpdateTeachSchedule)
}
