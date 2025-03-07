package routes

import (
	"intern_247/controllers"
	mdw "intern_247/middleware"

	"github.com/gofiber/fiber/v2"
)

func BranchRoute(branch fiber.Router) {
	// Chi nhánh
	branch.Post("/branches", mdw.Gate("branch", "create"), controllers.CreateBranch)
	branch.Get("/branches", mdw.Gate("branch", "list"), controllers.ListBranches)
	branch.Get("/branches/:id", mdw.Gate("branch", "read"), controllers.ReadBranch)
	branch.Put("/branches/:id", mdw.Gate("branch", "update"), controllers.UpdateBranch)
	branch.Delete("/branches/:id", mdw.Gate("branch", "delete"), controllers.DeleteBranch)
}
