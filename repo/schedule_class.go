package repo

import (
	"context"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/utils"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	db := app.Database.DB.Where("id IN ? AND center_id = ?", ids, centerId).Find(&schedules)
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
	db.Joins("INNER JOIN classes ON classes.`id` = schedule_classes.`class_id`").Where("classes.`status` != ? AND classes.`status` != ?", consts.CLASS_CANCELED, consts.CLASS_FINISHED)
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
func GetDetailScheduleByClassId(classId, centerId uuid.UUID) (models.ScheduleClass, error) {
	var schedule models.ScheduleClass
	query := app.Database.DB.Model(&models.ScheduleClass{})

	query.Preload("Childrens", func(db1 *gorm.DB) *gorm.DB {
		return db1.
			Omit("created_at", "updated_at").
			Order("start_date ASC, created_at ASC").
			Preload("WorkSession", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "title")
			}).
			Preload("Classrooms", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name")
			}).
			Preload("Asistant", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "full_name")
			}).
			Preload("Teacher", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "full_name")
			}).
			Preload("Childrens", func(db2 *gorm.DB) *gorm.DB {
				return db2.
					Omit("created_at", "updated_at").
					Order("created_at ASC").
					Preload("WorkSession", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "title")
					}).
					Preload("Classrooms", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "name")
					}).
					Preload("Asistant", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "full_name")
					}).
					Preload("Teacher", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "full_name")
					}).
					Preload("Attendancers", func(db *gorm.DB) *gorm.DB {
						return db.Limit(1)
					})
			}).
			Preload("Attendancers", func(db *gorm.DB) *gorm.DB {
				return db.Limit(1)
			})
	})

	query.Debug().
		Omit("created_at", "updated_at").
		Where("class_id = ? AND center_id = ? AND `type` IS NOT NULL", classId, centerId).
		First(&schedule)

	return schedule, query.Error
}

