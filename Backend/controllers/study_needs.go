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
	var (
		err   error
		entry repo.StudyNeeds
	)

	if err = c.BodyParser(&entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}

	entry.CenterId = *user.CenterId
	if err = entry.Create(); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.CreateFail, err.Error()), consts.CreateFailed)
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, nil)
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

// func GetStudyNeedsByStudentID(c *fiber.Ctx) error {
// 	// Lấy thông tin người dùng từ JWT token
// 	user, ok := c.Locals("user").(models.User)
// 	if !ok {
// 		return ResponseError(c, fiber.StatusForbidden, consts.InvalidInput, "Permission denied!")
// 	}

// 	// Lấy student_id từ URL parameters
// 	studentIDParam := c.Params("student_id")
// 	studentID, err := uuid.Parse(studentIDParam)
// 	if err != nil {
// 		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid Student ID")
// 	}

// 	var studyNeeds repo.StudyNeeds

// 	// Truy vấn cơ sở dữ liệu theo student_id và center_id
// 	if err := app.Database.DB.Where("student_id = ? AND center_id = ?", studentID, user.CenterId).First(&studyNeeds).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return ResponseError(c, fiber.StatusNotFound, consts.GetFailed, "Study needs not found")
// 		}
// 		logrus.Error(err)
// 		return ResponseError(c, fiber.StatusInternalServerError,
// 			fmt.Sprintf("%s: %s", consts.GetFailed, err.Error()), consts.GetFailed)
// 	}

// 	// Trả về dữ liệu thành công
// 	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, studyNeeds)
// }

func UpdateStudyNeeds(c *fiber.Ctx) error {
	// Lấy thông tin người dùng từ JWT token
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, consts.InvalidInput, "Permission denied!")
	}

	// Lấy student_id từ Path Variables
	studentIDParam := c.Params("student_id")
	studentID, err := uuid.Parse(studentIDParam)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid Student ID")
	}

	// Lấy dữ liệu từ body request
	var updatedStudyNeeds repo.StudyNeeds
	if err := c.BodyParser(&updatedStudyNeeds); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid request body")
	}

	// Tìm kiếm StudyNeeds của sinh viên với student_id
	var existingStudyNeeds repo.StudyNeeds
	if err := app.Database.DB.Where("student_id = ? AND center_id = ?", studentID, user.CenterId).First(&existingStudyNeeds).Error; err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, consts.GetFailed, "Study needs not found")
	}

	// Cập nhật các trường cần thiết
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

	// Lưu lại thay đổi
	if err := app.Database.DB.Save(&existingStudyNeeds).Error; err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFailed, "Failed to update study needs")
	}

	// Trả về kết quả thành công
	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, existingStudyNeeds)
}
