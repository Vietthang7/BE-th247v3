package controllers

import (
	"fmt"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"intern_247/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

type NewSubjectInput struct {
	Name string `json:"name"`

	Thumbnail     string         `json:"thumbnail"`
	IsActive      *bool          `json:"is_active"`
	FeeType       uint8          `json:"fee_type"` //1 - free, 2 - payment
	OriginFee     uint64         `json:"origin_fee"`
	DiscountFee   uint64         `json:"discount_fee"`
	Description   string         `json:"description"`
	TotalLessons  uint64         `json:"total_lessons"`
	InputRequire  string         `json:"input_require"`
	OutputRequire string         `json:"output_require"`
	Metadata      datatypes.JSON `json:"metadata,omitempty"`
	CategoryId    uuid.UUID      `json:"category_id"`
	TeacherIds    []uuid.UUID    `json:"teacher_ids"`
}
type SubjectUpdateInput struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Thumbnail     string         `json:"thumbnail"`
	IsActive      *bool          `json:"is_active"`
	FeeType       uint8          `json:"fee_type"` //1 - free, 2 - payment
	OriginFee     uint64         `json:"origin_fee"`
	DiscountFee   uint64         `json:"discount_fee"`
	Description   string         `json:"description"`
	TotalLessons  uint64         `json:"total_lessons"`
	InputRequire  string         `json:"input_require"`
	OutputRequire string         `json:"output_require"`
	Metadata      datatypes.JSON `json:"metadata,omitempty"`
	//CertIssuance  *bool          `json:"cert_issuance"`
	//CertificateId *uuid.UUID     `json:"certificate_id"`
	CategoryId *uuid.UUID  `json:"category_id"`
	TeacherIds []uuid.UUID `json:"teacher_ids"`
}

func CreateSubject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Lỗi quyền truy cập", consts.ERROR_PERMISSION_DENIED)
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusUnauthorized, "Lỗi quyền truy cập - trung tâm không xác định", consts.ERROR_PERMISSION_DENIED)
	}
	var (
		input   NewSubjectInput
		subject models.Subject
	)
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ", consts.InvalidReqInput)
	}
	if !utils.IsValidStrLen(input.Name, 250) {
		return ResponseError(c, fiber.StatusBadRequest, "Tên môn học không hợp lệ", consts.InvalidReqInput)
	}
	if input.OriginFee < input.DiscountFee {
		return ResponseError(c, fiber.StatusBadRequest, "Học phí giảm không thể lớn hơn học phí gốc", consts.ERROR_DATA_LONGER)
	}
	if input.FeeType != consts.FREE_SUBJECT && input.FeeType != consts.PAYMENT_SUBJECT {
		return ResponseError(c, fiber.StatusBadRequest, "Loại học phí không hợp lệ", consts.InvalidReqInput)
	}
	if input.Description == "" || input.InputRequire == "" || input.OutputRequire == "" {
		return ResponseError(c, fiber.StatusBadRequest, "Thông tin yêu cầu đầu vào và đầu ra không thể bỏ trống", consts.InvalidReqInput)
	}
	_, err := repo.GetCategoryByIdAndCenterId(input.CategoryId, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Danh mục không hợp lệ", consts.InvalidReqInput)
	}
	if input.TotalLessons < 1 {
		return ResponseError(c, fiber.StatusBadRequest, "Số buổi học phải lớn hơn 0", consts.InvalidReqInput)
	}
	teacherInputLen := len(input.TeacherIds)
	if teacherInputLen < 1 {
		return ResponseError(c, fiber.StatusBadRequest, "Cần chọn ít nhất một giảng viên", consts.InvalidReqInput)
	}
	teachers, _ := repo.GetTeachersByIdsAndCenterId(input.TeacherIds, *user.CenterId)
	if len(teachers) != teacherInputLen {
		return ResponseError(c, fiber.StatusBadRequest, "Không tìm thấy giảng viên", consts.ERROR_TEACHER_NOT_FOUND)
	}
	// Kiểm tra tên môn học đã tồn tại chưa
	if _, err = repo.GetSubjectByNameAndCenterId(input.Name, *user.CenterId); err == nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "Môn học đã tồn tại", consts.ERROR_SUBJECT_EXISTS)
	}

	if input.FeeType == consts.PAYMENT_SUBJECT {
		subject.OriginFee = input.OriginFee
		subject.DiscountFee = input.DiscountFee
	}
	subject.Name = input.Name
	subject.InputRequire = input.InputRequire
	subject.OutputRequire = input.OutputRequire
	subject.Description = input.Description
	subject.CategoryId = input.CategoryId
	subject.Metadata = input.Metadata
	subject.IsActive = input.IsActive
	subject.Thumbnail = input.Thumbnail
	subject.FeeType = input.FeeType
	subject.CenterId = *user.CenterId
	subject.TotalLessons = input.TotalLessons
	subject.CreatedBy = user.ID
	subject.Teachers = teachers
	subject.Code = utils.GenerateRandomCodeFormatByKey(consts.SUBJECT_CODE_PREFIX)
	newSubject, err := repo.CreateSubject(&subject)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Tạo môn học thất bại", consts.CreateFailed)
	}
	return ResponseSuccess(c, fiber.StatusCreated, "Thành công", newSubject)
}

