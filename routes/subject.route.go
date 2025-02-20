package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
	mdw "intern_247/middleware"
)

func SubjectRoutes(subject fiber.Router) {
	subject.Use(mdw.AdminAuthentication)
	subject.Post("/subject", mdw.Gate("subject", "create"), controllers.CreateSubject)
	subject.Patch("/subject", mdw.Gate("subject", "update"), controllers.UpdateSubject)
	subject.Delete("/subject", mdw.Gate("subject", "delete"), controllers.DeleteSubject)
	subject.Get("/subject", mdw.Gate("subject", "read"), controllers.GetDetailSubject)
	//subject.Get("/subjects", mdw.Gate("subject", "list"), controllers.GetListSubjects)
}
