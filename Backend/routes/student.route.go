package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
	mdw "intern_247/middleware"
)

func StudentRouter(student fiber.Router) {
	student.Use(mdw.AdminAuthentication)
	//potential-student : sinh viên tiềm năng
	//trial-student : sinh viên học thử
	student.Post("/students", mdw.Gate2("create", "student", "potential-student", "trial-student"), controllers.CreateStudent)
	student.Get("/students/:id", mdw.Gate2("read", "student", "potential-student", "trial-student"), controllers.ReadStudent)
	student.Put("/students/:id", mdw.Gate2("update", "student", "potential-student", "trial-student"), controllers.UpdateStudent)
}
