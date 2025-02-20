package repo

import (
	"github.com/google/uuid"
	"intern_247/app"
	"intern_247/models"
)

func GetClassesBySubjectIdAndCenterId(subjectId, centerId uuid.UUID) ([]models.Class, error) {
	var classes []models.Class
	db := app.Database.DB.Where("subject_id = ? AND center_id = ?", subjectId, centerId).Find(&classes)
	return classes, db.Error
}
