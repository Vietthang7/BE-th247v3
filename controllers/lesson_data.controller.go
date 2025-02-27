package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"intern_247/utils"
	"regexp"
	"strconv"
	"time"
)

type NewLessonDataInput struct {
	Name     string    `json:"name"`
	Type     int       `json:"type"`
	LessonId uuid.UUID `json:"lesson_id"`
	Position uint64    `json:"position"`
	// link video
	Url     string `json:"url"`
	Content string `json:"content"`
	//test
	PointTest    uint       `json:"point_test"` // bài kiểm tra có điểm là bao nhiêu
	TestId       uuid.UUID  `json:"test_id"`
	ExpiredAt    *time.Time `json:"expired_at"`
	AllowExpired bool       `json:"allow_expired"`
	//document
	FileSize    uint64     `json:"file_size"`
	FileName    string     `json:"file_name"`
	ContentType string     `json:"content_type"`
	DocumentId  *uuid.UUID `json:"document_id"`
}
type UpdateLessonDataInput struct {
	ID       uuid.UUID `json:"id"`
	Type     int       `json:"type"`
	Name     *string   `json:"name"`
	Position *uint64   `json:"position"`
	// link video
	Url     *string `json:"url"`
	Content *string `json:"content"`
	//test
	PointTest    *uint      `json:"point_test"`
	TestId       *uuid.UUID `json:"test_id"`
	ExpiredAt    *time.Time `json:"expired_at"`
	AllowExpired bool       `json:"allow_expired"`
	//document
	FileSize    uint64 `json:"file_size"`
	ContentType string `json:"content_type"`
	FileName    string `json:"file_name"`
	//dữ liệu bài học đã hoàn thành cấu hình điều kiện
	DoneType   int        `json:"done_type"`
	DoneValue  string     `json:"done_value"`
	DocumentId *uuid.UUID `json:"document_id"`
}

