package controllers

import (
	"fmt"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"intern_247/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type NewClassInput struct {
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	Type        int        `json:"type"`
	Description string     `json:"description"` //number and character
	BranchId    uuid.UUID  `json:"branch_id"`
	//PlanId       uuid.UUID  `json:"plan_id"` //ke hoach tuyen sinh
	//CurriculumId *uuid.UUID `json:"curriculum_id"`
	CategoryId uuid.UUID `json:"category_id"`
	SubjectId  uuid.UUID `json:"subject_id"`
	// CuratorId    uuid.UUID      `json:"curator_id"`
	GroupUrl string         `json:"group_url"`
	Metadata datatypes.JSON `json:"metadata"` // Đường dẫn nhóm học
}
type ClassUpdateInput struct {
	Id          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	Type        int        `json:"type"`
	Description string     `json:"description"` //number and character
	BranchId    uuid.UUID  `json:"branch_id"`
	//PlanId       uuid.UUID  `json:"plan_id"` //ke hoach tuyen sinh
	//CurriculumId *uuid.UUID `json:"curriculum_id"`
	CategoryId uuid.UUID `json:"category_id"`
	SubjectId  uuid.UUID `json:"subject_id"`
	// CuratorId    uuid.UUID      `json:"curator_id"`
	GroupUrl string         `json:"group_url"`
	Metadata datatypes.JSON `json:"metadata"`
}

