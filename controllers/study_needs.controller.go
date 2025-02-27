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
	"gorm.io/gorm"
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

	// Kiểm tra các trường bắt buộc
	if entry.StudentId == uuid.Nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Student ID không được để trống")
	}
	if entry.BranchId == nil || *entry.BranchId == uuid.Nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Branch ID không được để trống")
	}
	if len(entry.SubjectIds) != 1 {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Chỉ được nhập 1 môn học")
	}

	// Kiểm tra Student có tồn tại không
	if err := repo.CheckStudentExists(entry.StudentId); err != nil {
		return ResponseError(c, fiber.StatusNotFound, consts.GetFailed, "Không tìm thấy học viên")
	}

	// Kiểm tra Branch có đang hoạt động không
	if err := repo.CheckBranchIsActive(*entry.BranchId); err != nil {
		return ResponseError(c, fiber.StatusForbidden, consts.InvalidInput, err.Error())
	}

	// Kiểm tra Subject có tồn tại không
	if err := repo.CheckSubjectsExist(entry.SubjectIds); err != nil {
		return ResponseError(c, fiber.StatusNotFound, consts.GetFailed, "Môn học không tồn tại")
	}

	// Kiểm tra WorkSession có tồn tại không
	for _, timeSlot := range entry.TimeSlots {
		if err := repo.CheckWorkSessionExists(timeSlot.WorkSessionId); err != nil {
			return ResponseError(c, fiber.StatusNotFound, consts.GetFailed,
				fmt.Sprintf("Work session với ID %s không tồn tại", timeSlot.WorkSessionId))
		}
	}

	// Gán CenterId từ user
	entry.CenterId = *user.CenterId

	// Ghi dữ liệu vào database
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
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid StudyNeeds ID")
		}

		studyNeeds, err := repo.GetStudyNeedsByID(studyNeedsID, *user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusNotFound, consts.GetFailed, "Study needs not found")
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, studyNeeds)
	}

	// Nếu không có study_needs_id, trả về tất cả
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
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid StudyNeeds ID")
	}

	var updatedData repo.StudyNeeds
	if err := c.BodyParser(&updatedData); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid request body")
	}

	if updatedData.BranchId != nil {
		var branch models.Branch
		err := app.Database.DB.Where("id = ?", updatedData.BranchId).First(&branch).Error
		if err != nil {
			logrus.Error("Invalid BranchId:", err)
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid BranchId")
		}
	}

	if len(updatedData.SubjectIds) > 0 {
		for _, subjectId := range updatedData.SubjectIds {
			var subject models.Subject
			err := app.Database.DB.Where("id = ?", subjectId).First(&subject).Error
			if err != nil {
				logrus.Error("Invalid SubjectId:", err)
				return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid SubjectId")
			}
		}
	}

	if len(updatedData.SubjectIds) > 1 {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Chỉ được nhập 1 môn học")
	}

	if err := updatedData.Update(studyNeedsID, *user.CenterId); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFailed, "Failed to update study needs")
	}

	var updatedStudyNeeds repo.StudyNeeds
	if err := app.Database.DB.Where("id = ?", studyNeedsID).First(&updatedStudyNeeds).Error; err != nil {
		logrus.Error("Failed to fetch updated StudyNeeds:", err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFailed, "Failed to retrieve updated study needs")
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, updatedStudyNeeds)
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

	if err := app.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("student_id = ? AND center_id = ?", studyNeeds.StudentId, user.CenterId).
			Delete(&repo.StudentSchedule{}).Error; err != nil {
			logrus.Error("Failed to delete StudentSchedule:", err)
			return err
		}

		if err := tx.Where("student_id = ?", studyNeeds.StudentId).
			Delete(&repo.StudentSchedule{}).Error; err != nil {
			logrus.Error("Failed to delete StudentScheduleData:", err)
			return err
		}

		if err := tx.Where("student_id = ?", studyNeeds.StudentId).
			Delete(&repo.TimeSlot{}).Error; err != nil {
			logrus.Error("Failed to delete TimeSlots:", err)
			return err
		}

		if err := tx.Where("student_id = ?", studyNeeds.StudentId).
			Delete(&repo.Shift{}).Error; err != nil {
			logrus.Error("Failed to delete ShortShifts:", err)
			return err
		}

		if err := tx.Where("student_id = ?", studyNeeds.StudentId).
			Delete(&repo.StudentSubject{}).Error; err != nil {
			logrus.Error("Failed to delete StudentSubject:", err)
			return err
		}

		if err := tx.Delete(&studyNeeds).Error; err != nil {
			logrus.Error("Failed to delete StudyNeeds:", err)
			return err
		}

		return nil
	}); err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.DeleteFail, "Failed to delete study needs and related data")
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, "Deleted successfully")
}
