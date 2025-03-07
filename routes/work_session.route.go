package routes

import (
	"intern_247/controllers"
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
)

func WorkSessionRoute(WorkSession fiber.Router) {
	WorkSession.Post("/work-session", mdw.Gate("work_session", "create"), controllers.CreateWorkSession)
	WorkSession.Get("/work-sessions", mdw.Gate("work_session", "list"), controllers.GetListWorkSessions)
	WorkSession.Get("/work-sessions/schedule", mdw.Gate("work_session", "list"), controllers.ListWorkSessionForSchedule)
	WorkSession.Get("/work-session", mdw.Gate("work_session", "read"), controllers.GetWorkSessionDetail)
	WorkSession.Patch("/work-session", mdw.Gate("work_session", "update"), controllers.UpdateWorkSession)
	WorkSession.Delete("/work-session/:id", mdw.Gate("work_session", "delete"), controllers.DeleteWorkSession)
}
