package routes

import (
	"intern_247/controllers"
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
)

func SubjectRoutes(subject fiber.Router) {
	subject.Post("/subject", mdw.Gate("subject", "create"), controllers.CreateSubject)
	subject.Patch("/subject", mdw.Gate("subject", "update"), controllers.UpdateSubject)
	subject.Delete("/subject/:id", mdw.Gate("subject", "delete"), controllers.DeleteSubject)
	subject.Get("/subject", mdw.Gate("subject", "read"), controllers.GetDetailSubject)
	subject.Get("/subjects", mdw.Gate("subject", "list"), controllers.GetListSubjects)
	subject.Get("/subject/all", mdw.Gate("subject", "list"), controllers.GetAllSubject)
}
