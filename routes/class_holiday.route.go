package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
)

func ClassHolidayRoute(classHoliday fiber.Router) {

	classHoliday.Post("/class-holiday", controllers.CreateClassHoliday)
	classHoliday.Get("/class-holiday", controllers.GetListClassHoliday)
	classHoliday.Get("/class-holiday/:id", controllers.GetDetailClassHoliday)
	classHoliday.Patch("/class-holiday/:id", controllers.UpdateClassHoliday)
	classHoliday.Delete("/class-holiday/:id", controllers.DeleteClassHoliday)

}
