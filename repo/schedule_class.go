package repo

import (
	"context"
	"github.com/google/uuid"
	"intern_247/app"
	"intern_247/models"
)

// TeacherIsArranged Kiểm tra xem giáo viên có được sắp xếp dạy lớp nào không.
func TeacherIsArranged(teacherId uuid.UUID) bool {
	var (
		err         error
		entry       models.ScheduleClass
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	if err = app.Database.DB.WithContext(ctx).Where("teacher_id = ? OR asistant_id = ?", teacherId, teacherId).First(&entry).Error; err == nil {
		return true
	}
	return false
}
func GetScheduleClassById(id uuid.UUID, centerId uuid.UUID) (models.ScheduleClass, error) {
	var schedule models.ScheduleClass
	db := app.Database.DB.Debug().Where("id = ? AND center_id = ?", id, centerId).First(&schedule)
	return schedule, db.Error
}
func GetScheduleClassByIds(ids []uuid.UUID, centerId uuid.UUID) ([]models.ScheduleClass, error) {
	var schedules []models.ScheduleClass
	db := app.Database.DB.Debug().Where("id IN ? AND center_id = ?", ids, centerId).Find(&schedules)
	return schedules, db.Error
}
