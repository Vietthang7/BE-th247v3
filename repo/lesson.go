package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/models"
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
