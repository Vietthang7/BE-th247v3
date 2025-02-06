package repo

import (
	"github.com/google/uuid"
	"intern_247/app"
	"intern_247/models"
)

func GetSubjectByIdsAndCenterId(ids []*uuid.UUID, centerId uuid.UUID, isActive *bool) ([]*models.Subject, error) {
	var subjects []*models.Subject
	query := app.Database.DB.Where("center_id = ? AND id IN ?", centerId, ids)
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	query.Find(&subjects)
	return subjects, query.Error
}
