package repo

import (
	"context"
	"fmt"
	"intern_247/app"
	"intern_247/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateWorkSession(workSession *models.WorkSession) (int64, error) {
	query := app.Database.DB.Debug().Create(&workSession)
	return query.RowsAffected, query.Error
}

func GetListWorkSessionsByCenterId(centerId uuid.UUID) ([]models.WorkSession, error) {
	var workSessions []models.WorkSession
	tx := app.Database.DB.Model(&models.WorkSession{})
	tx.Preload("Branch", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})
	tx.Where("center_id = ?", centerId).Find(&workSessions)
	return workSessions, tx.Error
}

func ListWorkSessions(DB *gorm.DB) ([]models.WorkSession, error) {
	var (
		entries     []models.WorkSession
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).Find(&entries)
	return entries, err.Error
}

func CountWorkSession(DB *gorm.DB) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	DB.Model(&models.WorkSession{}).WithContext(ctx).Count(&count)
	return count
}

func GetWorkSessionByIdAndCenterId(id, centerId uuid.UUID) (models.WorkSession, error) {
	var workSession models.WorkSession
	tx := app.Database.DB.Model(&models.WorkSession{})
	tx.Preload("Branch", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})
	tx.Where("id = ? AND center_id = ?", id, centerId).First(&workSession)
	return workSession, tx.Error
}

func GetActiveWorkSessionByIdsAndBranchCenter(ids []uuid.UUID, branchId uuid.UUID, centerId uuid.UUID) ([]models.WorkSession, error) {
	var workSessions []models.WorkSession
	tx := app.Database.DB.Model(&models.WorkSession{})
	tx.Preload("Branch", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})
	tx.Debug().Where("id IN ? AND (branch_id = ? OR branch_id IS NULL) AND center_id = ? AND is_active = ?", ids, branchId, centerId, true).Find(&workSessions)
	return workSessions, tx.Error
}

func UpdateWorkSessionById(workSession *models.WorkSession) (int64, error) {
	query := app.Database.DB.Where("id = ?", workSession.ID).Updates(&workSession)
	if workSession.BranchId == nil {
		query.Update("branch_id", nil)
	}
	return query.RowsAffected, query.Error
}

func FindSessionByName(title string, centerId uuid.UUID) (models.WorkSession, error) {
	var workSession models.WorkSession
	query := app.Database.DB.Where("title = ? AND center_id = ?", title, centerId).First(&workSession)
	return workSession, query.Error
}

func DeleteWorkSessionsByIdsAndCenterId(ids []uuid.UUID, centerId uuid.UUID) (int64, error) {
	query := app.Database.DB.Where("center_id = ?", centerId).Delete(&models.WorkSession{}, ids)
	return query.RowsAffected, query.Error
}
func FirstWorkSession(DB *gorm.DB) (models.WorkSession, error) {
	var (
		entry       models.WorkSession
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).First(&entry)
	return entry, err.Error
}
func DeleteWorkSessions(DB *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return DB.WithContext(ctx).Delete(&models.WorkSession{}).Error
}

// true -> have dependencies
func IsWorkSessionHasDataDependencies(workSessionId, centerId uuid.UUID) bool {
	type DataDepend struct {
		ShiftTotal         int64 `json:"shift_total"`
		TimeSlotTotal      int64 `json:"time_slot_total"`
		ScheduleClassTotal int64 `json:"schedule_class_total"`
	}
	var dataTotal DataDepend
	db := app.Database.DB.Debug().Raw("SELECT (SELECT COUNT(*) FROM shifts WHERE work_session_id = ?  AND center_id = ?) as shift_total, (SELECT COUNT(*) FROM time_slots WHERE work_session_id = ? AND center_id = ?) as time_slot_total, (SELECT COUNT(*) FROM schedule_classes WHERE work_session_id = ? AND center_id = ?) as schedule_class_total", workSessionId, centerId, workSessionId, centerId, workSessionId, centerId).Scan(&dataTotal)
	if db.Error != nil {
		return true
	}
	return dataTotal.ShiftTotal > 0 || dataTotal.TimeSlotTotal > 0 || dataTotal.ScheduleClassTotal > 0
}

func CheckWorkSessionExists(workSessionID uuid.UUID) error {
	var count int64
	err := app.Database.DB.
		Model(&models.WorkSession{}).
		Where("id = ?", workSessionID).
		Count(&count).Error

	if err != nil {
		logrus.Error("Failed to check work session:", err)
		return err
	}

	if count == 0 {
		return fmt.Errorf("%s", "Work session không tồn tại")
	}

	return nil
}

func CheckWorkSessionIsActive(workSessionID uuid.UUID) error {
	var active bool
	err := app.Database.DB.
		Model(&models.WorkSession{}).
		Select("is_active").
		Where("id = ?", workSessionID).
		Scan(&active).Error

	if err != nil {
		logrus.Error("Failed to check WorkSession status:", err)
		return fmt.Errorf("lỗi kiểm tra trạng thái WorkSession")
	}

	if !active {
		return fmt.Errorf("WorkSession với ID %s không hoạt động", workSessionID)
	}

	return nil
}
