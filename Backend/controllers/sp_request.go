package controllers

import (
	"fmt"
	"intern_247/consts"
	"intern_247/repo"
	"intern_247/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
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

// func CreateSpRequest(c *fiber.Ctx) error {
// 	student, ok := c.Locals("student").(repo.Student)
// 	if !ok {
// 		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
// 	}

// 	var (
// 		err   error
// 		input CreatingSpRequestInput
// 	)
// 	if err = c.BodyParser(&input); err != nil {
// 		logrus.Error("Error invalid input: ", err.Error())
// 		return ResponseError(c, fiber.StatusBadRequest,
// 			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
// 	}

// 	if !utils.Contains([]int{consts.SRForLeave, consts.SRReserve, consts.SROther, consts.SRStopStudying}, input.Type) {
// 		logrus.Error("Error type")
// 		return ResponseError(c, fiber.StatusBadRequest, "Error type", consts.InvalidReqInput)
// 	}
// 	// val := validate.Validate(
// 	// 	&validators.StringIsPresent{
// 	// 		Name:  "title",
// 	// 		Field: input.Title,
// 	// 	}, &validators.StringIsPresent{
// 	// 		Name:  "content",
// 	// 		Field: input.Content,
// 	// 	},
// 	// )
// 	// switch input.Type {
// 	// case consts.SRForLeave:
// 	// 	val.Append(validate.Validate(
// 	// 		&validators.TimeIsPresent{
// 	// 			Name:  "leave_from_date",
// 	// 			Field: input.LeaveFromDate,
// 	// 		}, &validators.TimeIsPresent{
// 	// 			Name:  "leave_until_date",
// 	// 			Field: input.LeaveUntilDate,
// 	// 		},
// 	// 	))
// 	// case consts.SRReserve:
// 	// 	if len(input.SubjectIds) < 1 {
// 	// 		val.Add("", "Error zero length: subject_ids")
// 	// 	}
// 	// 	var studentSubject repo.StudentSubject
// 	// 	if studentSubject.Count("student_id = ? AND subject_id IN ?",
// 	// 		student.ID, input.SubjectIds) != int64(len(input.SubjectIds)) {
// 	// 		val.Add("", "Error invalid subject_ids")
// 	// 	}
// 	// }
// 	// if input.Type != consts.SRReserve && len(input.SubjectIds) > 0 {
// 	// 	val.Add("", "This type do not allow subject_ids")
// 	// }
// 	// if val.HasAny() {
// 	// 	return ResponseError(c, http.StatusBadRequest, val.Errors[val.Keys()[0]][0], consts.InvalidReqInput)
// 	// }

// 	entry := repo.SupportRequest{
// 		CenterID:    student.CenterId,
// 		CreatedBy:   student.ID,
// 		Title:       input.Title,
// 		Type:        input.Type,
// 		Content:     input.Content,
// 		Metadata:    input.Metadata,
// 		MakeUpClass: &input.MakeUpClass,
// 	}

// 	switch input.Type {
// 	case consts.SRForLeave:
// 		entry.LeaveFromDate = &input.LeaveFromDate
// 		entry.LeaveUntilDate = &input.LeaveUntilDate
// 	case consts.SRReserve:
// 		entry.SubjectIds = input.SubjectIds
// 	}

// 	if err = entry.Create(); err != nil {
// 		logrus.Error("Error creating SupportRequest: ", err.Error())
// 		return ResponseError(c, fiber.StatusInternalServerError,
// 			fmt.Sprintf("%s: %s", consts.CreateFail, err.Error()), consts.CreateFailed)
// 	}

// 	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, nil)
// }

func CreateSpRequest(c *fiber.Ctx) error {
	var (
		err   error
		input CreatingSpRequestInput
	)

	// Parse JSON input từ request body
	if err = c.BodyParser(&input); err != nil {
		logrus.Error("Error invalid input: ", err.Error())
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}

	// Kiểm tra type hợp lệ
	if !utils.Contains([]int{consts.SRForLeave, consts.SRReserve, consts.SROther, consts.SRStopStudying}, input.Type) {
		logrus.Error("Error type")
		return ResponseError(c, fiber.StatusBadRequest, "Error type", consts.InvalidReqInput)
	}

	// Tạo yêu cầu mới (Không ràng buộc vào student)
	entry := repo.SupportRequest{
		Title:       input.Title,
		Type:        input.Type,
		Content:     input.Content,
		Metadata:    input.Metadata,
		MakeUpClass: &input.MakeUpClass,
	}

	// Nếu là yêu cầu nghỉ, thêm ngày nghỉ
	if input.Type == consts.SRForLeave {
		entry.LeaveFromDate = &input.LeaveFromDate
		entry.LeaveUntilDate = &input.LeaveUntilDate
	}

	// Nếu là yêu cầu bảo lưu, thêm danh sách môn học
	if input.Type == consts.SRReserve {
		entry.SubjectIds = input.SubjectIds
	}

	// Tạo bản ghi trong DB
	if err = entry.Create(); err != nil {
		logrus.Error("Error creating SupportRequest: ", err.Error())
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.CreateFail, err.Error()), consts.CreateFailed)
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, nil)
}
