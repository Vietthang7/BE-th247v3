package repo

import (
	"github.com/google/uuid"
	"intern_247/app"
	"intern_247/models"
	"time"
)

func GetHolidayByDateBranchIdAndCenterId(startAt *time.Time, branchId uuid.UUID, centerId uuid.UUID) ([]models.Holiday, error) {
	var holidays []models.Holiday
	db := app.Database.DB.Where("end_day >= ? AND (branch_id is NULL OR branch_id = ?) AND center_id = ?", startAt.Format("2006-01-02"), branchId, centerId).Find(&holidays)
	return holidays, db.Error
}