func UpdateSubject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusUnauthorized, "Error Permission denied - center", consts.ERROR_PERMISSION_DENIED)
	}
	var (
		input         SubjectUpdateInput
		subjectUpdate *models.Subject
		_             uuid.UUID
	)
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	if !utils.IsValidStrLen(input.Name, 250) {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	if input.OriginFee < input.DiscountFee {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_DATA_LONGER)
	}
	if input.FeeType != consts.FREE_SUBJECT && input.FeeType != consts.PAYMENT_SUBJECT {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	if input.Description == "" || input.InputRequire == "" || input.OutputRequire == "" {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	if input.TotalLessons < 1 {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	// check name subject is exists
	if _, err := repo.GetSubjectByNameAndIdAndCenterId(input.Name, *user.CenterId, input.ID); err == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	subject, err := repo.GetSubjectByIdAndCenterId(input.ID, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
	}
	if subject.Code == "" {
		subject.Code = utils.GenerateRandomCodeFormatByKey(consts.SUBJECT_CODE_PREFIX)
		_, err = repo.UpdateSubject(&subject)
		if err != nil {
			return ResponseError(c, fiber.StatusInternalServerError, "Failed update code", consts.UpdateFailed)
		}
	}
	_ = subject.ID
	classes, _ := repo.GetClassesBySubjectIdAndCenterId(input.ID, *user.CenterId)

	//if input.CertificateId != nil {
	//	if subject.CertificateId != nil && *subject.CertificateId != *input.CertificateId {
	//
	//	}
	//}

	//check student subject registed
	students, _ := repo.GetStudentsBySubjectAndCenterId(input.ID, *user.CenterId)
	if (len(classes) > 0 || len(students) > 0) && (subject.DiscountFee != input.DiscountFee || subject.OriginFee != input.OriginFee || subject.TotalLessons != input.TotalLessons) {
		subject.ID = uuid.Nil
		subject.UpdatedAt = nil
	}
	if input.CategoryId != nil {
		_, err := repo.GetCategoryByIdAndCenterId(*input.CategoryId, *user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
		}
		subject.CategoryId = *input.CategoryId
	}
	teacherInputLen := len(input.TeacherIds)
	if teacherInputLen > 0 {
		teachers, err := repo.GetTeachersByIdsAndCenterId(input.TeacherIds, *user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
		}
		if len(teachers) != teacherInputLen {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_TEACHER_NOT_FOUND)
		}
		subject.Teachers = teachers
	}
	if input.FeeType == consts.PAYMENT_SUBJECT {
		subject.OriginFee = input.OriginFee
		subject.DiscountFee = input.DiscountFee
	}
	subject.Name = input.Name
	subject.InputRequire = input.InputRequire
	subject.OutputRequire = input.OutputRequire
	subject.Description = input.Description
	subject.TotalLessons = input.TotalLessons
	subject.Metadata = input.Metadata
	subject.IsActive = input.IsActive
	subject.Thumbnail = input.Thumbnail
	subject.FeeType = input.FeeType
	//subject.CertIssuance = input.CertIssuance
	if subject.ID == uuid.Nil {
		subject.ID = uuid.New()
		lessons, err1 := repo.GetAllLessonBySubjectIdAndCenterId(input.ID, *user.CenterId)
		for i := range lessons {
			lessons[i].ID = uuid.New()
			lessons[i].SubjectId = &subject.ID
			lessons[i].CreatedAt = nil
			lessons[i].UpdatedAt = nil
			for j := range lessons[i].Childrens {
				lessons[i].Childrens[j].ID = uuid.New()
				lessons[i].Childrens[j].SubjectId = &subject.ID
				lessons[i].Childrens[j].ParentId = &lessons[i].ID
				lessons[i].Childrens[j].CreatedAt = nil
				lessons[i].Childrens[j].UpdatedAt = nil
				for k := range lessons[i].Childrens[j].LessonDatas {
					lessons[i].Childrens[j].LessonDatas[k].ID = uuid.New()
					lessons[i].Childrens[j].LessonDatas[k].LessonId = lessons[i].Childrens[j].ID
					lessons[i].Childrens[j].LessonDatas[k].CreatedAt = nil
					lessons[i].Childrens[j].LessonDatas[k].UpdatedAt = nil
				}
			}
		}
		if err1 != nil {
			return ResponseError(c, fiber.StatusInternalServerError, "Failed get lesson", consts.UpdateFailed)
		}
		subject.Lessons = lessons
		//curriculums
		//curriculums, err := repo.GetCurriculumsBySubjectIdAndCenterId(subjectId, *user.CenterId)
		//if err != nil {
		//	return ResponseError(c, fiber.StatusInternalServerError, "Failed get relation curri", consts.UpdateFailed)
		//}
		//for i := range curriculums {
		//	for j := range curriculums[i].Subjects {
		//		if curriculums[i].Subjects[j].ID == subjectId {
		//			curriculums[i].Subjects[j].ID = subject.ID
		//			curriculums[i].ID = uuid.Nil
		//			curriculums[i].UpdatedAt = nil
		//		}
		//	}
		//}
		subjectUpdate, err = repo.CreateSubject(&subject)
		if err != nil {
			return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.UpdateFailed)
		}
		//_, err = repo.CreateCurriculums(curriculums)
		//if err != nil {
		//	return ResponseError(c, fiber.StatusNoContent, "failed when update curriculums", consts.CreateFailed)
		//}
	} else {
		subjectUpdate, err = repo.UpdateSubject(&subject)
	}
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed update", consts.UpdateFailed)
	}
	return ResponseSuccess(c, fiber.StatusOK, "success", subjectUpdate)
}

