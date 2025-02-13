package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
)

func CreatePermission(c *fiber.Ctx) error {
	var (
		entry models.Permission
		err   error
	)
	if err = c.BodyParser(&entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, err.Error())
	}
	if err = repo.CreatePermission(app.Database.DB, &entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.CreateFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry)
}

func CreatePermissionTag(c *fiber.Ctx) error {
	var (
		entry models.PermissionTag
		err   error
	)
	if err = c.BodyParser(&entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, err.Error())
	}
	if err = repo.CreatePermissionTag(app.Database.DB, &entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.CreateFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry)
}
func ListPermissionTags(c *fiber.Ctx) error {
	var (
		err     error
		entries []models.PermissionTag
		DB      = app.Database.DB
	)
	if entries, err = repo.FindPermissionTags(DB.Order(consts.DescCreatedAt).Where("parent_tag_id IS NULL")); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
	var permissionsGoingWith datatypes.JSON
	permissionsGoingWith = []byte(consts.PermissionsGoingWith)
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":                   entries,
		"permissions_going_with": permissionsGoingWith,
	})
}
func ReadPermissionTag(c *fiber.Ctx) error {
	entry, err := repo.FirstPermissionTag(app.Database.DB, c.Params("id"))
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
