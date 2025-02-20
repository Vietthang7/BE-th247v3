package controllers

import (
	"errors"
	"fmt"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"intern_247/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CreatingSpRequestInput struct {
	Title          string         `json:"title"`
	Content        string         `json:"content"`
	Type           int            `json:"type"` // Loại yêu cầu; internal/consts/config.go
	Metadata       datatypes.JSON `json:"metadata"`
	LeaveFromDate  time.Time      `json:"leave_from_date"`  // Ngày bắt đầu nghỉ
	LeaveUntilDate time.Time      `json:"leave_until_date"` // Ngày kết thúc nghỉ
	MakeUpClass    bool           `json:"make_up_class"`    // Yêu cầu học bù
	SubjectIds     []uuid.UUID    `json:"subject_ids"`
}

func CreateSpRequest(c *fiber.Ctx) error {
	student, ok := c.Locals("student").(repo.Student)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err   error
		input CreatingSpRequestInput
	)
	if err = c.BodyParser(&input); err != nil {
		logrus.Error("Error invalid input: ", err.Error())
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}

	if !utils.Contains([]int{consts.SRForLeave, consts.SRReserve, consts.SROther, consts.SRStopStudying}, input.Type) {
		logrus.Error("Error type")
		return ResponseError(c, fiber.StatusBadRequest, "Error type", consts.InvalidReqInput)
	}
	// val := validate.Validate(
	// 	&validators.StringIsPresent{
	// 		Name:  "title",
	// 		Field: input.Title,
	// 	}, &validators.StringIsPresent{
	// 		Name:  "content",
	// 		Field: input.Content,
	// 	},
	// )
	// switch input.Type {
	// case consts.SRForLeave:
	// 	val.Append(validate.Validate(
	// 		&validators.TimeIsPresent{
	// 			Name:  "leave_from_date",
	// 			Field: input.LeaveFromDate,
	// 		}, &validators.TimeIsPresent{
	// 			Name:  "leave_until_date",
	// 			Field: input.LeaveUntilDate,
	// 		},
	// 	))
	// case consts.SRReserve:
	// 	if len(input.SubjectIds) < 1 {
	// 		val.Add("", "Error zero length: subject_ids")
	// 	}
	// 	var studentSubject repo.StudentSubject
	// 	if studentSubject.Count("student_id = ? AND subject_id IN ?",
	// 		student.ID, input.SubjectIds) != int64(len(input.SubjectIds)) {
	// 		val.Add("", "Error invalid subject_ids")
	// 	}
	// }
	// if input.Type != consts.SRReserve && len(input.SubjectIds) > 0 {
	// 	val.Add("", "This type do not allow subject_ids")
	// }
	// if val.HasAny() {
	// 	return ResponseError(c, http.StatusBadRequest, val.Errors[val.Keys()[0]][0], consts.InvalidReqInput)
	// }

	entry := repo.SupportRequest{
		CenterID:    student.CenterId,
		CreatedBy:   student.ID,
		Title:       input.Title,
		Type:        input.Type,
		Content:     input.Content,
		Metadata:    input.Metadata,
		MakeUpClass: &input.MakeUpClass,
	}

	switch input.Type {
	case consts.SRForLeave:
		entry.LeaveFromDate = &input.LeaveFromDate
		entry.LeaveUntilDate = &input.LeaveUntilDate
	case consts.SRReserve:
		entry.SubjectIds = input.SubjectIds
	}

	if err = entry.Create(); err != nil {
		logrus.Error("Error creating SupportRequest: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.CreateFail, err.Error()), consts.CreateFailed)
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, nil)
}