func CreateLessonData(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	var input NewLessonDataInput
	if err := c.BodyParser(&input); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	if !utils.IsValidStrLen(input.Name, 250) {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_DATA_MAX_SIZE_250)
	}
	lesson, err := repo.GetLessonByIdAndCenterId(input.LessonId, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
	}
	fmt.Println("ok1")
	if lesson.ClassId != nil {
		_, err := repo.GetClassByIdAndCenterId(*lesson.ClassId, *user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
		}
	}
	if lesson.ParentId == nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}

	if !utils.Contains(consts.LESSON_DATAS_TYPE, input.Type) {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	if input.Type == consts.TEST_TYPE {
		// bài kiểm tra
		if utils.ContainSpecialCharacter(input.Name) {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_LESSON_DATA_TEST_NAME_CONTAIN_SPECIAL_CHARACTER)
		}
		if lesson.ClassId != nil {
			_, err := repo.GetLessonDataByNameAndClassId(input.Name, uuid.Nil, *lesson.ClassId, uuid.Nil)
			if err == nil {
				return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_LESSON_DATA_TEST_NAME_EXIST_IN_CLASS)
			}
		} else if lesson.SubjectId != nil {
			_, err = repo.GetLessonDataByNameAndClassId(input.Name, uuid.Nil, uuid.Nil, *lesson.SubjectId)
			if err == nil {
				return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_LESSON_DATA_TEST_NAME_EXIST_IN_CLASS)
			}
		}
	}
	lessonData, errCode := getCustomLessonData(input)
	lessonData.DocumentId = input.DocumentId
	if errCode != 0 {
		return ResponseError(c, fiber.StatusBadRequest, "Custom failed", strconv.Itoa(errCode))
	}
	lessonData.CenterId = *user.CenterId
	lessonData.CreatedBy = user.ID
	newLessonData, err := repo.CreateLessonData(lessonData)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.CreateFail, consts.CreateFailed)
	}
	return ResponseSuccess(c, fiber.StatusCreated, consts.CreateSuccess, newLessonData)
}
func getCustomLessonData(input NewLessonDataInput) (models.LessonData, int) {
	var (
		lessonData models.LessonData
		metadata   consts.LessonDataMetadata
	)
	lessonData.Name = input.Name
	lessonData.Type = input.Type
	lessonData.LessonId = input.LessonId
	lessonData.Position = input.Position
	switch input.Type {
	case consts.YOUTUBE_LINK_TYPE:
		{
			regex := `^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube(?:-nocookie)?\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|live\/|v\/)?)([\w\-]+)(\S+)?$`
			re := regexp.MustCompile(regex)
			if !re.MatchString(input.Url) { // kiểm tra xem đường link ytb có hợp lệ không
				return lessonData, consts.ERROR_LESSON_DATA_YOUTUBE_LINK_INVALID
			}
			metadata.Url = input.Url
			jsonData, err := json.Marshal(metadata)
			lessonData.Metadata = jsonData
			if err != nil {
				return lessonData, consts.ERROR_LESSON_DATA_LESSON_DATA_MARSHAL_METADATA_FAILED
			}
			return lessonData, 0
		}
	case consts.S3_LINK_TYPE:
		{
			metadata.Url = input.Url
			jsonData, err := json.Marshal(metadata)
			lessonData.Metadata = jsonData
			if err != nil {
				return lessonData, consts.ERROR_LESSON_DATA_LESSON_DATA_MARSHAL_METADATA_FAILED
			}
			return lessonData, 0
		}
	case consts.PARAGRAPH_TYPE:
		{
			if input.Content == "" {
				return lessonData, consts.ERROR_LESSON_DATA_PARAGRAPH_CONTENT_INVALID
			}
			metadata.Content = input.Content
			jsonData, err := json.Marshal(metadata)
			lessonData.Metadata = jsonData
			if err != nil {
				return lessonData, consts.ERROR_LESSON_DATA_LESSON_DATA_MARSHAL_METADATA_FAILED
			}
			return lessonData, 0
		}
	case consts.TEST_TYPE:
		{
			if input.PointTest == consts.LESSON_DATA_POINT_FOLLOW_TEST || input.PointTest == consts.LESSON_DATA_POINT_NOT_CHECK {
				metadata.PointTest = input.PointTest
			}
			if input.ExpiredAt != nil {
				if time.Now().After(*input.ExpiredAt) {
					return lessonData, consts.ERROR_LESSON_DATA_EXPIRED_TIME_INVALID
				}
				metadata.ExpiredAt = input.ExpiredAt
				metadata.AllowExpired = input.AllowExpired
			}
			metadata.TestId = &input.TestId
			var test models.TestService
			if err := app.Database.DB.Model(&models.TestService{}).Where("id = ?", input.TestId).First(&test).Error; err != nil {
				return lessonData, consts.ERROR_LESSON_DATA_TEST_ID_INVALID
			} else {
				metadata.TestName = test.Name
				metadata.Duration = test.Duration
				metadata.MaxAnswers = test.MaxAnswers
			}
			jsonData, err := json.Marshal(metadata)
			lessonData.Metadata = jsonData
			if err != nil {
				return lessonData, consts.ERROR_LESSON_DATA_LESSON_DATA_MARSHAL_METADATA_FAILED
			}
			return lessonData, 0
		}
	case consts.DOCUMENT_TYPE:
		{
			regex := `(?i)\.(pdf|doc|xls|xlsx|pptx|png|jpg)$`
			re := regexp.MustCompile(regex)
			if !re.MatchString(input.Url) {
				return lessonData, consts.ERROR_LESSON_DATA_URL_DOCUMENT_INVALID
			}
			metadata.Url = input.Url
			metadata.FileSize = input.FileSize
			metadata.FileName = input.FileName
			metadata.ContentType = input.ContentType
			jsonData, err := json.Marshal(metadata)
			lessonData.Metadata = jsonData
			if err != nil {
				return lessonData, consts.ERROR_LESSON_DATA_LESSON_DATA_MARSHAL_METADATA_FAILED
			}
			return lessonData, 0
		}
	default:
		return lessonData, consts.ERROR_LESSON_DATA_TYPE_LESSON_DATA_NOT_FOUND
	}
}
func UpdateLessonData(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	var (
		inputs []UpdateLessonDataInput
		ids    []uuid.UUID
		// classIds []uuid.UUID
	)
	updateSlice := make(map[uuid.UUID]*UpdateLessonDataInput)
	if err := c.BodyParser(&inputs); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	for i := range inputs {
		if inputs[i].Name != nil && !utils.IsValidStrLen(*inputs[i].Name, 250) {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
		}
		ids = append(ids, inputs[i].ID)
		updateSlice[inputs[i].ID] = &inputs[i]
	}
	lessonDatas, err := repo.GetLessonDatasPreloadClassByIdsAndCenterId(ids, *user.CenterId)
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	if len(lessonDatas) != len(inputs) {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	// for i := range lessonDatas {
	// 	if lessonDatas[i].Lesson.ClassId != nil {
	// 		classIds = append(classIds, *lessonDatas[i].Lesson.ClassId)
	// 	}
	// }
	// if len(classIds) == 1 {
	// 	class, err := repo.GetClassByIdAndCenterId(classIds[0], *user.CenterId)
	// 	if err != nil {
	// 		return ResponseError(c, fiber.StatusBadRequest, "Not found", consts.DataNotFound)
	// 	}
	// 	if class.StartAt != nil {
	// 		if time.Now().After(*class.StartAt) && len(class.StudentsClasses) > 0 {
	// 			return ResponseError(c, fiber.StatusBadRequest, "Not found", consts.ERROR_CLASS_CAN_NOT_UPDATE_INPROGRESS)
	// 		}
	// 	}
	// }

	if len(inputs) == 1 {
		lessonDataUpdate := inputs[0]
		if !utils.Contains(consts.LESSON_DATAS_TYPE, inputs[0].Type) {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
		}
		lessonData, errCode := getCustomUpdateLessonData(lessonDataUpdate, lessonDatas[0])
		if errCode != 0 {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, strconv.Itoa(errCode))
		}
		lessonDataFind, err := repo.GetLessonDataByIdAndCenterId(lessonData.ID, *user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
		}
		lessonData.LessonId = lessonDataFind.LessonId
		lessonData.DocumentId = inputs[0].DocumentId

		lessonData.Type = inputs[0].Type
		if lessonData.Type == consts.TEST_TYPE {
			if utils.ContainSpecialCharacter(lessonData.Name) {
				return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_LESSON_DATA_TEST_NAME_CONTAIN_SPECIAL_CHARACTER)
			}
			lesson, err := repo.GetLessonByIdAndCenterId(lessonData.LessonId, *user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
			}
			if lesson.ClassId != nil {
				_, err := repo.GetLessonDataByNameAndClassId(lessonData.Name, lessonData.ID, uuid.Nil, *lesson.ClassId)
				if err == nil {
					return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_LESSON_DATA_TEST_NAME_EXIST_IN_CLASS)
				}
			} else if lesson.SubjectId != nil {
				_, err := repo.GetLessonDataByNameAndClassId(lessonData.Name, lessonData.ID, uuid.Nil, *lesson.SubjectId)
				if err == nil {
					return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_LESSON_DATA_TEST_NAME_EXIST_IN_CLASS)
				}
			}
		}
		row, err := repo.UpdateLessonData(&lessonData)
		if err != nil || row == 0 {
			return ResponseError(c, fiber.StatusInternalServerError, consts.InvalidInput, consts.ERROR_INTERNAL_SERVER_ERROR)
		}
		return ResponseSuccess(c, fiber.StatusOK, "success", lessonData)
	}

	for i := range lessonDatas {
		if updateSlice[lessonDatas[i].ID] == nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
		}
		lessonMetaDataUpdate := updateSlice[lessonDatas[i].ID]
		lessonDataCustom, errCode := getCustomUpdateLessonData(*lessonMetaDataUpdate, lessonDatas[i])
		if errCode != 0 {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, strconv.Itoa(errCode))
		}
		lessonDatas[i] = lessonDataCustom
	}
	_, err = repo.UpdateMultipleLessonDatas(&lessonDatas)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFail, consts.UpdateFailed)
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, lessonDatas)
}
func getCustomUpdateLessonData(update UpdateLessonDataInput, lessonData models.LessonData) (models.LessonData, int) {
	var metadata consts.LessonDataMetadata
	_ = json.Unmarshal(lessonData.Metadata, &metadata)

	if update.Position != nil {
		lessonData.Position = *update.Position
	}
	if update.Name != nil {
		lessonData.Name = *update.Name
	}
	switch update.Type {
	case consts.YOUTUBE_LINK_TYPE:
		{
			if update.Url != nil {
				regex := `((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube(?:-nocookie)?\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|live\/|v\/)?)([\w\-]+)(\S+)?$`
				re := regexp.MustCompile(regex)
				if !re.MatchString(*update.Url) {
					return lessonData, consts.ERROR_LESSON_DATA_YOUTUBE_LINK_INVALID
				}
				metadata.Url = *update.Url
			}
			metadata.DoneType = update.DoneType
			metadata.DoneValue = update.DoneValue
		}
	case consts.S3_LINK_TYPE:
		{
			if update.Url != nil {
				metadata.Url = *update.Url
			}
		}
	case consts.PARAGRAPH_TYPE:
		{
			contentNotNil := update.Content != nil
			if contentNotNil && *update.Content == "" {
				return lessonData, consts.ERROR_LESSON_DATA_PARAGRAPH_CONTENT_INVALID
			}
			if contentNotNil {
				metadata.Content = *update.Content
			}
		}
	case consts.TEST_TYPE:
		{
			if update.PointTest != nil {
				if *update.PointTest != consts.LESSON_DATA_POINT_FOLLOW_TEST && *update.PointTest != consts.LESSON_DATA_POINT_NOT_CHECK {
					return lessonData, consts.ERROR_LESSON_DATA_POINT_TEST_INVALID
				}
				metadata.PointTest = *update.PointTest
			}
			if update.ExpiredAt != nil {
				if time.Now().After(*update.ExpiredAt) {
					return lessonData, consts.ERROR_LESSON_DATA_EXPIRED_TIME_INVALID
				}
				metadata.ExpiredAt = update.ExpiredAt
				metadata.AllowExpired = update.AllowExpired
			}
			if update.TestId != nil {
				metadata.TestId = update.TestId
			}
			var test models.TestService
			if err := app.Database.DB.Model(&models.TestService{}).Where("id = ?", *update.TestId).First(&test).Error; err != nil {
				return lessonData, consts.ERROR_LESSON_DATA_TEST_ID_INVALID
			} else {
				metadata.TestName = test.Name
				metadata.Duration = test.Duration
				metadata.MaxAnswers = test.MaxAnswers
			}
		}
	case consts.DOCUMENT_TYPE:
		{
			if update.Url != nil {
				regex := `(?i)\.(pdf|doc|xls|xlsx|pptx|png|jpg)$`
				re := regexp.MustCompile(regex)
				if !re.MatchString(*update.Url) {
					return lessonData, consts.ERROR_LESSON_DATA_URL_DOCUMENT_INVALID
				}
				metadata.Url = *update.Url
				metadata.FileSize = update.FileSize
				metadata.FileName = update.FileName
				metadata.ContentType = update.ContentType
			}
		}
	default:
		return lessonData, consts.ERROR_LESSON_DATA_TYPE_LESSON_DATA_NOT_FOUND
	}
	if update.DoneType == consts.LESSON_DATA_DONE_TYPE_ANSWER_QUESTION || update.DoneType == consts.LESSON_DATA_DONE_TYPE_WATCHED_CONTENT {
		metadata.DoneType = update.DoneType
		metadata.DoneValue = update.DoneValue
	}
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return lessonData, consts.ERROR_LESSON_DATA_LESSON_DATA_MARSHAL_METADATA_FAILED
	}
	lessonData.Metadata = jsonData
	return lessonData, 0
}
func GetDetailLessonData(c *fiber.Ctx) error {
	token, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	query := new(consts.Query)
	if err := c.QueryParser(query); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidInput)
	}
	lessonData, err := repo.GetLessonDataByIdAndCenterId(query.ID, token.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
	}
	return ResponseSuccess(c, fiber.StatusOK, "Success", lessonData)
}
func DeleteLessonData(c *fiber.Ctx) error {
	token, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	var (
		input consts.DeleteInput
	)
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidInput)
	}
	_, err = repo.GetLessonDataPreloadClassByIdsAndCenterId(input.Id, token.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "lesson not found", consts.CreateFailed)
	}
	row, err := repo.DeleteLessonDataByIdAndCenterId(input.Id, token.CenterId)
	if row < 1 || err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.DeleteFail, consts.DeletedFailed)
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, row)
}
func GetListLessonDatas(c *fiber.Ctx) error {
	token, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	query := new(consts.Query)
	if err := c.QueryParser(query); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidInput)
	}
	lessonData, err := repo.GetListLessonDatas(*query, token.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
	}

	return ResponseSuccess(c, fiber.StatusOK, "Success", lessonData)
}
