package controllers

//import (
//	"fmt"
//	"github.com/gofiber/fiber/v2"
//	"github.com/google/uuid"
//	"intern_247/consts"
//	"intern_247/models"
//	"intern_247/repo"
//)

//func ListStudentByEnrollmentPlan(c *fiber.Ctx) error {
//	user, ok := c.Locals("user").(models.User)
//	if !ok {
//		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", "")
//	}
//	var (
//		err        error
//		entries    []*models.Student
//		pagination consts.RequestTable
//		query      = "students.center_id = ?"
//		args       = []interface{}{*user.CenterId}
//	)
//	classId, err := uuid.Parse(c.Params("id"))
//	if err != nil {
//		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
//	}
//	if c.Query("alphabetName") == "true" {
//		pagination = consts.BindRequestTable(c, "name")
//		pagination.Dir = "asc"
//	} else {
//		pagination = consts.BindRequestTable(c, "created_at")
//	}
//	if pagination.Search != "" {
//		query += " AND students.full_name LIKE ? OR students.email LIKE ? OR students.phone LIKE ?"
//		args = append(args, "%"+pagination.Search+"%", "%"+pagination.Search+"%", "%"+pagination.Search+"%")
//	}
//	fmt.Println(args)
//	if entries, err = repo.ListStudentByEnrollmentPlan(classId, *user.CenterId, &pagination, query, args); err != nil {
//		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
//	}
//	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
//		"data":       entries,
//		"pagination": pagination,
//	})
//}
