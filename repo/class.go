package repo

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/models"
)

func GetClassesBySubjectIdAndCenterId(subjectId, centerId uuid.UUID) ([]models.Class, error) {
	var classes []models.Class
	db := app.Database.DB.Where("subject_id = ? AND center_id = ?", subjectId, centerId).Find(&classes)
	return classes, db.Error
}
func FirstClassroom(query interface{}, args []interface{}, preload ...string) (models.Classroom, error) {
	var (
		entry       models.Classroom
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		err         error
		DB          = app.Database.DB.WithContext(ctx).Where(query, args...)
	)
	defer cancel()
	if len(preload) > 0 {
		PreloadClassroom(DB, preload...)
		err = DB.First(&entry).Error
		if entry.Schedule != nil && len(entry.Schedule.RoomShifts) > 0 {
			ShortenRoomShifts(entry.Schedule)
		}
	} else {
		err = DB.First(&entry).Error
	}
	return entry, err
}

func PreloadClassroom(DB *gorm.DB, properties ...string) {
	for _, v := range properties {
		if v == "Branch" {
			DB.Preload("Branch", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name")
			})
		}
		if v == "Schedule" {
			DB.Preload("Schedule", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "classroom_id", "center_id").
					Preload("RoomShifts", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "work_session_id", "schedule_id", "day_of_week")
					}).Preload("TimeSlots", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "schedule_id", "work_session_id", "start_time", "end_time")
				})
			})
		}
	}
}
