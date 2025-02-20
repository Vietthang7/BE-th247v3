package controllers

import (
	"errors"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateNotification(c *fiber.Ctx) error {
	// Lấy thông tin người dùng từ token
	token, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
	}

	var (
		entry repo.Notification
	)

	// Phân tích body request
	if err := c.BodyParser(&entry); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}

	// Kiểm tra StudentID có hợp lệ không
	studentExists, err := repo.CheckStudentID(*entry.To)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.InvalidInput, "Lỗi kiểm tra ID sinh viên")
	}
	if !studentExists {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "ID sinh viên không tồn tại")
	}

	// Đặt ID người gửi từ token
	entry.From = &token.ID

	// Tạo thông báo
	err = entry.Create()
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.CreateFail, consts.CreateFailed)
	}

	return ResponseSuccess(c, fiber.StatusCreated, consts.CreateSuccess, entry)
}

func ListNotification(c *fiber.Ctx) error {
	//_, err := repo.GetTokenData(c)
	token, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
	}
	var (
		entry      repo.Notification
		entries    []*models.Notification
		pagination = consts.BindRequestTable(c, "created_at")
		DB         = pagination.CustomOptions(app.Database.DB)
	)

	// Lọc theo người nhận (`to`)
	if c.Query("to") != "" {
		recipientID := c.Query("to")           // Giữ nguyên như chuỗi (có thể là UUID)
		DB = DB.Where("`to` = ?", recipientID) // Lọc theo chuỗi ID
	}

	if c.Query("is_read") != "" {
		status, _ := strconv.ParseBool(c.Query("is_read"))
		DB = DB.Where("is_read = ?", status)
	}

	if entries, err = entry.Find(DB.Where("`from` = ?", token.ID)); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}

	// if entries, err = entry.Find(DB); err != nil {
	// 	logrus.Error(err)
	// 	return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	// }

	pagination.Total = entry.Count(DB.Offset(-1))
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}

// func MarkNotificationIsRead(c *fiber.Ctx) error {
// 	_, err := repo.GetTokenData(c)
// 	if err != nil {
// 		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
// 	}
// 	var input models.ReqIds
// 	if err := c.BodyParser(&input); err != nil {
// 		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid input!")
// 	}
// 	err = app.Database.DB.Model(&models.Notification{}).Where("id IN ?", input.Ids).Update("is_read", true).Error
// 	if err != nil {
// 		return ResponseError(c, fiber.StatusBadRequest, consts.UpdateFail, "Can't delete record!")
// 	}
// 	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, nil)
// }

func MarkNotificationIsRead(c *fiber.Ctx) error {
	// Lấy thông tin token người dùng
	_, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
	}

	var input models.ReqIds
	// Parse dữ liệu từ body request
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid input!")
	}

	// Kiểm tra xem các ID có tồn tại trong cơ sở dữ liệu không
	var count int64
	err = app.Database.DB.Model(&models.Notification{}).
		Where("id IN ?", input.Ids).
		Count(&count).Error

	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, "Lỗi kiểm tra ID")
	}

	if count != int64(len(input.Ids)) {
		// Nếu số lượng ID không khớp, nghĩa là có ID không tồn tại
		return ResponseError(c, fiber.StatusBadRequest, consts.UpdateFail, "ID không tồn tại")
	}

	// Nếu tất cả ID tồn tại, tiến hành cập nhật
	err = app.Database.DB.Model(&models.Notification{}).Where("id IN ?", input.Ids).Update("is_read", true).Error
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.UpdateFail, "Can't update records!")
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, nil)
}

func DeleteNotification(c *fiber.Ctx) error {
	_, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
	}
	var input models.ReqIds
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	err = app.Database.DB.Delete(&models.Notification{}, input.Ids).Error
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.DeleteFail, "Can't delete record!")
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, nil)
}

func UpdateNotification(c *fiber.Ctx) error {
	// Lấy ID từ URL Params
	notificationID := c.Params("id")
	if notificationID == "" {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "ID is required!")
	}

	// Kiểm tra Notification có tồn tại không
	var notification models.Notification
	err := app.Database.DB.First(&notification, "id = ?", notificationID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ResponseError(c, fiber.StatusNotFound, consts.GetFail, "Notification not found!")
		}
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, "Error fetching notification!")
	}

	// Cấu trúc dữ liệu cập nhật
	type UpdateNotificationInput struct {
		Title    *string `json:"title"` // Sử dụng con trỏ để phân biệt giữa null và rỗng
		Content  *string `json:"content"`
		Metadata *string `json:"metadata"`
	}

	var input UpdateNotificationInput

	// Parse dữ liệu từ body request
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid input data!")
	}

	// Tạo map lưu các trường cần cập nhật
	updates := map[string]interface{}{}

	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Content != nil {
		updates["content"] = *input.Content
	}
	if input.Metadata != nil {
		updates["metadata"] = *input.Metadata // Có thể cập nhật thành chuỗi rỗng hoặc giá trị mới
	}

	// Nếu không có trường nào để cập nhật
	if len(updates) == 0 {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "No fields to update!")
	}

	// Thực hiện cập nhật trong DB
	err = app.Database.DB.Model(&notification).Updates(updates).Error
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFail, "Failed to update notification!")
	}

	// Trả về kết quả thành công
	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, notification)
}
