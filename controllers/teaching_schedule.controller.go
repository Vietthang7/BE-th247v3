package controllers

import (
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateTeachSchedule xử lý request tạo lịch giảng dạy
func CreateTeachSchedule(c *fiber.Ctx) error {
	// Parse dữ liệu từ request body
	var form models.CreateTeachScheForm
	if err := c.BodyParser(&form); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid request body!")
	}

	// Gọi repository để xử lý logic tạo lịch giảng dạy
	newSchedule, err := repo.CreateTeachingSchedule(form)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.CreateFail, err.Error())
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, newSchedule)
}

// DeleteTeachSchedule xử lý request xóa lịch giảng dạy
func DeleteTeachSchedule(c *fiber.Ctx) error {
	scheduleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid schedule ID format!")
	}

	err = repo.DeleteTeachingSchedule(scheduleID)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.DeleteFail, err.Error())
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, "Teaching schedule deleted successfully!")
}

func GetTeachSchedule(c *fiber.Ctx) error {
	// Lấy ID từ request param
	scheduleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "ID không hợp lệ!")
	}

	// Gọi repo để lấy dữ liệu
	schedule, err := repo.ReadTeachSchedule(scheduleID)
	if err != nil {
		return ResponseError(c, fiber.StatusNotFound, consts.NotFound, "Không tìm thấy lịch giảng dạy!")
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, schedule)
}

func GetListTeachSchedule(c *fiber.Ctx) error {
	// Gọi repo để lấy dữ liệu
	schedule, err := repo.ListTeachSchedule()
	if err != nil {
		return ResponseError(c, fiber.StatusNotFound, consts.NotFound, "Không tìm thấy lịch giảng dạy!")
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, schedule)
}

// UpdateTeachScheduleController cập nhật lịch giảng dạy
func UpdateTeachSchedule(c *fiber.Ctx) error {
	// Lấy ID từ params
	scheduleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "ID không hợp lệ")
	}

	// Parse body request
	var form models.CreateTeachScheForm
	if err := c.BodyParser(&form); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid request body")
	}

	// Gọi hàm cập nhật từ repo
	schedule, err := repo.UpdateTeachSchedule(scheduleID, form)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFail, err.Error())
	}

	// Trả về kết quả JSON
	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, schedule)
}
