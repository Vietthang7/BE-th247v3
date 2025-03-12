package repo

import (
	"context"
	"fmt"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAllLessonBySubjectIdAndCenterId(subjectId, centerId uuid.UUID) ([]*models.Lesson, error) {
	var lessons []*models.Lesson
	db := app.Database.DB.Debug().Model(&models.Lesson{}).Where("center_id = ?", centerId).Omit("created_at", "updated_at")
	db.Preload("Childrens", func(db *gorm.DB) *gorm.DB {
		return db.Where("class_id IS NULL").Order("position ASC, updated_at ASC").Preload("LessonDatas")
	})
	db.Where("parent_id IS NULL AND subject_id = ? AND class_id IS NULL", subjectId)
	db.Order("position ASC, updated_at ASC").Find(&lessons)
	return lessons, db.Error
}
func DeleteLessonBySubjectIdAndCenterId(subjectId, centerId uuid.UUID) (int64, error) {
	db := app.Database.DB.Where("subject_id = ? AND center_id = ? AND class_id IS NULL", subjectId, centerId).Delete(&models.Lesson{})
	return db.RowsAffected, db.Error
}
func GetLessonByIdAndCenterId(id, centerId uuid.UUID) (models.Lesson, error) {
	var lesson models.Lesson
	db := app.Database.DB.Model(&models.Lesson{}).Where("id = ? AND center_id = ?", id, centerId).Omit("created_at", "updated_at").First(&lesson)
	return lesson, db.Error
}
func CreateLesson(lesson *models.Lesson) (*models.Lesson, error) {
	db := app.Database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, //kiểm tra nếu trùng id thì update nhưng fields phía dưới còn không thì sẽ tạo cái mới
		DoUpdates: clause.AssignmentColumns([]string{"name", "position", "free_trial"}),
	}).Create(&lesson)
	return lesson, db.Error
}
func CreateLessons(lessons *[]*models.Lesson) (*[]*models.Lesson, error) {
	db := app.Database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "position", "free_trial"}),
	}).Create(&lessons)
	return lessons, db.Error
}

func GetLessonsByIdsWithParentAndCenterId(ids []uuid.UUID, parentId, centerId uuid.UUID) ([]*models.Lesson, error) {
	var lessons []*models.Lesson
	db := app.Database.DB.Where("parent_id = ? AND center_id = ?", parentId, centerId).Find(&lessons, ids)
	return lessons, db.Error
}
func GetLessonsByIdsAndCenterId(ids []uuid.UUID, centerId uuid.UUID) ([]*models.Lesson, error) {
	var lessons []*models.Lesson
	db := app.Database.DB.Where("center_id = ?", centerId).Find(&lessons, ids)
	return lessons, db.Error
}
func UpdateLessonsByCenterId(lessons []*models.Lesson, centerId uuid.UUID) ([]*models.Lesson, error) {
	db := app.Database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Where("center_id = ?", centerId).Create(&lessons)
	return lessons, db.Error
}
func GetLessonIdsByParentIdAndCenterId(parentId uuid.UUID, centerId uuid.UUID) ([]models.Lesson, error) {
	var lessons []models.Lesson
	db := app.Database.DB.Select("id", "schedule_id").Where("parent_id = ? AND center_id = ?", parentId, centerId).Find(&lessons)
	return lessons, db.Error
}
func DeleteLessonByIdsAndCenterId(ids []uuid.UUID, centerId uuid.UUID) (int64, error) {
	db := app.Database.DB.Where("id IN ? AND center_id = ?", ids, centerId).Delete(&models.Lesson{})
	return db.RowsAffected, db.Error
}
func DeleteLessonsByParentIdAndCenterId(parentIds []uuid.UUID, centerId uuid.UUID) (int64, error) {
	db := app.Database.Where("parent_id IN ? AND center_id = ?", parentIds, centerId).Delete(&models.Lesson{})
	return db.RowsAffected, db.Error
}
func DeleteLessonDataByLessonIdsAndCenterId(lessonIds []uuid.UUID, centerId uuid.UUID) (int64, error) {
	db := app.Database.DB.Where("lesson_id IN ? AND center_id = ?", lessonIds, centerId).Delete(&models.LessonData{})
	return db.RowsAffected, db.Error
}
func GetDetailLessonByIdAndCenterId(query consts.Query, centerId uuid.UUID) (models.Lesson, error) {
	var lesson models.Lesson
	db := app.Database.DB.Model(&models.Lesson{}).Where("id = ? AND center_id = ?", query.ID, centerId).Omit("created_at", "updated_at")
	db.Preload("Childrens", func(db *gorm.DB) *gorm.DB {

		return db.Preload("ScheduleClass", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "start_date", "start_time", "end_time")
		}).Select("id", "parent_id", "name", "position", "free_trial", "is_live", "metadata", "schedule_id").Order("position ASC, updated_at ASC")
	})
	db.First(&lesson)
	return lesson, db.Error
}
func GetListLessons(query consts.Query, centerId uuid.UUID) ([]models.Lesson, error) {
	var lessons []models.Lesson
	db := app.Database.DB.Model(&models.Lesson{}).Where("center_id = ?", centerId).Omit("created_at", "updated_at")
	if query.Relation != "" {
		db.Where("subject_id = ? AND class_id IS NULL", query.Relation)
	}
	if query.Children == "true" {
		db.Preload("Childrens", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "parent_id", "name", "position", "free_trial", "is_live", "metadata").Order("position ASC, updated_at ASC")
		})
	}
	db.Where("parent_id IS NULL")
	db.Order("position ASC, updated_at ASC").Find(&lessons)
	return lessons, db.Error
}

func CountLessonData(query string, args []interface{}) (count int64) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	app.Database.DB.WithContext(ctx).Model(&models.LessonData{}).Where(query, args...).Count(&count)
	return
}
func FilterDetailLessonByLive(isLive bool, classId, centerId uuid.UUID) ([]models.Lesson, error) {
	var lessons []models.Lesson
	db := app.Database.DB.Where("class_id = ? AND center_id = ? AND parent_id IS NULL", classId, centerId)
	db.Preload("Childrens", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_live = ?", isLive).Order("position ASC, created_at ASC").Select("id", "name", "parent_id")
	})
	db.Order("position ASC, created_at ASC").Find(&lessons)
	if db.Error != nil {
		fmt.Println(db.Error.Error())
	}
	return lessons, db.Error
}
func FilterLessonsByLive(isLive bool, classId, centerId uuid.UUID) ([]models.Lesson, error) {
	var lessons []models.Lesson
	db := app.Database.DB.Where("is_live = ? AND class_id ? AND center_id = ?", isLive, classId, centerId).Find(&lessons)
	return lessons, db.Error
}
