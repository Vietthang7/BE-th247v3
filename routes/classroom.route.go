package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
	mdw "intern_247/middleware"
)

func ClassRoomRouter(classroom fiber.Router) {
	classroom.Post("/classrooms", mdw.Gate("classroom", "create"), controllers.CreateClassroom)
	classroom.Put("/classrooms/:id", mdw.Gate("classroom", "update"), controllers.UpdateClassroom)
	classroom.Get("/classrooms/:id", mdw.Gate("classroom", "read"), controllers.ReadClassroom)
	classroom.Get("/classrooms", mdw.Gate("classroom", "list"), controllers.ListClassrooms)
	classroom.Delete("/classrooms/:id", mdw.Gate("classroom", "delete"), controllers.DeleteClassroom)
}