func GetScheduleClassForStudent(query consts.Query, user TokenData) ([]models.ScheduleClass, error) {
	var schedules []models.ScheduleClass

	db := app.Database.Model(&models.ScheduleClass{}).
		Select(
			"schedule_classes.`id`", "schedule_classes.`name`",
			"schedule_classes.`start_time`", "schedule_classes.`end_time`",
			"schedule_classes.`class_id`", "schedule_classes.`teacher_id`",
			"schedule_classes.`start_date`", "schedule_classes.`index`",
		).
		Where("schedule_classes.`center_id` = ? AND schedule_classes.`type` IS NULL", user.CenterId)

	// Lọc theo lớp học
	if query.Class != "" {
		db.Where("class_id = ?", query.Class)
	}

	// Lọc theo ngày bắt đầu
	if query.StartAt != "" {
		startAt, err := utils.ConvertStringToTime(query.StartAt)
		if err != nil {
			return schedules, err
		}
		db.Where("DATE(schedule_classes.start_date) >= ?", startAt.Format("2006-01-02"))
	}

	// Lọc theo ngày kết thúc
	if query.EndAt != "" {
		endAt, err := utils.ConvertStringToTime(query.EndAt)
		if err != nil {
			return schedules, err
		}
		db.Where("DATE(schedule_classes.start_date) <= ?", endAt.Format("2006-01-02"))
	}

	db.Joins("INNER JOIN classes ON classes.id = schedule_classes.class_id")

	// Lọc theo chi nhánh
	if query.Branch != "" {
		db.Where("classes.branch_id = ?", query.Branch)
	}

	// Lọc theo giáo viên
	if query.Teacher != "" {
		db.Where("schedule_classes.`teacher_id` = ?", query.Teacher)
	}

	// Tìm kiếm theo tên lớp
	if query.Search != "" {
		db.Where("classes.name LIKE ?", "%"+query.Search+"%")
	}

	// Lọc theo phòng học
	if query.Classroom != "" {
		db.Joins("INNER JOIN schedule_classrooms ON schedule_classrooms.schedule_class_id = schedule_classes.`id`").
			Where("schedule_classrooms.classroom_id = ?", query.Classroom)
	}

	// Lọc theo môn học
	if query.Subject != "" {
		db.Joins("INNER JOIN subjects ON classes.subject_id = subjects.id").
			Where("subjects.id = ?", query.Subject)
	}

	// Load danh sách học sinh tham gia lớp học
	db.Preload("Attendancers", func(db *gorm.DB) *gorm.DB {
		if user.RoleId == consts.Student {
			db = db.Where("student_id = ?", user.ID)
		}
		return db.Select("class_id", "student_id", "class_session_id")
	})

	// Lọc theo học sinh
	if query.StudentId != "" {
		var student Student
		if err := student.First("id = ?", []interface{}{query.StudentId}); err != nil {
			return nil, err
		}

		if student.Type == consts.Trial {
			db.Joins("JOIN student_sessions ss ON (ss.student_id = ? AND schedule_classes.id = ss.session_id)", user.ID)
		} else {
			db.Joins("INNER JOIN student_classes as sc ON sc.`class_id` = schedule_classes.`class_id`").
				Where("sc.`student_id` = ? AND (sc.`status` != ? OR sc.`status` IS NULL)", query.StudentId, consts.Reserved)
		}
	}

	// Kiểm tra quyền của người dùng HR
	if user.RoleId == consts.CenterHR {
		if user.Position == consts.Teacher || user.Position == consts.TeachingAssistant {
			db.Where("schedule_classes.`teacher_id` = ? OR schedule_classes.`asistant_id` = ?", user.ID, user.ID)
		}
		if user.BranchId != nil {
			db.Where("classes.`branch_id` = ?", user.BranchId)
		}
	}

	// Loại bỏ các lớp bị hủy
	db.Where("classes.`status` != ?", consts.CLASS_CANCELED)

	// Load thông tin chi tiết của lớp học
	db.Preload("Class", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "metadata")
		}).Select("id", "name", "subject_id")
	})

	// Load thông tin giáo viên
	db.Preload("Teacher", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "full_name")
	})

	// Load thông tin phòng học
	db.Preload("Classrooms", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "is_online", "room_type", "metadata")
	})

	// Debug SQL query và sắp xếp theo ngày bắt đầu
	db.Debug().Order("schedule_classes.`start_date` ASC").Find(&schedules)

	return schedules, db.Error
}
func CountScheduleClass(query string, args ...interface{}) (count int64) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	if err := app.Database.DB.WithContext(ctx).Model(&models.ScheduleClass{}).Where(query, args...).Count(&count); err != nil {
		logrus.Error(err)
	}
	return
}
func GetScheduleClassByStudentId(studentId, centerId uuid.UUID) ([]models.ScheduleClass, error) {
	var scheduleClasses []models.ScheduleClass
	db := app.Database.DB.Where("schedule_classes.`center_id` = ?", centerId)
	db.Joins("INNER JOIN classes c ON c.id = schedule_classes.class_id")
	db.Joins("INNER JOIN student_classes sc ON schedule_classes.class_id = student_classes.class_id")
	db.Order("schedule_classes.`start_date` ASC , schedule_classes.`start_time`, schedule_classes.`index`")
	db.Where("student_classes.`student_id` = ? AND schedule_classes.`type` IS NULL AND c.`status` != ? AND (student_classes.status != ? OR student_classes.status IS NULL)", studentId, consts.CLASS_CANCELED, consts.Reserved).Find(&scheduleClasses)
	return scheduleClasses, db.Error
}
func GetScheduleClassesByClassIdsAndCenterId(classIds []uuid.UUID, centerId uuid.UUID) ([]models.ScheduleClass, error) {
	var schedules []models.ScheduleClass
	db := app.Database.DB.Table("schedule_classes as sc").
		Order("sc.`start_date` ASC , sc.`start_time` ASC, sc.`end_time` ASC,sc.`created_at` ASC").
		Joins("INNER JOIN classes c on c.id = sc.class_id").
		Where("sc.class_id IN (?) AND sc.center_id = ? AND sc.`start_date` IS NOT NULL AND c.status != ?", classIds, centerId, consts.CLASS_CANCELED).
		Preload("Class", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).Find(&schedules)
	return schedules, db.Error
}
