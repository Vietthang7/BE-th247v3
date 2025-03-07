package controllers

import (
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateClassHoliday(c *fiber.Ctx) error {

	user, err := repo.GetTokenData(c)
	if err != nil || user.RoleId == consts.Student {
		return ResponseError(c, fiber.StatusUnauthorized, "Error: Permission denied", consts.ERROR_PERMISSION_DENIED)
	}

	var classHoliday models.ClassHoliday
	if err := c.BodyParser(&classHoliday); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid request body", consts.InvalidReqInput)
	}

	if classHoliday.CenterId == uuid.Nil {
		classHoliday.CenterId = user.CenterId
	}

	if classHoliday.Name == "" {
		return ResponseError(c, fiber.StatusBadRequest, "Tên ngày nghỉ là bắt buộc", consts.ERROR_CLASS_HOLIDAY_REQUIRED)
	}

	if classHoliday.StartAt.IsZero() || classHoliday.EndAt.IsZero() || classHoliday.ClassID == uuid.Nil {
		return ResponseError(c, fiber.StatusBadRequest, "Missing required fields", consts.InvalidReqInput)
	}

	if !repo.IsClassIDExist(classHoliday.ClassID) {
		return ResponseError(c, fiber.StatusBadRequest, "Class ID không tồn tại", consts.ERROR_INVALID_CLASS_ID)
	}

	if classHoliday.EndAt.Before(classHoliday.StartAt) {
		return ResponseError(c, fiber.StatusBadRequest, "End date must be after start date", consts.ERROR_START_DAY_MUST_SMALLER_THAN_END_DAY)
	}

	if repo.IsClassHolidayExist(classHoliday.Name) {
		return ResponseError(c, fiber.StatusBadRequest, "Lịch nghỉ này đã có trên hệ thống, vui lòng thêm mới ngày nghỉ khác", consts.ERROR_HOLIDAY_ALREADY_EXIST_IN_SYSTEM)
	}

	if err := repo.CreateClassHoliday(&classHoliday); err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed to create class holiday", nil)
	}

	return ResponseSuccess(c, fiber.StatusCreated, "Class holiday created successfully", classHoliday)
}

func GetListClassHoliday(c *fiber.Ctx) error {

	user, err := repo.GetTokenData(c)
	if err != nil || user.RoleId == consts.Student {
		return ResponseError(c, fiber.StatusUnauthorized, "Error: Permission denied", consts.ERROR_PERMISSION_DENIED)
	}

	classHolidays, err := repo.GetListClassHoliday()
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed to fetch class holidays", nil)
	}

	return ResponseSuccess(c, fiber.StatusOK, "Success", classHolidays)
}

func GetDetailClassHoliday(c *fiber.Ctx) error {

	user, err := repo.GetTokenData(c)
	if err != nil || user.RoleId == consts.Student {
		return ResponseError(c, fiber.StatusUnauthorized, "Error: Permission denied", consts.ERROR_PERMISSION_DENIED)
	}

	idParam := c.Params("id")
	classHolidayID, err := uuid.Parse(idParam)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid ID format", consts.ERROR_INVALID_CLASS_HOLIDAY_ID)
	}

	classHoliday, err := repo.GetDetailClassHoliday(classHolidayID)
	if err != nil {
		return ResponseError(c, fiber.StatusNotFound, "Class holiday not found", consts.ERROR_CLASS_HOLIDAY_NOT_FOUND)
	}

	return ResponseSuccess(c, fiber.StatusOK, "Success", classHoliday)
}

func DeleteClassHoliday(c *fiber.Ctx) error {

	user, err := repo.GetTokenData(c)
	if err != nil || user.RoleId == consts.Student {
		return ResponseError(c, fiber.StatusUnauthorized, "Error: Permission denied", consts.ERROR_PERMISSION_DENIED)
	}

	idParam := c.Params("id")
	classHolidayID, err := uuid.Parse(idParam)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid ID format", consts.ERROR_INVALID_CLASS_HOLIDAY_ID)
	}

	exists, err := repo.IsClassHolidayExistByID(classHolidayID)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Error checking class holiday", nil)
	}
	if !exists {
		return ResponseError(c, fiber.StatusNotFound, "Class holiday not found", consts.ERROR_CLASS_HOLIDAY_NOT_FOUND)
	}

	if err := repo.DeleteClassHoliday(classHolidayID); err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed to delete class holiday", nil)
	}

	return ResponseSuccess(c, fiber.StatusOK, "Class holiday deleted successfully", nil)
}

func UpdateClassHoliday(c *fiber.Ctx) error {
	// Xác thực người dùng
	user, err := repo.GetTokenData(c)
	if err != nil || user.RoleId == consts.Student {
		return ResponseError(c, fiber.StatusUnauthorized, "Error: Permission denied", consts.ERROR_PERMISSION_DENIED)
	}

	// Lấy ID từ URL
	idParam := c.Params("id")
	classHolidayID, err := uuid.Parse(idParam)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid ID format", consts.ERROR_INVALID_CLASS_HOLIDAY_ID)
	}

	// Lấy thông tin ngày nghỉ từ DB để kiểm tra
	existingHoliday, err := repo.GetDetailClassHoliday(classHolidayID)
	if err != nil {
		return ResponseError(c, fiber.StatusNotFound, "Class holiday not found", consts.ERROR_CLASS_HOLIDAY_NOT_FOUND)
	}

	// Nhận dữ liệu từ request body
	var updateData models.ClassHoliday
	if err := c.BodyParser(&updateData); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid request body", consts.InvalidReqInput)
	}

	// Gán giá trị mới nếu có
	if updateData.Name != "" {
		existingHoliday.Name = updateData.Name
	}
	if updateData.Description != "" {
		existingHoliday.Description = updateData.Description
	}
	if !updateData.StartAt.IsZero() {
		existingHoliday.StartAt = updateData.StartAt
	}
	if !updateData.EndAt.IsZero() {
		existingHoliday.EndAt = updateData.EndAt
	}
	if updateData.ClassID != uuid.Nil {
		existingHoliday.ClassID = updateData.ClassID
	}
	if updateData.IsChanged != nil {
		existingHoliday.IsChanged = updateData.IsChanged
	}
	existingHoliday.IsAuto = updateData.IsAuto

	// Kiểm tra logic ngày nghỉ (EndAt phải sau StartAt)
	if existingHoliday.EndAt.Before(existingHoliday.StartAt) {
		return ResponseError(c, fiber.StatusBadRequest, "End date must be after start date", consts.ERROR_START_DAY_MUST_SMALLER_THAN_END_DAY)
	}

	// Cập nhật vào database
	if err := repo.UpdateClassHoliday(existingHoliday); err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed to update class holiday", nil)
	}

	// Trả về kết quả thành công
	return ResponseSuccess(c, fiber.StatusOK, "Class holiday updated successfully", existingHoliday)
}
