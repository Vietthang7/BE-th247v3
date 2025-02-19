package routes

import (
	"intern_247/controllers"
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
)

func StudentNeedsRouter(student fiber.Router) {

	student.Use(mdw.AdminAuthentication)

	// Nhu cầu học tập - Học viên
	student.Post("/study-needs", mdw.Gate2("create", "student", "potential-student", "trial-student"), controllers.CreateStudyNeeds)
	student.Get("/list-study-needs", mdw.Gate2("read", "student", "potential-student", "trial-student"), controllers.ReadStudyNeeds)
	student.Get("/study-needs/:id", mdw.Gate2("read", "student", "potential-student", "trial-student"), controllers.ReadStudyNeeds)
	student.Put("/study-needs/:id", mdw.Gate2("update", "student", "potential-student", "trial-student"), controllers.UpdateStudyNeeds)
	student.Delete("/study-needs/:id", mdw.Gate2("read", "student", "potential-student", "trial-student"), controllers.DeleteStudyNeeds)

}
