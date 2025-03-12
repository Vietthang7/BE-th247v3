package repo

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/utils"
	"time"
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
func GetSingleScheduleClassByClassId(classId, centerId uuid.UUID) (models.ScheduleClass, error) {
	var schedule models.ScheduleClass
	query := app.Database.DB.Model(&models.ScheduleClass{})
	query.Omit("created_at", "updated_at").Where("class_id = ? AND center_id = ? AND `type` IS NOT NULL", classId, centerId).First(&schedule)
	return schedule, query.Error
}

func GetListScheduleByClassId(classId uuid.UUID, query consts.Query, user TokenData) ([]models.ScheduleClass, error) {
	var schedules []models.ScheduleClass
	db := app.Database.DB.Model(&models.ScheduleClass{})
	if query.StartAt != "" {
		startAt, err := utils.ConvertStringToTime(query.StartAt)
		if err != nil {
			return schedules, err
		}
		db.Where("DATE(start_date) >= ?", startAt.Format("2006-01-02"))
	}
	if query.EndAt != "" {
		endAt, err := utils.ConvertStringToTime(query.EndAt)
		if err != nil {
			return schedules, err
		}
		db.Where("DATE(start_date) <= ?", endAt.Format("2006-01-02"))
	}
	if user.RoleId != consts.Student && query.StudentId != "" {
		db.Preload("Teacher", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "full_name")
		})
		db.Preload("Attendancers", func(db1 *gorm.DB) *gorm.DB {
			db1 = db1.Where("student_id = ?", query.StudentId)
			return db1.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "full_name")
			}).Select("class_id", "class_session_id", "student_id", "user_id")
		})
		db.Joins("INNER JOIN student_classes ON schedule_classes.`class_id` = student_classes.`class_id`").Where("student_classes.`student_id` = ?", query.StudentId)
	}
	db.Select("schedule_classes.`id`", "schedule_classes.`name`", "schedule_classes.`start_date`", "schedule_classes.`index`", "schedule_classes.`start_time`", "schedule_classes.`end_time`", "teacher_id", "schedule_classes.`work_session_id`").Where("schedule_classes.`class_id` = ? AND schedule_classes.`center_id` = ? AND `type` IS NULL", classId, user.CenterId).Order("schedule_classes.`start_date` ASC, schedule_classes.`start_time` ASC, schedule_classes.`index` ASC, schedule_classes.`created_at` ASC").Find(&schedules)
	return schedules, db.Error
}

func GetScheduleClassByTeacherIdAndLimitDate(classId, teacherId uuid.UUID, times []time.Time) ([]models.ScheduleClass, error) {
	var schedules []models.ScheduleClass
	var dateStrings []string
	// Vì sql không hỗ trợ  time.Time trực tiếp, ta cần chuyển times[] thành danh sách string có định dạng YYYY-MM-DD
	for _, t := range times {
		dateStrings = append(dateStrings, t.Format("2006-01-02"))
	}
	db := app.Database.DB.Order("schedule_classes.`start_date` ASC, schedule_classes.`start_time` ASC, schedule_classes.`index` ASC").Where("schedule_classes.`teacher_id` = ? AND DATE(schedule_classes.`start_date`) IN (?) AND schedule_classes.`class_id` != ?", teacherId, dateStrings, classId)
	db.Joins("INNER JOIN classes ON classes.`id` = schedule_class.`class_id`").Where("classes.`status` != ? AND classes.`status` != ?", consts.CLASS_CANCELED, consts.CLASS_FINISHED)
	db.Find(&schedules)
	return schedules, db.Error
}
func GetScheduleClassByAsistantIdAndLimitDate(classId, asistantId uuid.UUID, times []time.Time) ([]models.ScheduleClass, error) {
	var schedules []models.ScheduleClass
	var dateStrings []string
	for _, t := range times {
		dateStrings = append(dateStrings, t.Format("2006-01-02"))
	}
	db := app.Database.DB.Order("schedule_classes.`start_date` ASC, schedule_classes.`start_time` ASC, schedule_classes.`index` ASC").Where("schedule_classes.`asistant_id` = ? AND DATE(schedule_classes.`start_date`) IN (?) AND schedule_classes.`class_id` != ?", asistantId, dateStrings, classId)
	db.Joins("INNER JOIN classes ON classes.`id` = schedule_class.`class_id`").Where("classes.`status` != ? AND classes.`status` != ?", consts.CLASS_CANCELED, consts.CLASS_FINISHED)
	db.Find(&schedules)
	return schedules, db.Error
}
func GetScheduleClassByClassroomIdAndLimitDate(classId, classroomId uuid.UUID, times []time.Time) ([]models.ScheduleClass, error) {
	var schedules []models.ScheduleClass
	var dateStrings []string
	for _, t := range times {
		dateStrings = append(dateStrings, t.Format("2006-01-02"))
	}
	db := app.Database.DB.Model(&models.ScheduleClass{}).Order("schedule_classes.`start_date` ASC, schedule_classes.`start_time` ASC, schedule_classes.`index` ASC")
	db.Joins("INNER JOIN schedule_classrooms ON schedule_classes.`id` = schedule_classrooms.`schedule_class_id`")
	db.Joins("INNER JOIN classes ON classes.`id` = schedule_classes.`class_id`").Where("classes.`status` != ? AND classes.`status` != ?", consts.CLASS_CANCELED, consts.CLASS_FINISHED)
	db.Where("DATE(schedule_classes.`start_date`) IN (?) AND schedule_classrooms.`classroom_id` = ? AND schedule_classes.`class_id` != ?", dateStrings, classroomId, classId).Find(&schedules)
	return schedules, db.Error
}
func CreateScheduleClass(scheduleClass *models.ScheduleClass, listIdChanged []uuid.UUID, isChangedType bool, centerId uuid.UUID, lessons []*models.Lesson, class models.Class) (*models.ScheduleClass, error) {
	tx := app.Database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, err
	}
	if len(listIdChanged) > 0 {
		if err := tx.Select(clause.Associations).Where("class_id = ? AND is NOT IN ? AND center_id = ?", scheduleClass.ClassId, listIdChanged, centerId).Delete(&models.ScheduleClass{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Exec("DELETE FROM schedule_classrooms WHERE schedule_class_id IN (SELECT id FROM schedule_classes WHERE id IN ? AND class_id = ? AND center_id = ?)", listIdChanged, scheduleClass.ClassId, centerId).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if isChangedType {
		if err := tx.Where("class_id = ? AND center_id = ? AND `type` IS NULL", scheduleClass.ClassId, centerId).Delete(&models.ScheduleClass{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "start_date", "start_time", "end_time", "type", "work_session_id", "teacher_id", "asistant_id", "metadata"}),
	}).Session(&gorm.Session{FullSaveAssociations: true}).Create(&scheduleClass).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Model(&models.Lesson{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"schedule_id", "updated_at"}),
	}).Select("ID", "ScheduleId", "UpdatedAt").Create(&lessons).Error; err != nil {
		tx.Rollback()
		return nil, err

	}
	if err := tx.Model(&models.Class{}).Where("id = ?", class.ID).Update("end_at", class.EndAt).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return scheduleClass, tx.Commit().Error
}
