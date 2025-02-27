package controllers

import (
	"errors"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateBranch(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}
	var (
		entry models.Branch
		err   error
	)
	if err = c.BodyParser(&entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, err.Error())
	}

	if _, err = repo.FirstBranch(app.Database.DB.Where("name = ? AND center_id = ?",
		entry.Name, *user.CenterId)); err == nil {
		return ResponseError(c, 0, "Tên chi nhánh đã tồn tại.", "")
	}
	entry.UserId = user.ID
	entry.CenterId = user.CenterId
	if err = repo.CreateBranch(app.Database.DB, &entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.CreateFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry)
}

func ReadBranch(c *fiber.Ctx) error {
	entry, err := repo.FirstBranch(app.Database.DB.Where(map[string]interface{}{
		"deleted_at": nil,
		"id":         c.Params("id"),
	}))
	switch {
	case err == nil:
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, consts.NotFound, err.Error())
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
}

func ListBranches(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}

	// Gọi repo để lấy danh sách chi nhánh
	branches, err := repo.ListBranches(app.Database.DB, *user.CenterId)
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}

	// Trả về danh sách chi nhánh
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, branches)
}

func UpdateBranch(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}

	// Lấy branchId từ URL parameter
	branchId := c.Params("id")
	parsedBranchId, err := uuid.Parse(branchId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "ID không hợp lệ!")
	}

	// Kiểm tra xem chi nhánh có tồn tại hay không
	var existingBranch models.Branch
	if err := app.Database.DB.Where("id = ?", parsedBranchId).First(&existingBranch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ResponseError(c, fiber.StatusNotFound, consts.NotFound, "Chi nhánh không tồn tại!")
		}
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFail, "Không thể kiểm tra chi nhánh!")
	}

	// Parse dữ liệu từ request body
	var updatedBranch models.Branch
	if err := c.BodyParser(&updatedBranch); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, err.Error())
	}

	// Gán lại thông tin UserId và CenterId từ user
	updatedBranch.UserId = user.ID
	updatedBranch.CenterId = user.CenterId
	updatedBranch.ID = parsedBranchId

	// Gọi repo để cập nhật chi nhánh
	if err := repo.UpdateBranch(app.Database.DB, parsedBranchId, &updatedBranch); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFail, err.Error())
	}

	// Trả về kết quả thành công
	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, updatedBranch)
}

func DeleteBranch(c *fiber.Ctx) error {
	// Lấy user từ context (kiểm tra quyền truy cập)
	// user, ok := c.Locals("user").(models.User)
	// if !ok {
	// 	return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	// }

	// Lấy branchId từ URL parameter
	branchId := c.Params("id")
	parsedBranchId, err := uuid.Parse(branchId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "ID không hợp lệ!")
	}

	// Kiểm tra xem chi nhánh có tồn tại hay không
	var existingBranch models.Branch
	if err := app.Database.DB.Where("id = ?", parsedBranchId).First(&existingBranch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ResponseError(c, fiber.StatusNotFound, consts.InvalidInput, "Chi nhánh không tồn tại!")
		}
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.DeleteFail, "Không thể kiểm tra chi nhánh!")
	}

	// Gọi repo để xóa chi nhánh
	if err := repo.DeleteBranch(app.Database.DB, parsedBranchId); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.DeleteFail, "Không thể xóa chi nhánh!")
	}

	// Trả về kết quả thành công
	return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, nil)
}
