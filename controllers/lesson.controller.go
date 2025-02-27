package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"intern_247/utils"
	"time"
)

type NewLessonInput struct {
	Id           *uuid.UUID     `json:"id"`
	Name         string         `json:"name"`
	ParentId     *uuid.UUID     `json:"parent_id"`
	SubjectId    *uuid.UUID     `json:"subject_id"`
	ClassId      *uuid.UUID     `json:"class_id"`
	FreeTrial    *bool          `json:"free_trial"` // học thử miễn phí
	Position     uint64         `json:"position"`
	Metadata     datatypes.JSON `json:"metadata"`
	IsLive       *bool          `json:"is_live"`
	ChildLessons []ChildLesson  `json:"child_lessons"`
}
type ChildLesson struct {
	Id        *uuid.UUID     `json:"id"`
	Name      string         `json:"name"`
	FreeTrial *bool          `json:"free_trial"`
	Position  uint64         `json:"position"`
	IsLive    *bool          `json:"is_live"`
	Metadata  datatypes.JSON `json:"metadata"`
}
type LessonUpdateInput struct {
	ID        uuid.UUID       `json:"id"`
	Name      *string         `json:"name"`
	FreeTrial bool            `json:"free_trial"`
	Position  *uint64         `json:"position"`
	Metadata  *datatypes.JSON `json:"metadata"`
}

