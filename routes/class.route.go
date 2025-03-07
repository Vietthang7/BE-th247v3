package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
	mdw "intern_247/middleware"
)

func ClassRoute(class fiber.Router) {
	class.Use(mdw.AdminAuthentication)
	class.Post("/class", mdw.Gate("class", "create"), controllers.CreateClass)
	class.Get("/class/:id", mdw.Gate("class", "read"), controllers.GetDetailClass)
	class.Patch("/class", mdw.Gate("class", "update"), controllers.UpdateClass)
	class.Get("/classes", mdw.Gate("class", "list"), controllers.GetListClasses)
	//class.Get("/class/:id/student_enroll", mdw.Gate("class", "list"), controllers.ListStudentByEnrollmentPlan)
}
