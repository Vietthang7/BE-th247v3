package controllers

import "github.com/gofiber/fiber/v2"

func ResponseError(c *fiber.Ctx, code int, message interface{}, err interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  false,
		"code":    code,
		"error":   err,
		"message": message,
	})
}
func ResponseSuccess(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  true,
		"code":    code,
		"data":    data,
		"message": message,
	})
}
