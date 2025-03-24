package controllers

import (
	"errors"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ListStudentByEnrollmentPlan(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized", consts.ERROR_UNAUTHORIZED)
	}
	var (
		err     error
		entries []repo.StudentCanBeAddedIntoClass
		class   models.Class
	)
	class, err = repo.GetClassByIdAndCenterId(uuid.MustParse(c.Params("id")), *user.CenterId)
	switch {
	case err == nil:
		if entries, err = repo.FindStudentsCanBeAddedIntoClass(class,
			strings.TrimSpace(c.Query("search"))); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.GetFailed)
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, entries)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, "Error class not found: "+err.Error(), consts.DataNotFound)
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, "Error get class: "+err.Error(), consts.GetFailed)
	}
}

// func AddStudentToClass(c *fiber.Ctx) error {
// 	token, err := repo.GetTokenData(c)
// 	if err != nil {
// 		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
// 	}
// 	if token.RoleId == consts.Student {
// 		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
// 	}
// 	var (
// 		input      []models.StudentToClass
// 		studentIds []uuid.UUID
// 	)
// 	if err := c.BodyParser(&input); err != nil {
// 		return ResponseError(c, fiber.StatusBadRequest, err.Error(), consts.InvalidReqInput)
// 	}
// 	// Tạo bản đồ để lưu trữ các lớp học của từng học sinh
// 	studentClasses := make(map[uuid.UUID][]uuid.UUID)
// 	for i := range input {
// 		studentIds = append(studentIds, input[i].StudentId)
// 		studentClasses[input[i].StudentId] = append(studentClasses[input[i].StudentId], input[i].ClassId)
// 	}
// 	// Kiểm tra lịch trùng của từng học sinh
// 	// for i := range studentIds {
// 	// 	if len(studentClasses[studentIds[i]]) > 0 {
// 	// 		// lấy ra lịch cũ của học viên
// 	// 		scheduleClasses, err := repo.GetScheduleClassByStudentId(studentIds[i], token.CenterId)
// 	// 		if err != nil {
// 	// 			logrus.Error(err)
// 	// 			return ResponseError(c, fiber.StatusBadRequest, "failed get student schedule", consts.InvalidReqInput)
// 	// 		}
// 	// 		// để lấy lịch học của các lớp mới mà học sinh muốn tham gia.
// 	// 		newScheduleClasses, err := repo.GetScheduleClassesByClassIdsAndCenterId(studentClasses[studentIds[i]], token.CenterId)
// 	// 		if err != nil {
// 	// 			logrus.Error(err)
// 	// 			return ResponseError(c, fiber.StatusBadRequest, "failed get student schedule", consts.InvalidReqInput)
// 	// 		}
// 	// 		// Kiểm tra xung đột lịch học giữa các lớp mới
// 	// 		for i := range newScheduleClasses {
// 	// 			for j := i + 1; j < len(newScheduleClasses); j++ {
// 	// 				if newScheduleClasses[i].ClassId != newScheduleClasses[j].ClassId {
// 	// 					scStartAt1 := utils.MixedDateAndTime(newScheduleClasses[i].StartDate, newScheduleClasses[i].StartTime)
// 	// 					scEndAt1 := utils.MixedDateAndTime(newScheduleClasses[i].StartDate, newScheduleClasses[i].EndTime)
// 	// 					scStartAt2 := utils.MixedDateAndTime(newScheduleClasses[j].StartDate, newScheduleClasses[j].StartTime)
// 	// 					scEndAt2 := utils.MixedDateAndTime(newScheduleClasses[j].StartDate, newScheduleClasses[j].EndTime)
// 	// 					if utils.IsTimeRangeOverlap(*scStartAt1, *scEndAt1, *scStartAt2, *scEndAt2) {
// 	// 						return ResponseError(c, fiber.StatusBadRequest, "invalid", consts.ERROR_CLASS_STUDENT_CONFLICT_SCHEDULE)
// 	// 					}
// 	// 				}
// 	// 			}
// 	// 		}
// 	// 		// Kiểm tra xung đội lịch học giữa các lớp cũ và mới
// 	// 		for _, sc := range scheduleClasses {
// 	// 			scStartAt := utils.MixedDateAndTime(sc.StartDate, sc.StartTime)
// 	// 			scEndAt := utils.MixedDateAndTime(sc.StartDate, sc.EndTime)
// 	// 			for _, nsc := range newScheduleClasses {
// 	// 				nscStartAt := utils.MixedDateAndTime(nsc.StartDate, nsc.StartTime)
// 	// 				nscEndAt := utils.MixedDateAndTime(nsc.StartDate, nsc.EndTime)
// 	// 				if utils.IsTimeRangeOverlap(*scStartAt, *scEndAt, *nscStartAt, *nscEndAt) {
// 	// 					return ResponseError(c, fiber.StatusBadRequest, nsc.Class.Name, consts.ERROR_CLASS_STUDENT_CONFLICT_SCHEDULE)
// 	// 				}
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }
// 	if err := repo.AddStudentToClass(input, token, c); err != nil {
// 		return ResponseError(c, fiber.StatusInternalServerError, "invalid", consts.ERROR_INTERNAL_SERVER_ERROR)
// 	}
// 	return ResponseSuccess(c, fiber.StatusCreated, "Thêm học viên thành công", consts.CREATE_SUCCESS)
// }

func AddStudentToClass(c *fiber.Ctx) error {
	token, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
	}
	if token.RoleId == consts.Student {
		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
	}
	var input models.AddStudentsToClassInput
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, err.Error(), consts.InvalidReqInput)
	}
	if err := repo.AddStudentToClass(input, token, c); err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "invalid", consts.ERROR_INTERNAL_SERVER_ERROR)
	}
	return ResponseSuccess(c, fiber.StatusCreated, "Thêm học viên thành công", consts.CREATE_SUCCESS)
}

func ListStudentInClass(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}

	var (
		err        error
		entries    []*models.Student
		pagination consts.RequestTable
		query      = "center_id = ?"
		args       = []interface{}{*user.CenterId}
	)

	classId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}

	if c.Query("alphabetName") == "true" {
		pagination = consts.BindRequestTable(c, "name")
		pagination.Dir = "asc"
	} else {
		pagination = consts.BindRequestTable(c, "created_at")
	}
	if pagination.Search != "" {
		query += " AND (full_name LIKE ? OR email LIKE ? OR phone LIKE ?)"
		args = append(args, "%"+pagination.Search+"%", "%"+pagination.Search+"%", "%"+pagination.Search+"%")
	}
	if entries, err = repo.ListStudentInClass(classId, &pagination, query, args); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
	pagination.Total = repo.CountStudentInClass(classId)
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}

func RemoveStudentFromClass(c *fiber.Ctx) error {
	token, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
	}
	if token.RoleId == consts.Student {
		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
	}

	var input models.RemoveStudentsFromClassInput
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, err.Error(), consts.InvalidReqInput)
	}

	if err := repo.RemoveStudentFromClass(input, token, c); err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "invalid", consts.ERROR_INTERNAL_SERVER_ERROR)
	}

	return ResponseSuccess(c, fiber.StatusOK, "Xóa học viên thành công", consts.DELETE_SUCCESS)
}
