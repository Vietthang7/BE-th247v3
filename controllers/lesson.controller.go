package controllers

import (
	"fmt"
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
		input                NewLessonInput
		lesson               models.Lesson
		childLessonsData     []*models.Lesson
		childLessonIds       []uuid.UUID
		lessonLiveInputCount int
		lessonLiveCount      int
	)
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
	}
	childLessonsLen := len(input.ChildLessons)
	for i := range input.ChildLessons {
		if input.ChildLessons[i].IsLive != nil && *input.ChildLessons[i].IsLive && input.ChildLessons[i].Id == nil {
			lessonLiveInputCount++
		}
	}

	if input.IsLive != nil && *input.IsLive && input.Id == nil {
		lessonLiveInputCount++
	}

	if input.SubjectId != nil && input.ClassId != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Just choose between subject or class", consts.DataNotFound)
	}
	// just children
	if input.ParentId != nil {
		parentLesson, err := repo.GetLessonByIdAndCenterId(*input.ParentId, user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
		}
		if parentLesson.ParentId != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.InvalidReqInput)
		}
		isLessonInClass := parentLesson.ClassId != nil
		isLessonInSubject := parentLesson.SubjectId != nil && parentLesson.ClassId == nil
		if (isLessonInClass || isLessonInSubject) && lessonLiveInputCount > 0 {
			if isLessonInSubject {
				subject, err := repo.GetSubjectByIdAndCenterId(*parentLesson.SubjectId, user.CenterId)
				if err != nil {
					return ResponseError(c, fiber.StatusBadRequest, "get subject failed", consts.DataNotFound)
				}
				lessons, err := repo.GetLessonsTypeLiveBySubjectIdAndCenterId(subject.ID, user.CenterId)
				if err != nil {
					return ResponseError(c, fiber.StatusBadRequest, "Failed when get Lesssons", consts.InvalidInput)
				}

				lessonLiveCount = len(lessons)
				if (lessonLiveCount + lessonLiveInputCount) > int(subject.TotalLessons) {
					return ResponseError(c, fiber.StatusBadRequest, "Total lesson input subject is large", consts.ERROR_LESSON_IS_LIVE_MAXIMUM)
				}
			} else if isLessonInClass {
				class, err := repo.GetClassByIdAndCenterId(*parentLesson.ClassId, user.CenterId)
				if err != nil {
					return ResponseError(c, fiber.StatusBadRequest, "Get subject by class failed", consts.DataNotFound)
				}
				lessons, err := repo.GetLessonsTypeLiveByClassIdAndCenterId(class.ID, user.CenterId)
				if err != nil {
					return ResponseError(c, fiber.StatusBadRequest, "Failed when get Lesssons", consts.DataNotFound)
				}
				lessonLiveCount = len(lessons)
				if (lessonLiveCount + lessonLiveInputCount) > int(class.TotalLessons) {
					return ResponseError(c, fiber.StatusBadRequest, "Total lesson input class is large", consts.ERROR_LESSON_IS_LIVE_MAXIMUM)
				}
			}
		}
		//add list child lesson to parent
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
	if !utils.IsValidStrLen(input.Name, 250) && input.Id == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_DATA_LONGER)
	}
	if input.SubjectId != nil {
		subject, err := repo.GetSubjectByIdAndCenterId(*input.SubjectId, user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
		}
		if lessonLiveInputCount > 0 {
			lessons, err := repo.GetAllLessonBySubjectIdAndCenterId(subject.ID, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, "Failed when get Lesssons", consts.InvalidInput)
			}

			for i := range lessons {
				for j := range lessons[i].Childrens {
					if lessons[i].Childrens[j].IsLive != nil && *lessons[i].Childrens[j].IsLive {
						lessonLiveCount++
					}
				}

			}
			if (lessonLiveCount + lessonLiveInputCount) > int(subject.TotalLessons) {
				return ResponseError(c, fiber.StatusBadRequest, "Total lesson input subject is large", consts.ERROR_LESSON_IS_LIVE_MAXIMUM)
			}
		}
	}
	if input.ClassId != nil {
		class, err := repo.GetClassByIdAndCenterId(*input.ClassId, user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
		}
		lesson.ClassId = &class.ID
		if lessonLiveInputCount > 0 {
			lessons, err := repo.GetLessonsByClassIdAndCenterId(class.ID, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, "Failed when get Lesssons", consts.DataNotFound)
			}
			for i := range lessons {
				if lessons[i].IsLive != nil && *lessons[i].IsLive {
					lessonLiveCount++
				}
			}
			if (lessonLiveCount + lessonLiveInputCount) > int(class.TotalLessons) {
				return ResponseError(c, fiber.StatusBadRequest, "Total lesson input class is large", consts.ERROR_LESSON_IS_LIVE_MAXIMUM)
			}
		}
	}
	if input.Id != nil {
		oldLesson, err := repo.GetLessonByIdAndCenterId(*input.Id, user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
		}
		if oldLesson.SubjectId != nil && lessonLiveInputCount > 0 && oldLesson.ClassId == nil { //for subject
			subject, err := repo.GetSubjectByIdAndCenterId(*oldLesson.SubjectId, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
			}
			lessons, err := repo.GetAllLessonBySubjectIdAndCenterId(subject.ID, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, "Failed when get Lesssons", consts.InvalidInput)
			}
			for i := range lessons {
				for j := range lessons[i].Childrens {
					if lessons[i].Childrens[j].IsLive != nil && *lessons[i].Childrens[j].IsLive {
						lessonLiveCount++
					}
				}
			}
			if (lessonLiveCount + lessonLiveInputCount) > int(subject.TotalLessons) {
				return ResponseError(c, fiber.StatusBadRequest, "Total lesson input subject is large old", consts.ERROR_LESSON_IS_LIVE_MAXIMUM)
			}
		}
		if oldLesson.ClassId != nil && lessonLiveInputCount > 0 { //for class
			class, err := repo.GetClassByIdAndCenterId(*oldLesson.ClassId, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
			}
			lesson.ClassId = &class.ID

			lessons, err := repo.GetLessonsByClassIdAndCenterId(class.ID, user.CenterId)
			if err != nil {
				return ResponseError(c, fiber.StatusBadRequest, "Failed when get Lesssons", consts.DataNotFound)
			}
			for i := range lessons {
				if lessons[i].IsLive != nil && *lessons[i].IsLive {
					lessonLiveCount++
				}
			}
			if (lessonLiveCount + lessonLiveInputCount) > int(class.TotalLessons) {
				return ResponseError(c, fiber.StatusBadRequest, fmt.Sprintf("Bạn đã tạo thừa %d buổi học ở "+
					"Nội dung học tập. Vui lòng tạo đủ theo cài đặt môn học.",
					lessonLiveCount+lessonLiveInputCount-int(class.TotalLessons)), nil)
			}
		}
		lesson = oldLesson
	}
	if input.Name != "" {
		lesson.Name = input.Name
	}
	lesson.SubjectId = input.SubjectId
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
					return ResponseError(c, fiber.StatusInternalServerError, consts.InvalidInput, consts.ERROR_INTERNAL_SERVER_ERROR)
				}
				lessonChild.ID = *input.ChildLessons[i].Id
				childLessonIds = append(childLessonIds, *input.ChildLessons[i].Id)
			}
			childLessonsData = append(childLessonsData, lessonChild)
		}
		childLessonIdsLen := len(childLessonIds)
		if childLessonIdsLen > 0 {
			// get all child lessons with center and parent and ids
			lessonChilds, _ := repo.GetLessonsByIdsWithParentAndCenterId(childLessonIds, newLesson.ID, user.CenterId)
			if childLessonIdsLen != len(lessonChilds) {
				return ResponseError(c, fiber.StatusInternalServerError, consts.InvalidInput, consts.InvalidReqInput)
			}
		}
		newLessonChilds, err := repo.CreateLessons(&childLessonsData)
		if err != nil {
			return ResponseError(c, fiber.StatusInternalServerError, consts.InvalidInput, consts.ERROR_INTERNAL_SERVER_ERROR)
		}
		newLesson.Childrens = *newLessonChilds
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
		ids = append(ids, inputs[i].ID)
	}
	lessons, _ := repo.GetLessonsByIdsAndCenterId(ids, user.CenterId)
	if len(lessons) != len(inputs) {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
	}
	if len(lessons) > 0 && lessons[0].ClassId != nil {
		_, err := repo.GetClassByIdAndCenterId(*lessons[0].ClassId, user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, "Not found", consts.DataNotFound)
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
