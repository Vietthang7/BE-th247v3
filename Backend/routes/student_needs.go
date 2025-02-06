package routes

import (
	"intern_247/controllers"
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
)

func StudentNeedsRouter(student fiber.Router) {

	student.Use(mdw.AdminAuthentication)

	// Nhu cầu học tập - Học viên
	student.Post("/study-needs", controllers.CreateStudyNeeds)
	// student.Get("/list-study-needs", controllers.ReadStudyNeeds)
	// student.Get("/study-needs/:student_id", controllers.GetStudyNeedsByStudentID)
	student.Get("/list-study-needs", controllers.ReadStudyNeeds)
	student.Get("/study-needs/:student_id", controllers.ReadStudyNeeds)
	student.Put("/study-needs/:student_id", controllers.UpdateStudyNeeds)

	// admin.Put("/study-needs/:studentId", mdw.Gate2("update", "student", "potential-student", "trial-student"), controllers.UpdateStudyNeeds)

}
