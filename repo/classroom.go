package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
)

func TsCreateClassroom(entry *models.Classroom, slots []models.TimeSlot, shifts []models.ShortShift) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(&entry).Error; err != nil {
			return err
		}
		if len(shifts) > 0 && len(slots) > 0 {
			var scheduleId uuid.UUID
			if scheduleId, err = TsCreateRoomScheduleByClassroom(tx, entry.ID, *entry.CenterId); err != nil {
				return err
			}
			if err = CreateRoomScheduleData(tx, *entry, scheduleId, slots, shifts); err != nil {
				return err
			}
		}
		return nil
	})
}

func TsUpdateClassroom(entry *models.Classroom, origin models.Classroom, slots []models.TimeSlot, shifts []models.ShortShift) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	tx := app.Database.DB.WithContext(ctx).Begin()
	defer func() {
		if recover() != nil || err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = tx.Model(&models.Classroom{}).Where("id = ?", entry.ID).Updates(&entry).Error; err != nil {
		logrus.Error(err)
		return
	}
	if *origin.IsOnline != *entry.IsOnline {
		if *entry.IsOnline {
			// Nếu cập nhật phòng offline thành phòng online thì xóa dữ liệu lịch phòng
			if err = TsDeleteRoomSchedule(tx.Where("classroom_id = ?", entry.ID)); err != nil {
				logrus.Error(err)
				return
			}
			if err = TsDeleteTimeSlot(tx.Where("classroom_id = ?", entry.ID)); err != nil {
				logrus.Error(err)
				return
			}
			if err = TsDeleteShift(tx.Where("classroom_id = ?", entry.ID)); err != nil {
				logrus.Error(err)
				return
			}
		} else {
			if len(shifts) > 0 && len(slots) > 0 {
				var scheduleId uuid.UUID
				if scheduleId, err = TsCreateRoomScheduleByClassroom(tx, entry.ID, *entry.CenterId); err != nil {
					return
				}
				if err = CreateRoomScheduleData(tx, *entry, scheduleId, slots, shifts); err != nil {
					return
				}
			}
		}
	} else {
		if !*entry.IsOnline {
			if origin.Schedule != nil {
				if err = UpdateRoomScheduleData(tx, *origin.Schedule, slots, shifts); err != nil {
					logrus.Error(err)
					return
				}
			} else {
				if len(slots) > 0 && len(shifts) > 0 {
					var scheduleId uuid.UUID
					if scheduleId, err = TsCreateRoomScheduleByClassroom(tx, entry.ID, *entry.CenterId); err != nil {
						return
					}
					if err = CreateRoomScheduleData(tx, *entry, scheduleId, slots, shifts); err != nil {
						return
					}
				}
			}
		}
	}
	return tx.Commit().Error
}
func FindClassrooms(p *consts.RequestTable, query interface{}, args []interface{}) ([]models.Classroom, error) {
	var (
		entries     []models.Classroom
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = p.CustomOptions(app.Database.DB)
	)
	defer cancel()
	if p.Search != "" {
		DB = DB.Where("name LIKE ?", "%"+p.Search+"%")
	}
	err := DB.WithContext(ctx).Where(query, args...).Preload("Branch", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Find(&entries)
	return entries, err.Error
}
func CountClassroom(query interface{}, args []interface{}) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	app.Database.DB.Where(query, args...).Model(&models.Classroom{}).WithContext(ctx).Count(&count)
	return count
}
func ClassroomIsArranged(classroomId string) (bool, error) {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	if err := app.Database.DB.WithContext(ctx).Raw("SELECT COUNT(*) FROM `schedule_classrooms`\n"+"WHERE classroom_id = ?;", classroomId).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func DeleteClassroom(entry *models.Classroom) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	if *entry.IsOnline {
		return app.Database.DB.WithContext(ctx).Delete(entry).Error
	} else {
		// If deleting offline room then delete room schedule's data
		return app.Database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err = tx.Delete(entry).Error; err != nil {
				return err
			}
			if err = TsDeleteRoomSchedule(tx.Where("classroom_id = ?", entry.ID)); err != nil {
				logrus.Error(err)
				return err
			}
			if err = TsDeleteTimeSlot(tx.Where("classroom_id = ?", entry.ID)); err != nil {
				logrus.Error(err)
				return err
			}
			if err = TsDeleteShift(tx.Where("classroom_id = ?", entry.ID)); err != nil {
				logrus.Error(err)
				return err
			}
			return nil
		})
	}
}
func GetClassroomsByIdsAndBranchCenterId(ids []uuid.UUID, branchId, centerId uuid.UUID) ([]models.Classroom, error) {
	var classrooms []models.Classroom
	db := app.Database.DB.Where("id IN ? AND branch_id = ? AND center_id = ?", ids, branchId, centerId).Find(&classrooms)
	return classrooms, db.Error
}