func DeleteSubject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusUnauthorized, "Error Permission denied - center", consts.ERROR_PERMISSION_DENIED)
	}

	// Lấy ID từ params
	idParam := c.Params("id")
	subjectID, err := uuid.Parse(idParam)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid ID format", consts.InvalidReqInput)
	}

	_, err = repo.GetSubjectByIdAndCenterId(subjectID, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid ID", consts.InvalidReqInput)
	}

	classes, err := repo.GetClassesBySubjectIdAndCenterId(subjectID, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Error retrieving classes", consts.ERROR_INTERNAL_SERVER_ERROR)
	}
	if len(classes) > 0 {
		return ResponseError(c, fiber.StatusBadRequest, "Cannot delete subject with associated classes", consts.ERROR_CAN_NOT_DELETE_SUBJECT_HAS_CLASS)
	}

	// Kiểm tra xem có học sinh nào đã đăng ký môn học này không
	students, err := repo.GetStudentsBySubjectAndCenterId(subjectID, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.ERROR_INTERNAL_SERVER_ERROR)
	}
	if len(students) > 0 {
		return ResponseError(c, fiber.StatusBadRequest, "Cannot delete subject with registered students", consts.ERROR_CAN_NOT_DELETE_SUBJECT_HAS_STUDENT)
	}

	row, err := repo.DeleteSubjectByIdAndCenterId(subjectID, *user.CenterId)
	if row < 1 || err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Error deleting subject", consts.ERROR_INTERNAL_SERVER_ERROR)
	}

	_, _ = repo.DeleteLessonBySubjectIdAndCenterId(subjectID, *user.CenterId)
	return ResponseSuccess(c, fiber.StatusOK, "Success", row)
}

func GetDetailSubject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusUnauthorized, "Error Permission denied - center", consts.ERROR_PERMISSION_DENIED)
	}
	q := new(consts.Query)
	if err := c.QueryParser(q); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	if q.ID == uuid.Nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	fmt.Println("ok")
	subject, err := repo.GetDetailSubjectByIdAndCenterId(q.ID, *user.CenterId)
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
	}
	return ResponseSuccess(c, fiber.StatusOK, "Success", subject)
}

func GetListSubjects(c *fiber.Ctx) error {
	user, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusUnauthorized, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}
	query := new(consts.Query)
	if err := c.QueryParser(query); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, err.Error())
	}
	subjects, pagination, err := repo.GetListSubjectsByCenterId(*query, user)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "invalid", err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, "success", fiber.Map{"subjects": subjects, "pagination": pagination})
}
func GetAllSubject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusUnauthorized, "Error Permission denied - center", consts.ERROR_PERMISSION_DENIED)
	}
	query := new(consts.Query)
	if err := c.QueryParser(query); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, err.Error())
	}
	subjects, err := repo.GetAllSubjectByCenterId(*query, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "invalid", consts.DataNotFound)
	}
	return ResponseSuccess(c, fiber.StatusOK, "success", subjects)
}
