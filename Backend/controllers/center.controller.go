package controllers

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ReadCenter(c *fiber.Ctx) error {
	var (
		err   error
		token repo.TokenData
	)
	if token, err = repo.GetTokenData(c); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusUnauthorized,
			fmt.Sprintf("Error unauthorized: %s", err.Error()), consts.ERROR_UNAUTHORIZED)
	}
	//jsonData, _ := json.MarshalIndent(token, "", "  ")
	//fmt.Println(string(jsonData))

	logrus.Debug(c.GetReqHeaders())
	entry, err := repo.ReadCenter(app.Database.DB.Where("id = ?", token.CenterId))
	switch {
	case err == nil:
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ResponseError(c, fiber.StatusNotFound, consts.NotFound, err.Error())
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
}

func CreateCenter(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}
	var (
		entry models.Center
		err   error
	)
	if err = c.BodyParser(&entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, err.Error())
	}
	// check center is exists with user id
	_, err = repo.GetCenterIDByUserID(user.ID)
	//center is exists
	if err == nil {
		return ResponseError(c, fiber.StatusConflict, consts.DataExists, "Duplicate center")
	}
	entry.UserId = &user.ID
	if err = repo.CreateCenter(app.Database.DB, &entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.CreateFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry)
}
func UpdateCenter(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}
	DB := app.Database.DB.Where("id = ?", user.CenterId)
	entry, err := repo.ReadCenter(DB)
	switch {
	case err == nil:
		if err = c.BodyParser(&entry); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, err.Error())
		}
		if user.RoleId == consts.CenterOwner {
			entry.UserId = &user.ID
		}
		if err = repo.UpdateCenter(DB, &entry); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFail, err.Error())
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, consts.NotFound, err.Error())
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
}
