package controllers

import (
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func CreateCenter(c *fiber.Ctx) error {
	var (
		entry models.Center
		err   error
	)
	if err = c.BodyParser(&entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, err.Error())
	}
	//// check center is exists with user id
	//_, err = repo.GetCenterIDByUserID(uuid.Nil) // Sử dụng `uuid.Nil` để bypass kiểm tra
	////center is exists
	//if err == nil {
	//	return ResponseError(c, fiber.StatusConflict, consts.DataExists, "Duplicate center")
	//}
	//entry.UserId = nil // Không gán `UserId` nếu không có thông tin người dùng
	if err = repo.CreateCenter(app.Database.DB, &entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.CreateFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry)
}
