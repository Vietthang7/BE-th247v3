package controllers

import (
	"errors"
	"fmt"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateDocsCategory(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err   error
		entry repo.DocsCategory
	)

	if err = c.BodyParser(&entry); err != nil {
		logrus.Error("Error invalid input: ", err.Error())
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}

	// Parse body vào struct entry
	if err = c.BodyParser(&entry); err != nil {
		logrus.Error("Error invalid input: ", err.Error())
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}

	if entry.Count("center_id = ? AND name = ?", []interface{}{user.CenterId, strings.TrimSpace(entry.Name)}) > 0 {
		logrus.Error("DocsCategory is exist")
		return ResponseError(c, fiber.StatusConflict,
			"Tên danh mục tài liệu đã tồn tại", consts.DocsCategoryExistence)
	}

	entry.CenterID = user.CenterId
	entry.CreatedBy = user.ID
	if err = entry.Create(); err != nil {
		logrus.Error("Error creating DocsCategory: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.CreateFail, err.Error()), consts.CreateFailed)
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, nil)
}

func ReadDocsCategory(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err   error
		entry repo.DocsCategory
	)

	err = entry.First("id = ? AND center_id = ?", []interface{}{c.Params("id"), user.CenterId})
	switch {
	case err == nil:
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error("DocsCategory not exist")
		return ResponseError(c, fiber.StatusNotFound, err.Error(), consts.DataNotFound)
	default:
		logrus.Error("Error finding DocsCategory: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}

func ListDocsCategories(c *fiber.Ctx) error {
	var (
		err   error
		token repo.TokenData
	)
	if token, err = repo.GetTokenData(c); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusUnauthorized,
			fmt.Sprintf("Error unauthorized: %s", err.Error()), consts.ERROR_PERMISSION_DENIED)
	}

	var (
		entry      repo.DocsCategory
		entries    repo.DocsCategories
		pagination = consts.BindRequestTable(c, "created_at")
		query      = "center_id = ?"
		args       = []interface{}{token.CenterId}
	)

	if c.Query("active") != "" {
		isActive, _ := strconv.ParseBool(c.Query("active"))
		query += " AND is_active = ?"
		args = append(args, isActive)
	}

	if entries, err = entry.Find(&pagination, query, args, "Creator"); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}

	if len(entries) > 0 {
		var docs repo.DocsCategory //
		for i, v := range entries {
			entries[i].TotalDocs = docs.Count("category_id = ?", []interface{}{v.ID})
		}
	}

	pagination.Total = entry.Count(query, args)
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}

func UpdateDocsCategory(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err   error
		entry repo.DocsCategory
	)

	err = entry.First("id = ? AND center_id = ?", []interface{}{c.Params("id"), user.CenterId})
	switch {
	case err == nil:
		if err = c.BodyParser(&entry); err != nil {
			logrus.Error("Error invalid input: ", err.Error())
			return ResponseError(c, fiber.StatusBadRequest,
				fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
		}

		if entry.Count("center_id = ? AND name = ? AND id <> ?",
			[]interface{}{user.CenterId, strings.TrimSpace(entry.Name), entry.ID}) > 0 {
			logrus.Error("DocsCategory is exist")
			return ResponseError(c, fiber.StatusConflict,
				"Tên Danh mục tài liệu đã tồn tại", consts.DocsCategoryExistence)
		}

		if err = entry.Update(); err != nil {
			logrus.Error("Error updating customer source: ", err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.DeleteFail, err.Error()), consts.DeletedFailed)
		}

		return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error("Customer source not exist")
		return ResponseError(c, fiber.StatusNotFound, err.Error(), consts.DataNotFound)
	default:
		logrus.Error("Error finding customer source: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}

func DeleteDocsCategory(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err   error
		entry repo.DocsCategory
	)
	err = entry.First("id = ? AND center_id = ?", []interface{}{c.Params("id"), user.CenterId})
	switch {
	case err == nil:
		var document repo.DocsCategory //
		if document.Count("category_id = ?", []interface{}{entry.ID}) > 0 {
			logrus.Error("Error category has documents")
			return ResponseError(c, fiber.StatusBadRequest,
				"Error category has documents", consts.DocsCategoryIsAssigned)
		}

		if err = entry.Delete(); err != nil {
			logrus.Error("Error deleting customer source: ", err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.DeleteFail, err.Error()), consts.DeletedFailed)
		}

		return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, nil)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error("Customer source not exist")
		return ResponseError(c, fiber.StatusNotFound, err.Error(), consts.DataNotFound)
	default:
		logrus.Error("Error finding customer source: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}