func CreateClass(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}
	var (
		input  NewClassInput
		class  models.Class
		active bool
	)
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.InvalidReqInput)
	}
	codeLength := len(input.Code)
	if input.Name == "" {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_CLASS_NAME_REQUIRED)
	}
	if !utils.IsValidStrLen(input.Name, 250) {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_DATA_MAX_SIZE_250)
	}
	if codeLength < consts.CLASS_CODE_MIN_SIZE {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_CLASS_CODE_MIN_SIZE_6)
	}
	if codeLength > consts.CLASS_CODE_MAX_SIZE {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_CLASS_CODE_MAX_SIZE_20)
	}
	if utils.ContainSpecialCharacter(input.Code) {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_CLASS_CODE_NOT_SUPPORT_SPECIAL_CHARACTER)
	}
	if input.EndAt != nil {
		if time.Time(input.StartAt).After(time.Time(*input.EndAt)) {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_START_TIME_MUST_SMALLER_THAN_END_TIME)
		}
	}
	class, err := repo.GetClassByCodeAndCenterId(input.Code, uuid.Nil, *user.CenterId)
	if err == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_CLASS_CODE_DUPLICATED)
	}
	if input.Type != consts.CLASS_TYPE_ONLINE && input.Type != consts.CLASS_TYPE_OFFLINE && input.Type != consts.CLASS_TYPE_HYBRID {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_TYPE_NOT_FOUND)
	}
	//TODO - CHECK PLAN apply
	active = true
	subject, err := repo.GetSubjectByIdAndCenterId(input.SubjectId, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.ERROR_SUBJECT_NOT_FOUND, consts.InvalidReqInput)
	}
	if subject.IsActive != nil && !*subject.IsActive {
		return ResponseError(c, fiber.StatusBadRequest, "Subject is not active", consts.InvalidReqInput)
	}
	//if input.CurriculumId != nil {
	//	isExist := repo.IsExistCurriculumInCenter(*input.CurriculumId, *user.CenterId, &active)
	//	if !isExist {
	//		return ResponseError(c, fiber.StatusBadRequest, "Curriculum not active", consts.InvalidReqInput)
	//	}
	//}
	isExist := repo.IsExistBranchInCenter(input.BranchId, *user.CenterId, &active)
	if !isExist {
		return ResponseError(c, fiber.StatusBadRequest, "Branch not active or not exists", consts.InvalidReqInput)
	}
	class.Name = input.Name
	class.Code = input.Code
	class.Description = input.Description
	class.GroupUrl = input.GroupUrl
	class.CreatedBy = user.ID
	//class.PlanId = input.PlanId
	class.BranchId = &input.BranchId
	class.SubjectId = input.SubjectId
	//class.CurriculumId = input.CurriculumId
	class.StartAt = &input.StartAt
	class.EndAt = input.EndAt
	class.CenterId = *user.CenterId
	class.Type = input.Type
	class.TotalLessons = subject.TotalLessons
	class.Metadata = input.Metadata
	newClass, err := repo.CreateClass(&class)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.ERROR_INTERNAL_SERVER_ERROR)
	}
	return ResponseSuccess(c, fiber.StatusCreated, "Success", newClass)
}
func GetDetailClass(c *fiber.Ctx) error {
	token, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	isChildren := c.QueryBool("children")
	class, err := repo.GetDetailClassByIdAndCenterId(id, token.CenterId, token, isChildren)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
	}
	if c.Query("clone") == "true" {
		class.Name = fmt.Sprintf("%s (%s)", class.Name, strconv.Itoa(class.TotalChild+1))
	}
	return ResponseSuccess(c, fiber.StatusOK, "success", class)
}
func UpdateClass(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}
	var (
		input           ClassUpdateInput
		class           models.Class
		active          bool
		isChangeType    bool
		isChangeSubject bool
	)
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.InvalidReqInput)
	}
	class, err := repo.GetClassByIdAndCenterId(input.Id, *user.CenterId)
	codeLength := len(input.Code)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Not found", consts.DataNotFound)
	}
	if class.Status == consts.CLASS_CANCELED {
		return ResponseError(c, fiber.StatusBadRequest, "Lớp học đã được hủy", consts.InvalidReqInput)
	}
	if class.StartAt != nil {
		if time.Now().After(*class.StartAt) && len(class.StudentsClasses) > 0 {
			if input.Description != class.Description || input.GroupUrl != class.GroupUrl || (input.EndAt != nil && class.EndAt.Before(*input.EndAt)) || input.Metadata != nil {
				goto UpdateMoreInfo
			}
			return ResponseError(c, fiber.StatusBadRequest, "Not found", consts.ERROR_CLASS_CAN_NOT_UPDATE_INPROGRESS)
		}
	}
	if input.Name == "" {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_CLASS_NAME_REQUIRED)
	}
	if !utils.IsValidStrLen(input.Name, 250) {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_DATA_MAX_SIZE_250)
	}
	if input.EndAt != nil {
		if time.Time(input.StartAt).After(time.Time(*input.EndAt)) {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_START_DAY_MUST_SMALLER_THAN_END_DAY)
		}
	}
	if codeLength < consts.CLASS_CODE_MIN_SIZE {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_CLASS_CODE_MIN_SIZE_6)
	}
	if codeLength > consts.CLASS_CODE_MAX_SIZE {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_CLASS_CODE_MAX_SIZE_20)
	}
	if utils.ContainSpecialCharacter(input.Code) {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_CLASS_CODE_NOT_SUPPORT_SPECIAL_CHARACTER)
	}
	if oldClass, err := repo.GetClassByCodeAndCenterId(input.Code, input.Id, *user.CenterId); err == nil {
		if oldClass.ID != uuid.Nil {
			return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_CLASS_CODE_DUPLICATED)
		}
	}
	if input.Type != consts.CLASS_TYPE_HYBRID && input.Type != consts.CLASS_TYPE_OFFLINE && input.Type != consts.CLASS_TYPE_ONLINE {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_TYPE_NOT_FOUND)
	}
	if class.Type != input.Type {
		isChangeType = true
	}
	//TODO - CHECK PLAN apply
	active = true
	if subject, err := repo.GetSubjectByIdAndCenterId(input.SubjectId, *user.CenterId); err == nil {
		class.TotalLessons = subject.TotalLessons
	} else {
		return ResponseError(c, fiber.StatusBadRequest, "subject not found", consts.InvalidReqInput)
	}
	if input.SubjectId != class.SubjectId {
		isChangeSubject = true
	}
	fmt.Println(*user.CenterId)
	if isExist := repo.IsExistBranchInCenter(input.BranchId, *user.CenterId, &active); !isExist {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid - b", consts.InvalidReqInput)
	}
	class.Code = input.Code
	class.Name = input.Name
	class.BranchId = &input.BranchId
	class.SubjectId = input.SubjectId
	class.StartAt = &input.StartAt
	class.EndAt = input.EndAt
	//class.PlanId = input.PlanId
	class.Type = input.Type

	// Cập nhật trạng thái của lớp học dựa trên thời gian
	if class.StartAt != nil && time.Now().Before(*class.StartAt) {
		class.Status = consts.CLASS_COMING_SOON
	} else if class.EndAt != nil && time.Now().After(*class.EndAt) {
		class.Status = consts.CLASS_FINISHED
	} else {
		class.Status = consts.CLASS_IN_PROGRESS
	}