func CreateLessons(c *fiber.Ctx) error {
	user, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	if user.RoleId == consts.Student {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	var (
		input            NewLessonInput
		lesson           models.Lesson
		childLessonsData []*models.Lesson
		childLessonIds   []uuid.UUID
		hasLive          bool
	)
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	// Kiểm tra xem có bài học nào là live không
	childLessonsLen := len(input.ChildLessons)
	for i := range input.ChildLessons {
		if input.ChildLessons[i].IsLive != nil && *input.ChildLessons[i].IsLive {
			hasLive = true
			break
		}
	}
	if hasLive && input.ClassId == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	if (input.IsLive != nil && *input.IsLive) || hasLive {
		if input.ClassId != nil {
			_, err := repo.GetClassByIdAndCenterId(*input.ClassId, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, "Not found", consts.DataNotFound)
			}
		}
	}
	// just children
	if input.ParentId != nil {
		parentLesson, err := repo.GetLessonByIdAndCenterId(*input.ParentId, user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, "Not found", consts.DataNotFound)
		}
		if parentLesson.ParentId != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
		}
		if parentLesson.ClassId != nil && (input.IsLive != nil && *input.IsLive) || hasLive {
			_, err := repo.GetClassByIdAndCenterId(*parentLesson.ClassId, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, "Not found", consts.DataNotFound)
			}
		}
		if childLessonsLen > 0 {
			for i := range input.ChildLessons {
				if !utils.IsValidStrLen(input.ChildLessons[i].Name, 250) {
					return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_DATA_MAX_SIZE_250)
				}
				lessonChild := &models.Lesson{Name: input.ChildLessons[i].Name, Position: input.ChildLessons[i].Position, ParentId: &parentLesson.ID, SubjectId: parentLesson.SubjectId, FreeTrial: input.ChildLessons[i].FreeTrial, CenterId: user.CenterId, CreatedBy: user.ID, ClassId: parentLesson.ClassId, Metadata: input.ChildLessons[i].Metadata, IsLive: input.ChildLessons[i].IsLive}
				childLessonsData = append(childLessonsData, lessonChild)
			}
			newLessonChilds, err := repo.CreateLessons(&childLessonsData)
			if err != nil {
				return ResponseError(c, fiber.StatusInternalServerError, consts.InvalidInput, consts.ERROR_INTERNAL_SERVER_ERROR)
			}
			return ResponseSuccess(c, fiber.StatusCreated, consts.CreateSuccess, newLessonChilds)
		}
		//create single child
		if !utils.IsValidStrLen(input.Name, 250) {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_DATA_LONGER)
		}
		if input.SubjectId != nil {
			_, err := repo.GetSubjectByIdAndCenterId(*input.SubjectId, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
			}
		}
		if input.Id != nil {
			oldLesson, err := repo.GetLessonByIdAndCenterId(*input.Id, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
			}
			lesson = oldLesson
		}
		lesson.SubjectId = input.SubjectId
		lesson.ParentId = &parentLesson.ID
		lesson.Name = input.Name
		lesson.FreeTrial = input.FreeTrial
		lesson.Position = input.Position
		lesson.CreatedBy = user.ID
		lesson.CenterId = user.CenterId
		lesson.ClassId = parentLesson.ClassId
		lesson.IsLive = input.IsLive
		lesson.Metadata = input.Metadata
		newLesson, err := repo.CreateLesson(&lesson)
		if err != nil {
			return ResponseError(c, fiber.StatusInternalServerError, "Failed", consts.CreateFailed)
		}
		return ResponseSuccess(c, fiber.StatusCreated, consts.CreateSuccess, newLesson)
	}
	//lesson parent and child
	if !utils.IsValidStrLen(input.Name, 250) {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_DATA_LONGER)
	}
	if input.SubjectId != nil {
		_, err := repo.GetSubjectByIdAndCenterId(*input.SubjectId, user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
		}
	}
	if input.ClassId != nil {
		class, err := repo.GetClassByIdAndCenterId(*input.ClassId, user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
		}
		lesson.ClassId = &class.ID
	}
	if input.Id != nil {
		oldLesson, err := repo.GetLessonByIdAndCenterId(*input.Id, user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
		}
		lesson = oldLesson
	}
	lesson.SubjectId = input.SubjectId
	lesson.Name = input.Name
	lesson.FreeTrial = input.FreeTrial
	lesson.Position = input.Position
	lesson.CreatedBy = user.ID
	lesson.CenterId = user.CenterId
	lesson.Metadata = input.Metadata
	newLesson, err := repo.CreateLesson(&lesson)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed", consts.CreateFailed)
	}
	//lesson childs in parent
	if childLessonsLen > 0 {
		for i := range input.ChildLessons {
			if !utils.IsValidStrLen(input.ChildLessons[i].Name, 250) {
				return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_DATA_LONGER)
			}
			lessonChild := &models.Lesson{Name: input.ChildLessons[i].Name, Position: input.ChildLessons[i].Position, ParentId: &newLesson.ID, SubjectId: newLesson.SubjectId, FreeTrial: input.ChildLessons[i].FreeTrial, CenterId: user.CenterId, CreatedBy: user.ID, ClassId: newLesson.ClassId, Metadata: input.ChildLessons[i].Metadata, IsLive: input.ChildLessons[i].IsLive}
			if input.ChildLessons[i].Id != nil {
				if input.Id == nil {
					// Nếu bài học con đã tồn tại thì bài học cha cũng phải tồn lại , nếu không
					return ResponseError(c, fiber.StatusInternalServerError, consts.InvalidInput, consts.ERROR_INTERNAL_SERVER_ERROR)
				}
				lessonChild.ID = *input.ChildLessons[i].Id
				childLessonIds = append(childLessonIds, *input.ChildLessons[i].Id)
			}
			childLessonsData = append(childLessonsData, lessonChild)
		}
		childLessonIdsLen := len(childLessonIds)
		if childLessonIdsLen > 0 {
			// lấy tất cả các bài học con với center, parent và id
			lessonChilds, _ := repo.GetLessonsByIdsWithParentAndCenterId(childLessonIds, newLesson.ID, user.CenterId)
			if childLessonIdsLen != len(lessonChilds) {
				return ResponseError(c, fiber.StatusInternalServerError, consts.InvalidInput, consts.InvalidReqInput)
			}
		}
		newLesonChilds, err := repo.CreateLessons(&childLessonsData)
		if err != nil {
			return ResponseError(c, fiber.StatusInternalServerError, consts.InvalidInput, consts.ERROR_INTERNAL_SERVER_ERROR)
		}
		newLesson.Childrens = *newLesonChilds
	}
	return ResponseSuccess(c, fiber.StatusCreated, consts.CreateSuccess, newLesson)
}
func UpdateLessons(c *fiber.Ctx) error {
	user, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	var (
		inputs []*LessonUpdateInput
		ids    []uuid.UUID
	)
	updateSlice := make(map[uuid.UUID]*LessonUpdateInput)
	if err := c.BodyParser(&inputs); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	for i, v := range inputs {
		if inputs[i].Name != nil && !utils.IsValidStrLen(*inputs[i].Name, 250) {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
		}
		updateSlice[v.ID] = v
		ids = append(ids, v.ID)
	}
	lessons, _ := repo.GetLessonsByIdsAndCenterId(ids, user.CenterId)
	if len(lessons) != len(inputs) {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
	}
	if len(lessons) > 0 {
		for _, lesson := range lessons {
			if lesson.ClassId != nil {
				_, err := repo.GetClassByIdAndCenterId(*lesson.ClassId, user.CenterId)
				if err != nil {
					return ResponseError(c, fiber.StatusBadRequest, "Not found", consts.DataNotFound)
				}
			}
		}
	}
	for i := range lessons {
		if updateSlice[lessons[i].ID] == nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
		}
		lessonUpdate := updateSlice[lessons[i].ID]
		if lessonUpdate.Name != nil {
			lessons[i].Name = *lessonUpdate.Name
		}
		if lessonUpdate.Position != nil {
			lessons[i].Position = *lessonUpdate.Position
		}
		if lessonUpdate.Metadata != nil {
			lessons[i].Metadata = *lessonUpdate.Metadata
		}
		lessons[i].FreeTrial = &lessonUpdate.FreeTrial
	}
	lessonsUpdate, err := repo.UpdateLessonsByCenterId(lessons, user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed", consts.UpdateFailed)
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, lessonsUpdate)
}
func DeleteLesson(c *fiber.Ctx) error {
	user, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	var (
		inputs      consts.DeleteInput
		childIds    []uuid.UUID
		scheduleIds []uuid.UUID
	)
	if err := c.BodyParser(&inputs); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	lessonChildIds, err := repo.GetLessonIdsByParentIdAndCenterId(inputs.Id, user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "get lesson child  err", consts.ERROR_INTERNAL_SERVER_ERROR)
	}
	lesson, err := repo.GetLessonByIdAndCenterId(inputs.Id, user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "get lesson err", consts.ERROR_INTERNAL_SERVER_ERROR)
	}
	if lesson.ClassId != nil {
		_, err := repo.GetClassByIdAndCenterId(*lesson.ClassId, user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, "Not found", consts.DataNotFound)
		}
		if lesson.ScheduleId != uuid.Nil {
			schedule, err := repo.GetScheduleClassById(lesson.ScheduleId, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, "Failed to get schedule", consts.DataNotFound)
			}
			if time.Now().After(*utils.MixedDateAndTime(schedule.StartDate, schedule.StartTime)) {
				return ResponseError(c, fiber.StatusBadRequest, "Not found", consts.ERROR_SCHEDULE_CLASS_LEARNED_CAN_NOT_DELETE)
			}
		}
		// chuong hoc
		if lesson.ParentId == nil && len(lessonChildIds) > 0 {
			for i := range lessonChildIds {
				if lessonChildIds[i].ScheduleId != uuid.Nil {
					scheduleIds = append(scheduleIds, lessonChildIds[i].ScheduleId)
				}
			}
			schedules, err := repo.GetScheduleClassByIds(scheduleIds, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, "Failed to get schedules", consts.DataNotFound)
			}
			for i := range schedules {
				if time.Now().After(*utils.MixedDateAndTime(schedules[i].StartDate, schedules[i].StartTime)) {
					return ResponseError(c, fiber.StatusBadRequest, "Not found schedules", consts.ERROR_SCHEDULE_CLASS_LEARNED_CAN_NOT_DELETE)
				}
			}
		}
	}
	row, err := repo.DeleteLessonByIdsAndCenterId([]uuid.UUID{inputs.Id}, user.CenterId)
	if err != nil || row < 1 {
		return ResponseError(c, fiber.StatusInternalServerError, consts.DeleteFail, consts.ERROR_INTERNAL_SERVER_ERROR)
	}
	if len(lessonChildIds) > 0 {
		row2, _ := repo.DeleteLessonsByParentIdAndCenterId([]uuid.UUID{inputs.Id}, user.CenterId)
		for i := range lessonChildIds {
			childIds = append(childIds, lessonChildIds[i].ID)
		}
		row3, _ := repo.DeleteLessonDataByLessonIdsAndCenterId(childIds, user.CenterId)
		return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, row+row2+row3)
	} else {
		row3, _ := repo.DeleteLessonDataByLessonIdsAndCenterId(childIds, user.CenterId)
		return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, row+row3)
	}
}
func GetDetailLesson(c *fiber.Ctx) error {
	token, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	query := new(consts.Query)
	if err := c.QueryParser(query); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidInput)
	}
	lesson, err := repo.GetDetailLessonByIdAndCenterId(*query, token.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
	}
	return ResponseSuccess(c, fiber.StatusOK, "Success", lesson)
}
func GetListLessons(c *fiber.Ctx) error {
	token, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_UNAUTHORIZED)
	}
	query := new(consts.Query)
	if err := c.QueryParser(query); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidInput)
	}
	lesson, err := repo.GetListLessons(*query, token.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
	}
	return ResponseSuccess(c, fiber.StatusOK, "Success", lesson)
}
