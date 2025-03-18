package routes

import (
	"intern_247/controllers"
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
)

func ClassRoute(class fiber.Router) {
	class.Post("/class", mdw.Gate("class", "create"), controllers.CreateClass)
	class.Get("/class/:id", mdw.Gate("class", "read"), controllers.GetDetailClass)
	class.Patch("/class", mdw.Gate("class", "update"), controllers.UpdateClass)
	class.Get("/classes", mdw.Gate("class", "list"), controllers.GetListClasses)
	class.Delete("/class/:id", mdw.Gate("class", "delete"), controllers.DeleteClass)
	class.Patch("/class/cancel/:id", mdw.Gate("class", "cancel"), controllers.CanceledClass)
	//class.Get("/class/:id/student_enroll", mdw.Gate("class", "list"), controllers.ListStudentByEnrollmentPlan)
	//class.Get("/class/:id/student_enroll", mdw.Gate("class", "list"), controllers.ListStudentByEnrollmentPlan)
}
