package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
)

func GetLessonDataByNameAndClassId(name string, lessonDataId, classId, subjectId uuid.UUID) (models.LessonData, error) {
	var lessonData models.LessonData
	db := app.Database.DB.Debug().Where("lesson_data.name = ? AND lesson_data.`type` = ?", name, consts.TEST_TYPE).Omit("lesson_data.created_at", "lesson_data.updated_at", "lesson_data.position")
	if classId != uuid.Nil {
		db.Joins("INNER JOIN lessons ON lesson_data.lesson_id = lessons.id").Where("lesson_data.`class_id` = ?", classId)
	}
	if subjectId != uuid.Nil {
		db.Joins("INNER JOIN lessons ON lesson_data.lesson_id = lessons.id").Where("lesson_data.`subject` = ?", subjectId)
	}
	//update
	if lessonDataId != uuid.Nil {
		db.Where("lesson_data.`id` != ?", lessonDataId)
	}
	db.First(&lessonData)
	return lessonData, db.Error
}
func CreateLessonData(lessonData models.LessonData) (*models.LessonData, error) {
	db := app.Database.DB.Debug().Create(&lessonData)
	return &lessonData, db.Error
}

func GetLessonDatasPreloadClassByIdsAndCenterId(ids []uuid.UUID, centerId uuid.UUID) ([]models.LessonData, error) {
	var lessonDatas []models.LessonData
	db := app.Database.DB.Debug().Where("center_id = ?", centerId)
	db.Preload("Lesson", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Class")
	})
	db.Find(&lessonDatas, ids)
	return lessonDatas, db.Error
}
func GetLessonDataByIdAndCenterId(id, centerId uuid.UUID) (models.LessonData, error) {
	var lessonData models.LessonData
	db := app.Database.DB.Where("id = ? AND center_id = ?", id, centerId).Omit("created_at", "updated_at", "position").First(&lessonData)
	return lessonData, db.Error
}
func UpdateLessonData(lessonData *models.LessonData) (int64, error) {
	db := app.Database.DB.Updates(&lessonData)
	return db.RowsAffected, db.Error
}
func UpdateMultipleLessonDatas(lessonDatas *[]models.LessonData) (int64, error) {
	db := app.Database.DB.Select("ID", "Position").Save(&lessonDatas)
	return db.RowsAffected, db.Error
}
func GetLessonDataPreloadClassByIdsAndCenterId(id, centerId uuid.UUID) (models.LessonData, error) {
	var lessonData models.LessonData
	db := app.Database.DB.Where("id = ? AND center_id = ?", id, centerId)
	db.Preload("Lesson", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Class")
	})
	db.First(&lessonData)
	return lessonData, db.Error
}
func DeleteLessonDataByIdAndCenterId(id uuid.UUID, centerId uuid.UUID) (int64, error) {
	db := app.Database.DB.Where("id = ? AND center_id = ?", id, centerId).Delete(&models.LessonData{})
	return db.RowsAffected, db.Error
}
func GetListLessonDatas(query consts.Query, centerId uuid.UUID) ([]models.LessonData, error) {
	var lessonDatas []models.LessonData
	db := app.Database.DB.Where("center_id = ?", centerId).Omit("created_at", "updated_at")

	if query.Relation != "" {
		db.Where("lesson_id = ?", query.Relation)
	}
	db.Order("position ASC, updated_at ASC").Find(&lessonDatas)
	return lessonDatas, db.Error
}