func ListSpRequests(c *fiber.Ctx) error {
	var (
		err   error
		token repo.TokenData
	)
	if token, err = repo.GetTokenData(c); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusUnauthorized,
			fmt.Sprintf("Error unauthorized: %s", err.Error()), consts.ERROR_UNAUTHORIZED)
	}

	var (
		entry      repo.SupportRequest
		entries    repo.SupportRequests
		query      = "support_requests.center_id = ?"
		args       = []interface{}{token.CenterId}
		pagination = consts.BindRequestTable(c, "created_at")
	)

	if token.RoleId == consts.Student {
		query += " AND created_by = ?"
		args = append(args, token.ID)
	}
	if c.Query("type") != "" {
		query += " AND support_requests.type = ?"
		args = append(args, c.Query("type"))
	}
	if c.Query("agree") != "" {
		agree, _ := strconv.ParseBool(c.Query("agree"))
		query += " AND agree = ?"
		args = append(args, agree)
	}
	countQuery := query
	if c.Query("resolved") != "" {
		resolved, _ := strconv.ParseBool(c.Query("resolved"))
		if resolved {
			query += " AND agree IS NOT NULL"
		} else {
			query += " AND agree IS NULL"
		}
	}

	if token.RoleId == consts.Student {
		if pagination.Search != "" {
			query += " AND title LIKE ?"
			countQuery += " AND title LIKE ?"
			args = append(args, "%"+pagination.Search+"%")
		}
		entries, err = entry.Find(&pagination, query, args, "")
	} else {
		entries, err = entry.Find(&pagination, query, args, pagination.Search, "Creator", "Responder")
	}
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}

	pagination.Total = entry.Count(countQuery, args)
	var totalResolved int64
	countQuery += " AND agree IS NOT NULL"
	totalResolved = entry.Count(countQuery, args)

	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"count_by_resolve_status": fiber.Map{
			"resolved":   totalResolved,
			"unresolved": pagination.Total - totalResolved,
		},
		"data":       entries,
		"pagination": pagination,
	})
}

