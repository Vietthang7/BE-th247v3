package controllers

import (
	"intern_247/consts"
	"intern_247/repo"

	"github.com/gofiber/fiber/v2"
)

func CreateNotification(c *fiber.Ctx) error {
	token, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
	}

	var (
		entry repo.Notification
	)

	if err := c.BodyParser(&entry); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	entry.From = &token.ID
	err = entry.Create()
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.CreateFail, consts.CreateFailed)
	}
	return ResponseSuccess(c, fiber.StatusCreated, consts.CreateSuccess, entry)
}
