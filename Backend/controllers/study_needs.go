package controllers

import (
	"fmt"
	"intern_247/app"
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
	if err := repo.CheckStudentHasBranch(entry.StudentId); err == nil {
		return ResponseError(c, fiber.StatusConflict, consts.InvalidInput, "Học viên đã được gán chi nhánh trước đó")
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

	studyNeedsIDParam := c.Params("id")
	if studyNeedsIDParam != "" {
		studyNeedsID, err := uuid.Parse(studyNeedsIDParam)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid Study Needs ID")
		}

		studyNeeds, err := repo.GetStudyNeedsByID(studyNeedsID, *user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusNotFound, consts.GetFailed, "Study needs not found")
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, studyNeeds)
	}

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

	studyNeedsIDParam := c.Params("id")
	studyNeedsID, err := uuid.Parse(studyNeedsIDParam)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid Study Needs ID")
	}

	var updatedStudyNeeds repo.StudyNeeds
	if err := c.BodyParser(&updatedStudyNeeds); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid request body")
	}

	var existingStudyNeeds repo.StudyNeeds
	if err := app.Database.DB.Where("id = ? AND center_id = ?", studyNeedsID, user.CenterId).First(&existingStudyNeeds).Error; err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, consts.GetFailed, "Study needs not found")
	}

	if updatedStudyNeeds.BranchId != nil {
		var branch repo.Branch
		if err := app.Database.DB.Where("id = ?", updatedStudyNeeds.BranchId).First(&branch).Error; err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusNotFound, consts.GetFailed, "Không tìm thấy chi nhánh")
		}
		if branch.IsActive == nil || !*branch.IsActive {
			return ResponseError(c, fiber.StatusForbidden, consts.InvalidInput, "Chi nhánh không hoạt động")
		}
		existingStudyNeeds.BranchId = updatedStudyNeeds.BranchId
	}

	if updatedStudyNeeds.StudyGoals != "" {
		existingStudyNeeds.StudyGoals = updatedStudyNeeds.StudyGoals
	}
	if updatedStudyNeeds.TeacherRequirements != "" {
		existingStudyNeeds.TeacherRequirements = updatedStudyNeeds.TeacherRequirements
	}
	if updatedStudyNeeds.IsOnlineForm != nil {
		existingStudyNeeds.IsOnlineForm = updatedStudyNeeds.IsOnlineForm
	}
	if updatedStudyNeeds.IsOfflineForm != nil {
		existingStudyNeeds.IsOfflineForm = updatedStudyNeeds.IsOfflineForm
	}
	if updatedStudyNeeds.StudyingStartDate != nil {
		existingStudyNeeds.StudyingStartDate = updatedStudyNeeds.StudyingStartDate
	}
	if updatedStudyNeeds.BranchId != nil {
		existingStudyNeeds.BranchId = updatedStudyNeeds.BranchId
	}

	if err := app.Database.DB.Save(&existingStudyNeeds).Error; err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFailed, "Failed to update study needs")
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, existingStudyNeeds)
}

func DeleteStudyNeeds(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, consts.InvalidInput, "Permission denied!")
	}

	studyNeedsIDParam := c.Params("id")
	studyNeedsID, err := uuid.Parse(studyNeedsIDParam)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid Study Needs ID")
	}

	var studyNeeds repo.StudyNeeds
	if err := app.Database.DB.Where("id = ? AND center_id = ?", studyNeedsID, user.CenterId).
		First(&studyNeeds).Error; err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, consts.GetFailed, "Study needs not found")
	}

	if err := app.Database.DB.Delete(&studyNeeds).Error; err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.DeleteFail, "Failed to delete study needs")
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, "Deleted successfully")
}