func UpdateSpRequests(c *fiber.Ctx) error {
	student, ok := c.Locals("student").(repo.Student)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err   error
		entry repo.SupportRequest
	)

	err = entry.First("id = ? AND center_id = ? AND agree IS NULL",
		[]interface{}{c.Params("id"), student.CenterId})
	switch {
	case err == nil:
		var input CreatingSpRequestInput
		if err = c.BodyParser(&input); err != nil {
			logrus.Error("Error invalid input: ", err.Error())
			return ResponseError(c, fiber.StatusBadRequest,
				fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
		}

		// if !utils.Contains([]int{consts.SRForLeave, consts.SRReserve, consts.SROther, consts.SRStopStudying}, input.Type) {
		// 	logrus.Error("Error type")
		// 	return ResponseError(c, fiber.StatusBadRequest, "Error type", consts.InvalidReqInput)
		// }

		// val := validate.Validate(
		// 	&validators.StringIsPresent{
		// 		Name:  "title",
		// 		Field: input.Title,
		// 	}, &validators.StringIsPresent{
		// 		Name:  "content",
		// 		Field: input.Content,
		// 	},
		// )

		// switch input.Type {
		// case consts.SRForLeave:
		// 	val.Append(validate.Validate(
		// 		&validators.TimeIsPresent{
		// 			Name:  "leave_from_date",
		// 			Field: input.LeaveFromDate,
		// 		}, &validators.TimeIsPresent{
		// 			Name:  "leave_until_date",
		// 			Field: input.LeaveUntilDate,
		// 		},
		// 	))
		// case consts.SRReserve:
		// 	if len(input.SubjectIds) < 1 {
		// 		val.Add("", "Error zero length: subject_ids")
		// 	}
		// 	var studentSubject repo.StudentSubject
		// 	if studentSubject.Count("student_id = ? AND subject_id IN ?",
		// 		student.ID, input.SubjectIds) != int64(len(input.SubjectIds)) {
		// 		val.Add("", "Error invalid subject_ids")
		// 	}
		// }
		// if input.Type != consts.SRReserve && len(input.SubjectIds) > 0 {
		// 	val.Add("", "This type does not allow subject_ids")
		// }
		// if val.HasAny() {
		// 	return ResponseError(c, http.StatusBadRequest, val.Errors[val.Keys()[0]][0], consts.InvalidReqInput)
		// }

		entry.Title = input.Title
		entry.Type = input.Type
		entry.Content = input.Content
		entry.Metadata = input.Metadata
		entry.MakeUpClass = &input.MakeUpClass
		switch input.Type {
		case consts.SRForLeave:
			entry.LeaveFromDate = &input.LeaveFromDate
			entry.LeaveUntilDate = &input.LeaveUntilDate
		case consts.SRReserve:
			entry.SubjectIds = input.SubjectIds
		}

		if err = entry.Update(); err != nil {
			logrus.Error("Error updating SupportRequest: ", err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.UpdateFail, err.Error()), consts.UpdateFailed)
		}

		return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error("SupportRequest does not exist")
		return ResponseError(c, fiber.StatusNotFound, err.Error(), consts.DataNotFound)
	default:
		logrus.Error("Error finding SupportRequest: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}

func ReadSpRequests(c *fiber.Ctx) error {
	var (
		err   error
		token repo.TokenData
	)
	if token, err = repo.GetTokenData(c); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusUnauthorized,
			fmt.Sprintf("Error unauthorized: %s", err.Error()), consts.ERROR_UNAUTHORIZED)
	}

	var entry repo.SupportRequest
	err = entry.First("id = ? AND center_id = ?",
		[]interface{}{c.Params("id"), token.CenterId}, "Responder", "Subjects")
	switch {
	case err == nil:
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error("SupportRequest does not exist")
		return ResponseError(c, fiber.StatusNotFound, err.Error(), consts.DataNotFound)
	default:
		logrus.Error("Error finding SupportRequest: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}

func DeleteSpRequests(c *fiber.Ctx) error {
	student, ok := c.Locals("student").(repo.Student)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err   error
		entry repo.SupportRequest
	)

	err = entry.First("id = ? AND center_id = ? AND agree IS NULL",
		[]interface{}{c.Params("id"), student.CenterId})
	switch {
	case err == nil:
		if err = entry.Delete(); err != nil {
			logrus.Error("Error deleting SupportRequest: ", err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.DeleteFail, err.Error()), consts.DeletedFailed)
		}

		return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, nil)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error("SupportRequest does not exist")
		return ResponseError(c, fiber.StatusNotFound, err.Error(), consts.DataNotFound)
	default:
		logrus.Error("Error finding SupportRequest: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}

func RespondSpRequests(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	type RespondSpRequestInput struct {
		Ids      uuid.UUIDs     `json:"ids"`
		Agree    *bool          `json:"agree"`
		Content  string         `json:"content"`
		Metadata datatypes.JSON `json:"metadata"`
	}

	var (
		err   error
		input RespondSpRequestInput
	)
	if err = c.BodyParser(&input); err != nil {
		logrus.Error("Error invalid input: ", err.Error())
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}

	// val := validate.Validate()
	// if len(input.Ids) < 1 {
	// 	logrus.Error("Error invalid input: invalid length of ids")
	// 	val.Add("", fmt.Sprintf("invalid length of 'ids'"))
	// }
	var entry repo.SupportRequest
	if entry.Count("center_id = ? AND id IN ? AND agree IS NULL",
		[]interface{}{user.CenterId, input.Ids}) < int64(len(input.Ids)) {
		logrus.Error("Error invalid input: invalid 'ids'")
		// val.Add("", fmt.Sprintf("invalid 'ids'"))
	}
	// if input.Agree == nil {
	// 	logrus.Error("Error invalid input: 'agree' can not be null")
	// 	val.Add("", fmt.Sprintf("'agree' can not be null"))
	// }
	// if val.HasAny() {
	// 	return ResponseError(c, http.StatusBadRequest, val.Errors[val.Keys()[0]][0], consts.InvalidReqInput)
	// }

	now := time.Now()
	entry.Agree = input.Agree
	entry.Metadata = input.Metadata
	entry.Response = input.Content
	entry.RespondedBy = &user.ID
	entry.ResolveDate = &now

	if err = entry.Respond(input.Ids); err != nil {
		logrus.Error("Error updating SupportRequest: ", err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.UpdateFail, err.Error()), consts.UpdateFailed)
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, nil)
}
