package controllers

import (
	"errors"
	"fmt"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateDocument(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var entry repo.Document
	if err := c.BodyParser(&entry); err != nil {
		logrus.Error("Error invalid input: ", err.Error())
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}

	entry.CenterID = user.CenterId
	entry.CreatedBy = user.ID

	if err := entry.Create(); err != nil {
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.CreateFail, err.Error()), consts.CreateFailed)
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry)
}

func ReadDocument(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err   error
		entry repo.Document
	)

	err = entry.First("id = ? AND center_id = ?", []interface{}{c.Params("id"), user.CenterId})
	switch {
	case err == nil:
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error("Document not exist")
		return ResponseError(c, fiber.StatusNotFound, err.Error(), consts.DataNotFound)
	default:
		logrus.Error("Error finding Document: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}

func ListDocuments(c *fiber.Ctx) error {
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
		entry      repo.Document
		entries    repo.Documents
		pagination = consts.BindRequestTable(c, "created_at")
		query      = "center_id = ?"
		args       = []interface{}{token.CenterId}
	)

	if c.Query("categoryId") != "" {
		query += " AND category_id = ?"
		args = append(args, c.Query("categoryId"))
	}
	if c.Query("type") != "" {
		query += " AND type = ?"
		args = append(args, c.Query("type"))
	}
	if c.Query("startTime") != "" {
		var startTime time.Time
		if startTime, err = time.Parse("2006-01-02", c.Query("startTime")); err != nil {
			return ResponseError(c, fiber.StatusBadRequest, "Error query 'startTime'", consts.InvalidReqInput)
		}
		query += " AND created_at >= ?"
		args = append(args, startTime)
	}
	if c.Query("endTime") != "" {
		var endTime time.Time
		if endTime, err = time.Parse("2006-01-02", c.Query("endTime")); err != nil {
			return ResponseError(c, fiber.StatusBadRequest, "Error query 'endTime'", consts.InvalidReqInput)
		}
		query += " AND created_at <= ?"
		args = append(args, endTime)
	}

	if pagination.Search != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+pagination.Search+"%")
	}

	if entries, err = entry.Find(&pagination, query, args, "Creator", "Category"); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}

	pagination.Total = entry.Count(query, args)
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}

func DeleteDocument(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err   error
		entry repo.Document
	)

	err = entry.First("id = ? AND center_id = ?", []interface{}{c.Params("id"), user.CenterId})
	switch {
	case err == nil:
		if repo.CountLessonData("document_id = ?", []interface{}{c.Params("id")}) > 0 {
			return ResponseError(c, fiber.StatusBadRequest, "Oops! Something went wrong.",
				consts.DocumentCannotDelete)
		}

		if err = entry.Delete(); err != nil {
			logrus.Error("Error deleting Document: ", err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.DeleteFail, err.Error()), consts.DeletedFailed)
		}

		return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, nil)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error("Document not exist")
		return ResponseError(c, fiber.StatusNotFound, err.Error(), consts.DataNotFound)
	default:
		logrus.Error("Error finding Document: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}

func UpdateDocument(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err   error
		entry repo.Document
	)

	err = entry.First("id = ? AND center_id = ?", []interface{}{c.Params("id"), user.CenterId})
	switch {
	case err == nil:
		if repo.CountLessonData("document_id = ?", []interface{}{c.Params("id")}) > 0 {
			return ResponseError(c, fiber.StatusBadRequest, "Oops! Something went wrong.",
				consts.DocumentCannotUpdate)
		}

		if err = c.BodyParser(&entry); err != nil {
			logrus.Error("Error invalid input: ", err.Error())
			return ResponseError(c, fiber.StatusBadRequest,
				fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
		}

		if err = entry.Update(); err != nil {
			logrus.Error("Error updating Document: ", err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.DeleteFail, err.Error()), consts.DeletedFailed)
		}

		return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error("Document not exist")
		return ResponseError(c, fiber.StatusNotFound, err.Error(), consts.DataNotFound)
	default:
		logrus.Error("Error finding Document: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}
