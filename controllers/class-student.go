package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"strings"
)

func ListStudentByEnrollmentPlan(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized", consts.ERROR_UNAUTHORIZED)
	}
	var (
		err     error
		entries []repo.StudentCanBeAddedIntoClass
		class   models.Class
	)
	class, err = repo.GetClassByIdAndCenterId(uuid.MustParse(c.Params("id")), *user.CenterId)
	switch {
	case err == nil:
		if entries, err = repo.FindStudentsCanBeAddedIntoClass(class,
			strings.TrimSpace(c.Query("search"))); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.GetFailed)
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, entries)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, "Error class not found: "+err.Error(), consts.DataNotFound)
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, "Error get class: "+err.Error(), consts.GetFailed)
	}
}
