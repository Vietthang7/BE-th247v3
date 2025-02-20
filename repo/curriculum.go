package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/models"
)

func GetCurriculumsBySubjectIdAndCenterId(subjectId, centerId uuid.UUID) ([]models.Curriculum, error) {
	var curriculums []models.Curriculum
	db := app.Database.DB.Model(&models.Curriculum{}).Joins("JOIN curriculum_subjects cs ON cs.curriculum_id = curriculums.id AND cs.subject_id = ?", subjectId).Where("curriculums.center_id = ?", centerId)
	db.Preload("Subjects", func(db *gorm.DB) *gorm.DB {
		return db.Select("id")
	})
	db.Find(&curriculums)
	return curriculums, db.Error
}

func CreateCurriculums(curriculums []models.Curriculum) ([]models.Curriculum, error) {
	query := app.Database.DB.Create(&curriculums)
	return curriculums, query.Error
}
