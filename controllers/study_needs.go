package controllers

import (
	"fmt"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func CreateStudyNeeds(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, consts.InvalidInput, "Permission denied!")
	}

	var entry repo.StudyNeeds
	if err := c.BodyParser(&entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}

	if err := repo.CheckStudentExists(entry.StudentId); err != nil {
		return ResponseError(c, fiber.StatusNotFound, consts.GetFailed, "Không tìm thấy học viên")
	}
	if entry.BranchId == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Branch ID không được để trống")
	}
	if err := repo.CheckBranchIsActive(*entry.BranchId); err != nil {
		return ResponseError(c, fiber.StatusForbidden, consts.InvalidInput, err.Error())
	}

	entry.CenterId = *user.CenterId

	if err := entry.Create(); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.CreateFail, err.Error()), consts.CreateFailed)
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry)
}

func ReadStudyNeeds(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, consts.InvalidInput, "Permission denied!")
	}

	studentIDParam := c.Params("student_id")
	if studentIDParam != "" {
		studentID, err := uuid.Parse(studentIDParam)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid Student ID")
		}

		studyNeeds, err := repo.GetStudyNeedsByStudentID(studentID, *user.CenterId)
		if err != nil || len(studyNeeds) == 0 {
			return ResponseError(c, fiber.StatusNotFound, consts.GetFailed, "Study needs not found")
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, studyNeeds)
	}

	// Nếu không có student_id, trả về tất cả
	studyNeeds, err := repo.GetAllStudyNeeds(*user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFailed, "Failed to retrieve study needs")
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, studyNeeds)
}

func UpdateStudyNeeds(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, consts.InvalidInput, "Permission denied!")
	}

	studentIDParam := c.Params("student_id")
	studentID, err := uuid.Parse(studentIDParam)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid Student ID")
	}

	var updatedStudyNeeds repo.StudyNeeds
	if err := c.BodyParser(&updatedStudyNeeds); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid request body")
	}

	if err := updatedStudyNeeds.Update(studentID, *user.CenterId); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFailed, "Failed to update study needs")
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, updatedStudyNeeds)
}
