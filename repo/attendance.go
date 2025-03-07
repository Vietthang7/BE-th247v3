package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/models"
)

func GetAttendanceByScheduleIdsAndClassId(scheduleIds []uuid.UUID, classId uuid.UUID) ([]models.SessionAttendance, error) {
	var attendancers []models.SessionAttendance
	db := app.Database.DB.Where("class_id = ? AND class_session_id IN (?)", classId, scheduleIds).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, full_name")
	}).Find(&attendancers)
	return attendancers, db.Error
}
