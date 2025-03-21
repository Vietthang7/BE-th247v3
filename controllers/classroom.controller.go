package controllers

import (
	"errors"
	"fmt"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateClassroom(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}
	var (
		err  error
		form models.CreateClassroomForm
	)
	if err = c.BodyParser(&form); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}
	// name validation
	_, err = repo.FirstClassroom("center_id = ? AND name = ?", []interface{}{*user.CenterId, form.Name})
	switch {
	case err == nil:
		return ResponseError(c, fiber.StatusConflict, "Phòng học đã tồn tại", consts.ClassroomExistence)
	case errors.Is(err, gorm.ErrRecordNotFound):
		entry := models.Classroom{
			Name:     form.Name,
			CenterId: user.CenterId,
			BranchId: &form.BranchId,
			IsOnline: form.IsOnline,
			RoomType: form.RoomType,
			Metadata: form.Metadata,
			Slots:    form.Slots,
			IsActive: form.IsActive,
		}
		if err = repo.TsCreateClassroom(&entry, form.TimeSlots, form.ShortShifts); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.CreateFail, err.Error()), consts.CreateFailed)
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry)
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, "Error get classroom: "+err.Error(),
			consts.GetFailed)
	}
}
func UpdateClassroom(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}
	var (
		entry models.Classroom
		form  models.CreateClassroomForm
		err   error
	)
	entry, err = repo.FirstClassroom("id = ? AND center_id = ?", []interface{}{c.Params("id"),
		*user.CenterId}, "Schedule")
	switch {
	case err == nil:
		origin := entry
		if err = c.BodyParser(&form); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusBadRequest, fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
		}
		if form.IsOnline == nil {
			return ResponseError(c, fiber.StatusBadRequest, "Error null is_online", consts.InvalidReqInput)
		}
		entry.BranchId = &form.BranchId
		entry.Name = form.Name
		entry.IsOnline = form.IsOnline
		entry.RoomType = form.RoomType
		entry.Metadata = form.Metadata
		entry.Slots = form.Slots
		entry.IsActive = form.IsActive

		// name validation
		if _, err = repo.FirstClassroom("id <> ? AND center_id = ? AND name = ?", []interface{}{c.Params("id"), *user.CenterId, entry.Name}); err == nil {
			return ResponseError(c, fiber.StatusConflict, "Phòng học đã tồn tại", consts.ClassroomExistence)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, "Error get classroom: "+err.Error(),
				consts.GetFailed)
		}
		if err = repo.TsUpdateClassroom(&entry, origin, form.TimeSlots, form.ShortShifts); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.UpdateFail, err.Error()), consts.DeletedFailed)
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, nil)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, consts.NotFound, consts.DataNotFound)
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}
func ReadClassroom(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}
	var (
		entry models.Classroom
		err   error
	)
	entry, err = repo.FirstClassroom("id = ? AND center_id = ?", []interface{}{c.Params("id"), *user.CenterId}, "Branch", "Schedule")
	switch {
	case err == nil:
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, consts.NotFound+": "+err.Error(), consts.DataNotFound)
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail+": "+err.Error(), consts.GetFailed)
	}
}

func ListClassrooms(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}
	var (
		err        error
		entries    []models.Classroom
		pagination = consts.BindRequestTable(c, "created_at")
		query      = "center_id = ?"
		args       = []interface{}{user.CenterId}
	)
	if user.BranchId != nil {
		query += " AND branch_id = ?"
		args = append(args, *user.BranchId)
	}
	if c.Query("branch") != "" {
		query += " AND branch_id = ?"
		args = append(args, c.Query("branch"))
	}
	if c.Query("online") != "" {
		query += " AND is_online = ?"
		isOnline, _ := strconv.ParseBool(c.Query("online"))
		args = append(args, isOnline)
	}
	if c.Query("active") != "" {
		query += " AND is_active = ?"
		isActive, _ := strconv.ParseBool(c.Query("active"))
		args = append(args, isActive)
	}
	if entries, err = repo.FindClassrooms(&pagination, query, args); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
	pagination.Total = repo.CountClassroom(query, args)
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}
func DeleteClassroom(c *fiber.Ctx) error {
	var (
		err   error
		entry models.Classroom
	)
	entry, err = repo.FirstClassroom("id = ?", []interface{}{c.Params("id")})
	switch {
	case err == nil:
		var isArranged bool
		if isArranged, err = repo.ClassroomIsArranged(c.Params("id")); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
		} else {
			if isArranged {
				return ResponseError(c, fiber.StatusBadRequest,
					"Đã có lớp học được gán. Không thể xóa", consts.ClassroomIsArranged)
			}
		}
		if err = repo.DeleteClassroom(&entry); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.DeleteFail, err.Error()), consts.DeletedFailed)
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, nil)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound,
			fmt.Sprintf("%s: %s", consts.NotFound, err.Error()), consts.DataNotFound)
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}