UpdateMoreInfo:
	class.Description = input.Description
	class.GroupUrl = input.GroupUrl
	class.EndAt = input.EndAt
	// class.CuratorId = input.CuratorId
	class.Metadata = input.Metadata
	newClass, err := repo.UpdateClassByIdAndCenterId(&class, isChangeType, isChangeSubject)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Error", consts.ERROR_INTERNAL_SERVER_ERROR)
	}
	return ResponseSuccess(c, fiber.StatusOK, "Update success", newClass)
}
func GetListClasses(c *fiber.Ctx) error {
	tokenData, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "invalid", consts.ERROR_PERMISSION_DENIED)
	}
	var (
		q          consts.Query
		Subjectids uuid.UUIDs
	)
	//passPercent := float64(consts.PASS_CONDITION / 100)
	type CustomListClasses struct {
		models.Class
		CorrectTotal int `json:"correct_total"`
	}
	//var classesCustom []CustomListClasses
	if err := c.QueryParser(&q); err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "InvalidInput", consts.InvalidReqInput)
	}
	if q.StudentId != "" {
		if student_id, err := uuid.Parse(q.StudentId); err != nil {
			return ResponseError(c, fiber.StatusInternalServerError, "InvalidInput "+err.Error(), consts.InvalidReqInput)
		} else {
			_, Subjectids, _ = repo.CountUnclassifiedsubjects(student_id)
		}
	}
	classes1, _ := repo.GetAllClasses(tokenData.CenterId)
	_ = repo.SaveAllStatusClasses(classes1)
	classes, pagination, overview, err := repo.GetListClassesByQueryAndCenterId(q, tokenData.CenterId, tokenData)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Error Permission denied", err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, "Success", fiber.Map{"classes": classes, "pagination": pagination, "subjects": Subjectids, "overview": overview})
}

func DeleteClass(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Lỗi quyền truy cập", consts.ERROR_PERMISSION_DENIED)
	}

	classId := c.Params("id") // Lấy ID lớp học từ tham số URL

	// Chuyển đổi ID lớp học từ tham số URL thành UUID
	classUUID, err := uuid.Parse(classId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "ID lớp học không hợp lệ", consts.ERROR_INVALID_CLASS_ID)
	}

	// Kiểm tra xem lớp học có tồn tại và thuộc trung tâm hiện tại không
	_, err = repo.GetClassByIdAndCenterId(classUUID, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusNotFound, "Lớp học không tìm thấy", consts.ERROR_CLASS_NOT_FOUND)
	}

	// Thực hiện xóa lớp học mà không kiểm tra trạng thái IsActive
	err = repo.DeleteClass(classUUID)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Lỗi khi xóa lớp học", consts.ERROR_INTERNAL_SERVER_ERROR)
	}

	// Trả về phản hồi thành công
	return ResponseSuccess(c, fiber.StatusOK, "Xóa lớp học thành công", nil)
}

func CanceledClass(c *fiber.Ctx) error {
	// Lấy thông tin người dùng từ context
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Lỗi quyền truy cập", consts.ERROR_PERMISSION_DENIED)
	}

	// Lấy ID lớp học từ tham số URL
	classId := c.Params("id")
	classUUID, err := uuid.Parse(classId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "ID lớp học không hợp lệ", consts.ERROR_CLASS_NOT_FOUND)
	}

	// Lấy centerId từ thông tin người dùng (giả sử rằng người dùng có quyền trên một center nhất định)
	centerId := *user.CenterId // Giả sử user.CenterId là kiểu uuid.UUID và có giá trị

	// Tìm lớp học theo classId và centerId sử dụng hàm GetClassByIdAndCenterId
	class, err := repo.GetClassByIdAndCenterId(classUUID, centerId)
	if err != nil {
		return ResponseError(c, fiber.StatusNotFound, "Lớp học không tìm thấy hoặc không thuộc trung tâm này", consts.ERROR_CLASS_NOT_FOUND)
	}

	// Kiểm tra xem lớp học đã bị hủy chưa
	if class.Status == consts.CLASS_CANCELED {
		return ResponseError(c, fiber.StatusBadRequest, "Lớp học đã bị hủy rồi", consts.ERROR_CLASS_CANCELED_OR_FINISHED_NOT_CANCEL)
	}

	// Lấy lý do hủy từ request body
	var requestBody struct {
		CancelReason string `json:"cancel_reason"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Dữ liệu không hợp lệ", consts.InvalidReqInput)
	}

	// Cập nhật trạng thái lớp học thành CLASS_CANCELED và lưu lý do hủy
	class.Status = consts.CLASS_CANCELED
	class.CancelReason = requestBody.CancelReason

	// Cập nhật lớp học trong cơ sở dữ liệu
	err = repo.UpdateClassStatus(&class)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Lỗi khi cập nhật trạng thái lớp học", consts.ERROR_INTERNAL_SERVER_ERROR)
	}

	// Trả về phản hồi thành công
	return ResponseSuccess(c, fiber.StatusOK, "Lớp học đã được hủy thành công", class)
}
