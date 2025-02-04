package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"time"
)

func CreateStudent(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}
	var (
		err   error
		entry repo.Student
	)
	if err = c.BodyParser(&entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}
	// data validation
	if entry.Type == consts.Official || entry.Type == consts.Trial {
		var existence repo.LoginInfo
		if err = existence.First("email = ? OR phone = ?", []interface{}{entry.Email, entry.Phone}); err == nil {
			var errExist = ""
			if entry.Email != "" && existence.Email == entry.Email {
				errExist += "Email đã tồn tại."
				return ResponseError(c, fiber.StatusConflict, errExist, consts.EmailDuplication)
			}
			if entry.Phone != "" && existence.Phone == entry.Phone {
				errExist += "Số điện thoại đã tồn tại."
				return ResponseError(c, fiber.StatusConflict, errExist, consts.PhoneDuplication)
			}
		}
	} else {
		var (
			existence repo.Student
			errExist  = ""
		)
		if err = existence.First("email = ? OR phone = ?", []interface{}{entry.Email, entry.Phone}); err == nil {
			if entry.Email != "" && existence.Email == entry.Email {
				errExist += "Email đã tồn tại"
				return ResponseError(c, fiber.StatusConflict, errExist, consts.EmailDuplication)
			}
			if entry.Phone != "" && existence.Phone == entry.Phone {
				errExist += "Số điện thoại đã tồn tại"
				return ResponseError(c, fiber.StatusConflict, errExist, consts.PhoneDuplication)
			}
		}
		var userExistence repo.User
		if err = userExistence.First("email = ? OR phone = ?", []interface{}{entry.Email, entry.Phone}); err == nil {
			if entry.Email != "" && existence.Email == entry.Email {
				errExist += "Email đã tồn tại."
				return ResponseError(c, fiber.StatusConflict, errExist, consts.EmailDuplication)
			}
			if entry.Phone != "" && existence.Phone == entry.Phone {
				errExist += "Số điện thoại đã tồn tại."
				return ResponseError(c, fiber.StatusConflict, errExist, consts.PhoneDuplication)
			}
		}
	}
	entry.CenterId = *user.CenterId
	entry.BranchId = user.BranchId
	if entry.Type == consts.Official {
		entry.IsOfficialAt = time.Now()
	}
	if err = entry.Create(); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.CreateFail, err.Error()), consts.CreateFailed)
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry.ID)
}
func ReadStudent(c *fiber.Ctx) error {
	_, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, consts.InvalidInput, "Permission denied!")
	}
	var (
		err   error
		entry repo.Student
	)
	if err = entry.First("id = ?", []interface{}{c.Params("id")}, "Province", "District", "CustomerSource", "ContactChannel"); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusForbidden, consts.GetFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		//"can_not_update": canNotUpdate,
		"entry": entry,
	})
}

func UpdateStudent(c *fiber.Ctx) error {
	//_, ok := c.Locals("user").(models.User)
	//if !ok {
	//	return ResponseError(c, fiber.StatusForbidden, "Permission denied!", consts.Forbidden)
	//}
	//var (
	//	err   error
	//	entry repo.Student
	//)
	//err = entry.First("id",[]interface{}{c.Params("id")})
	//switch  {
	//case err == nil:
	//
	//}
	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, nil)
}
