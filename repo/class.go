package repo

import (
	"context"
	"errors"
	"fmt"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/helpers"
	"intern_247/models"
	"intern_247/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

	// Preload dữ liệu đầy đủ
	if len(preload) > 0 {
		PreloadClassroom(DB, preload...)
	}

	// Đảm bảo load cả Address của Branch
	err = DB.Preload("Branch", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, address")
	}).Preload("Schedule").First(&entry).Error

	// Rút gọn RoomShifts nếu có
	if err == nil && entry.Schedule != nil && len(entry.Schedule.RoomShifts) > 0 {
		ShortenRoomShifts(entry.Schedule)
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

func GetClassByIdAndCenterId(id, centerId uuid.UUID) (models.Class, error) {
	var class models.Class
	db := app.Database.DB.Where("id = ? AND center_id = ?", id, centerId)
	db.Preload("StudentsClasses", func(db *gorm.DB) *gorm.DB {
		return db.Limit(1)
	})
	db.Preload("Subject", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "code")
	})
	db.First(&class)
	return class, db.Error
}
func GetClassByCodeAndCenterId(code string, classId, centerId uuid.UUID) (models.Class, error) {
	var class models.Class
	db := app.Database.DB.Where("code = ? AND center_id = ?", code, centerId)
	if classId != uuid.Nil {
		db.Where("id != ?", classId)
	}
	db.First(&class)
	return class, db.Error
}
func CreateClass(newClass *models.Class) (*models.Class, error) {
	db := app.Database.DB.Create(&newClass)
	return newClass, db.Error
}
func GetDetailClassByIdAndCenterId(id, centerId uuid.UUID, token TokenData, isChildren bool) (models.Class, error) {
	var class models.Class
	db := app.Database.DB.Where("id = ? AND center_id = ?", id, centerId)
	db.Preload("Branch", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})
	db.Preload("Subject", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name").Preload("Teachers", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "full_name")
		})
	})
	db.Preload("Creater", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "full_name")
	})

	db.Preload("Curator", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "full_name")
	})
	db.Preload("ScheduleClass", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Childrens", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Teacher", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "full_name")
			}).Preload("Asistant", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "full_name")
			}).Select("*").Preload("Childrens")
		}).Select("*").Where("parent_id IS NULL").Order("start_date DESC")
	})
	if isChildren {
		db.Preload("Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Select("*").Where("parent_id IS NULL").
				Preload("Childrens", func(db *gorm.DB) *gorm.DB {
					return db.Preload("ScheduleClass", func(db *gorm.DB) *gorm.DB {
						return db.Preload("Teacher", func(db *gorm.DB) *gorm.DB {
							return db.Select("id", "full_name")
						}).Preload("Classrooms", func(db *gorm.DB) *gorm.DB {
							return db.Select("*").Preload("Branch", func(db *gorm.DB) *gorm.DB {
								return db.Select("id", "name", "address")
							})
						})
					}).Preload("LessonDatas", func(db *gorm.DB) *gorm.DB {
						return db.Preload("Progress", func(db *gorm.DB) *gorm.DB {
							return db.Where("student_id = ?", token.ID)
						})
					})
				})
		})
	}
	db.First(&class)
	//if token.RoleId == consts.Student {
	//	class.CurrentLessonData, _ = GetCurrentLessonData(token.ID, class.ID)
	//}
	return class, db.Error
}
func UpdateClassByIdAndCenterId(class *models.Class, isChangeType, isChangeSubject bool) (models.Class, error) {
	tx := app.Database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Where("id = ? AND center_id = ?", class.ID, class.CenterId).Updates(&class).Error; err != nil {
		tx.Rollback()
		return *class, err
	}
	if isChangeType {
		if err := tx.Where("class_id = ?", class.ID).Delete(&models.StudentClasses{}).Error; err != nil {
			tx.Rollback()
			return *class, err
		}
	}
	if isChangeSubject {
		var lessons []models.Lesson
		var lessonUids []uuid.UUID
		if err := tx.Where("class_id = ?", class.ID).Find(&lessons).Error; err != nil {
			tx.Rollback()
			return *class, err
		}
		for i := range lessons {
			lessonUids = append(lessonUids, lessons[i].ID)
		}
		if err := tx.Where("class_id = ?", class.ID).Delete(&models.Lesson{}).Error; err != nil {
			tx.Rollback()
			return *class, err
		}
		if err := tx.Where("lesson_id IN ?", lessonUids).Delete(&models.LessonData{}).Error; err != nil {
			tx.Rollback()
			return *class, err
		}
	}
	return *class, tx.Commit().Error
}

func CountUnclassifiedsubjects(studentID uuid.UUID) (int64, uuid.UUIDs, error) {
	var count int64
	var unclassifiedSubjectIDs uuid.UUIDs

	var student models.Student
	err := app.Database.DB.Preload("Subjects").First(&student, "id = ?", studentID).Error
	if err != nil {
		return 0, nil, err
	}
	var subjectIds []uuid.UUID
	for _, subject := range student.Subjects {
		subjectIds = append(subjectIds, subject.ID)
	}
	// Remove duplicate subject IDs
	uniqueSubjectIds := utils.UniqueSliceElements(subjectIds)
	// Tìm môn học mà học sinh đã được phân vào lớp
	var assignedSubject []uuid.UUID
	err = app.Database.DB.Model(&models.StudentClasses{}).Select("classes.subject_id").Joins("JOIN classes ON student_classes.class_id = classes.id").
		Where("student_classes.student_id = ?", studentID).Pluck("classes.subject_id", &assignedSubject).Error
	assignedSubjectsMap := make(map[uuid.UUID]bool)
	for _, id := range assignedSubject {
		assignedSubjectsMap[id] = true
	}
	for _, id := range uniqueSubjectIds {
		if !assignedSubjectsMap[id] {
			unclassifiedSubjectIDs = append(unclassifiedSubjectIDs, id)
		}
	}
	count = int64(len(unclassifiedSubjectIDs))

	if err != nil {
		return 0, nil, err
	}
	return count, unclassifiedSubjectIDs, nil
}
func GetAllClasses(centerId uuid.UUID) ([]models.Class, error) {
	var classes []models.Class
	db := app.Database.DB.Where("center_id = ?", centerId).Find(&classes)
	return classes, db.Error
}
func SaveAllStatusClasses(classes []models.Class) error {

	db := app.Database.DB.Model(&models.Class{}).Select("Status", "ID").Save(&classes)
	if db.Error != nil {
		fmt.Println(db.Error.Error())
	}
	return db.Error
}
func GetListClassesByQueryAndCenterId(q consts.Query, centerId uuid.UUID, token TokenData) ([]*models.Class, consts.Pagination, models.ClassOverview, error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		classes     []*models.Class
		pagination  consts.Pagination
		overview    models.ClassOverview
	)
	defer cancel()
	db := app.Database.DB.WithContext(ctx).Model(&models.Class{}).Where("classes.center_id = ?", centerId)
	db.Preload("Branch", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})
	db.Preload("Subject", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})
	db.Preload("ScheduleClass", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "class_id").Preload("Classrooms", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "is_online")
		})
	})
	if q.Search != "" {
		db.Where("classes.name LIKE ? OR classes.code LIKE ?", "%"+q.Search+"%", "%"+q.Search+"%")
	}
	if q.Subject != "" {
		db.Where("classes.subject_id = ?", q.Subject)
	}
	if q.Branch != "" {
		db.Where("classes.branch_id = ?", q.Branch)
	}
	if q.Type > 0 {
		db.Where("classes.type = ?", q.Type)
	}
	if q.StartAt != "" {
		startAt, err := utils.ConvertStringToTime(q.StartAt)
		if err != nil {
			return classes, pagination, overview, err
		}
		db.Where("DATE(classes.`start_at`) >= ?", startAt.Format("2006-01-02"))
	}
	if q.EndAt != "" {
		endAt, err := utils.ConvertStringToTime(q.EndAt)
		if err != nil {
			return classes, pagination, overview, err
		}
		db.Where("DATE(classes.`start_at`) <= ?", endAt.Format("2006-01-02"))
	}
	if q.Classroom != "" {
		db.Joins("INNER JOIN schedule_classes ON schedule_classes.`class_id` = classes.`id` AND schedule_classes.deleted_at IS NULL INNER JOIN schedule_classrooms ON schedule_classrooms.`schedule_class_id` = schedule_classes.`id`")
		db.Where("schedule_classrooms.`classroom_id` = ?", q.Classroom)
	}
	classCount := db.Session(&gorm.Session{}).Select("COUNT(DISTINCT(CASE WHEN classes.status = ? THEN classes.id END))  as coming_soon, COUNT(DISTINCT(CASE WHEN classes.status = ? THEN classes.id END)) as in_progress, COUNT(DISTINCT(CASE WHEN classes.status = ? THEN classes.id END)) as finished, COUNT(DISTINCT(CASE WHEN classes.status = ? THEN classes.id END)) as canceled", consts.CLASS_COMING_SOON, consts.CLASS_IN_PROGRESS, consts.CLASS_FINISHED, consts.CLASS_CANCELED)
	if token.BranchId != nil {
		classCount.Where("classes.branch_id = ?", *token.BranchId)
	}
	if helpers.IsStudent(token.RoleId) {
		classCount.Joins("JOIN (SELECT DISTINCT sc.`class_id` FROM student_classes sc WHERE  sc.student_id = ? AND (sc.`status` IS NULL OR sc.`status` != ?)) sc ON sc.`class_id` = classes.`id`", token.ID, consts.Reserved)
	}
	if helpers.IsTeacherOrAsistant(token.RoleId, token.Position) {
		db.Joins("JOIN schedule_classes ON classes.id = schedule_classes.class_id")
		db.Where("schedule_classes.teacher_id = ? OR schedule_classes.asistant_id = ?", token.ID, token.ID)
		//counting
		classCount.Joins("JOIN (SELECT DISTINCT sc.class_id FROM  schedule_classes sc WHERE (sc.teacher_id = ? OR sc.asistant_id = ?) AND sc.deleted_at IS NULL) sc ON classes.id = sc.class_id", token.ID, token.ID)
	}
	if q.Status > 0 {
		db.Where("classes.`status` = ?", q.Status)
	}
	var classIdsByStudent uuid.UUIDs
	// Xếp lớp cho học viên
	if q.StudentId != "" {
		// Lọc danh sách lớp theo học viên.
		var student Student
		var subjectIdsByStudent uuid.UUIDs
		err := student.First("id = ?", []interface{}{q.StudentId}, "StudyNeeds")
		if err != nil {
			return classes, pagination, overview, db.Error
		}
		err = app.Database.DB.Model(&models.StudentClasses{}).Where("student_id = ?", q.StudentId).Pluck("class_id", &classIdsByStudent).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return classes, pagination, overview, db.Error
		}
		if _, subjectIdsByStudent, err = CountUnclassifiedsubjects(student.ID); err != nil || len(subjectIdsByStudent) <= 0 {
			return classes, pagination, overview, db.Error
		}
		if student.StudyNeeds != nil {
			isOffline := student.StudyNeeds.IsOfflineForm != nil && *student.StudyNeeds.IsOfflineForm
			isOnline := student.StudyNeeds.IsOnlineForm != nil && *student.StudyNeeds.IsOnlineForm
			if student.StudyNeeds.StudyingStartDate != nil {
				db.Where("? <= classes.start_at", *student.StudyNeeds.StudyingStartDate)
			}
			if isOffline && !isOnline {
				db.Where("classes.type = ?", consts.CLASS_TYPE_OFFLINE)
			} else if !isOffline && isOnline {
				db.Where("classes.type = ?", consts.CLASS_TYPE_ONLINE)
			}
		}
		db.Joins("JOIN subjects s1 ON s1.id = classes.subject_id")
		db.Where(`
                EXISTS (
                SELECT 1 
                FROM student_subjects ss 
                JOIN subjects s3 ON s3.id = ss.subject_id
                WHERE ss.student_id = ?
                AND s1.code = s3.code
                )`, student.ID)
		db.Where(`
                NOT EXISTS (
                SELECT 1 
                FROM student_classes sc
                JOIN classes c ON c.id = sc.class_id
                JOIN subjects s2 ON s2.id = c.subject_id
                WHERE sc.student_id = ?
                AND s1.code = s2.code
                )
                `, student.ID)
		db.Preload("ScheduleClass", func(db *gorm.DB) *gorm.DB {
			if q.ScheduleLength > 0 {
				db = db.Limit(q.ScheduleLength)
			}
			return db.Preload("Childrens", func(db *gorm.DB) *gorm.DB {
				return db.Select("*").Preload("Childrens")
			}).Select("*").Where("parent_id IS NULL").Order("start_date DESC")
		})
		if student.BranchId != nil && *student.BranchId != uuid.Nil {
			db.Where("classes.branch_id = ?", *student.BranchId)
		}
		db.Not("classes.status IN (?)", []int{consts.CLASS_CANCELED, consts.CLASS_FINISHED})
	}
	// Lọc danh sách lớp cho role học viên
	if helpers.IsStudent(token.RoleId) {
		db.Joins("JOIN student_classes sc ON sc.class_id = classes.id").Where("sc.student_id = ? AND sc.status IS NULL OR sc.status != ?", token.ID, consts.Reserved).
			Select("DISTINCT(classes.id)", "classes.name", "classes.type", "classes.metadata", "classes.start_at", "classes.total_lessons", "classes.code", "classes.`status`", "classes.cancel_reason"). // Thêm cancel_reason vào select
			Preload("StudentClasses", func(db *gorm.DB) *gorm.DB {
				return db.Where("student_id = ?", token.ID).Select("class_id", "student_id", "progress")
			}).
			Preload("ScheduleClass", func(db *gorm.DB) *gorm.DB {
				if q.ScheduleLength > 0 {
					db = db.Limit(q.ScheduleLength)
				}
				return db.Preload("Childrens", func(db *gorm.DB) *gorm.DB {
					return db.Preload("Teacher", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "full_name")
					}).Preload("Asistant", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "full_name")
					}).Select("*").Preload("Childrens")
				}).Select("*").Where("parent_id IS NULL").Order("start_date DESC")
			}).
			Where("classes.status != ?", consts.CLASS_CANCELED)
	} else if q.StudentId != "" {
		db.Select("DISTINCT(classes.id)", "classes.name, classes.metadata", "classes.type", "classes.start_at", "classes.curriculum_id", "classes.subject_id", "classes.branch_id", "classes.classroom_id", "classes.total_lessons", "classes.`status`", "classes.cancel_reason").Order("classes.start_at ASC") // Thêm cancel_reason vào select
	} else {
		db.Select("DISTINCT(classes.id)", "classes.created_at", "classes.name, classes.metadata", "classes.type", "classes.start_at", "classes.end_at", "classes.curriculum_id", "classes.subject_id", "classes.branch_id", "classes.classroom_id", "classes.total_lessons", "classes.`status`", "classes.`code`", "classes.cancel_reason").Order("classes.created_at DESC") // Thêm cancel_reason vào select
	}

	if !helpers.IsStudent(token.RoleId) {
		db.Preload("Exams", func(db *gorm.DB) *gorm.DB {
			return db.Select("class_id", "student_id", "result")
		})
		db.Preload("StudentsClasses", func(db *gorm.DB) *gorm.DB {
			return db.Select("class_id", "student_id")
		})
		db.Preload("SessionAttendancers", func(db *gorm.DB) *gorm.DB {
			return db.Select("class_id", "student_id")
		})
		db.Preload("Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("`position` ASC, created_at ASC").Preload("LessonDatas", func(db1 *gorm.DB) *gorm.DB {
				db1 = db1.Where("`type` = ?", consts.TEST_TYPE)
				return db1.Order("`position` ASC, created_at ASC").Select("id", "lesson_id", "metadata", "type")
			}).Order("`position` ASC").Select("id", "class_id")
		})
	}
	if token.BranchId != nil {
		db.Where("classes.branch_id = ?", *token.BranchId)
	}
	classCount.Scan(&overview)
	db.Offset(q.GetOffset()).Limit(q.GetPageSize()).Find(&classes)
	_ = db.Group("classes.id").Count(&pagination.TotalResults)
	pagination.CurrentPage = q.GetPage()
	pagination.TotalPages = pagination.GetTotalPages(q.GetPageSize())
	db.Offset(q.GetOffset()).Limit(q.GetPageSize()).Find(&classes)
	var classIds uuid.UUIDs
	for _, class := range classes {
		classIds = append(classIds, class.ID)
		class.TotalStudent = CountStudentInClass(class.ID)
		if token.RoleId == consts.Student {
			class.LessonLearned = CountLessonLearned(class.ID, token.ID)
		} else if q.StudentId != "" {
			_, found := lo.Find(classIdsByStudent, func(cc uuid.UUID) bool {
				return cc == class.ID
			})
			if found {
				class.IsAdded = &found
			} else {
				class.IsAdded = &found
			}
		}
	}

	if q.StudentId != "" && len(classIds) > 0 {
		var shiftCount int64
		app.Database.DB.WithContext(ctx).Raw("SELECT COUNT(*) FROM shifts\nWHERE deleted_at IS NULL AND "+
			"student_id = ?", q.StudentId).Scan(&shiftCount)
		if shiftCount > 0 {
			var matchClassIds uuid.UUIDs
			if err := app.Database.DB.WithContext(ctx).Raw("WITH t1 AS (\n\tSELECT day_of_week, work_session_id "+
				"FROM `shifts` WHERE `student_id` = ?\n), t2 AS (\n\tSELECT ts.work_session_id work_session_id, "+
				"t1.day_of_week FROM `time_slots` ts\n\tJOIN t1 ON t1.work_session_id = ts.work_session_id\n\t"+
				"WHERE ts.student_id = ?\n), t3 AS (\n\tSELECT sc.`class_id`, COUNT(*) total_match FROM "+
				"`schedule_classes` sc\n\tJOIN t2 ON t2.work_session_id = sc.work_session_id AND "+
				"DAYOFWEEK(sc.start_date) = t2.day_of_week + 1\n\tWHERE sc.`class_id` IN ?\n\tAND sc.`deleted_at` IS "+
				"NULL AND sc.parent_id IS NOT NULL\n\tGROUP BY sc.`class_id`\n)\n\nSELECT sc.class_id FROM "+
				"`schedule_classes` sc\nLEFT JOIN t3 ON t3.class_id = sc.class_id\nWHERE sc.`class_id` IN ?\nAND "+
				"sc.`deleted_at` IS NULL AND sc.parent_id IS NOT NULL\nGROUP BY sc.class_id, t3.total_match\nHAVING "+
				"COUNT(*) - COALESCE(t3.total_match, 0) = 0", q.StudentId, q.StudentId, classIds, classIds).
				Scan(&matchClassIds).Error; err != nil {
				return nil, consts.Pagination{}, models.ClassOverview{}, err
			}
			var filteredClasses []*models.Class
			for _, v := range classes {
				if utils.Contains(matchClassIds, v.ID) {
					filteredClasses = append(filteredClasses, v)
				}
			}
			classes = filteredClasses
		}
	}

	return classes, pagination, overview, db.Error
}

func ListStudentInClass(classId uuid.UUID, p *consts.RequestTable, query interface{}, args []interface{}) ([]*models.Student, error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()

	var class models.Class
	if err := app.Database.DB.WithContext(ctx).Preload("Students", func(db *gorm.DB) *gorm.DB {
		return p.CustomOptions(db).Where(query, args...).Preload("StudyNeeds")
	}).First(&class, "id = ?", classId).Error; err != nil {
		return nil, err
	}

	studentIds := make([]uuid.UUID, len(class.Students))
	for i, student := range class.Students {
		studentIds[i] = student.ID
	}

	var studentClasses []models.StudentClasses
	if err := app.Database.DB.Where("student_id IN (?) AND class_id = ? ", studentIds, classId).Find(&studentClasses).Error; err != nil {
		return nil, err
	}

	classMap := make(map[uuid.UUID]*time.Time)
	statusMap := make(map[uuid.UUID]int64)
	for _, sc := range studentClasses {
		classMap[sc.StudentId] = sc.CreatedAt
		if sc.Status != nil {
			statusMap[sc.StudentId] = *sc.Status
		}
	}

	for _, student := range class.Students {
		switch {
		case statusMap[student.ID] == consts.Reserved:
			student.Status = consts.Reserved
		case class.StartAt != nil && time.Now().Before(*class.StartAt):
			student.Status = consts.GoingStudy
		case class.EndAt != nil && time.Now().After(*class.EndAt):
			student.Status = consts.StudyDone
		default:
			student.Status = consts.Studying
		}
		if addedAt, exists := classMap[student.ID]; exists {
			student.AddedAt = addedAt
		}
	}

	var trialStudents []*models.Student
	if err := app.Database.DB.WithContext(ctx).Model(&models.Student{}).
		Joins("JOIN student_sessions ss ON ss.student_id = students.id").
		Joins("JOIN students s ON s.id = ss.student_id").
		Where("ss.class_id = ? AND s.type = ?", classId, consts.Trial).Distinct().
		Find(&trialStudents).Error; err != nil {
		logrus.Error(err)
		return class.Students, err
	}
	if len(trialStudents) > 0 {
		for i := range trialStudents {
			trialStudents[i].Status = consts.TrialStatus
		}
	}

	class.Students = append(class.Students, trialStudents...)
	return class.Students, nil
}

func CountStudentInClass(classId uuid.UUID) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	app.Database.DB.Model(&models.StudentClasses{}).WithContext(ctx).Where("class_id = ?", classId).Count(&count)
	return count
}

func RemoveStudentFromClass(input models.RemoveStudentsFromClassInput, token TokenData, c *fiber.Ctx) error {
	tx := app.Database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var class models.Class
	if err := tx.Model(&models.Class{}).Where("id = ? AND center_id = ?", input.ClassId, token.CenterId).First(&class).Error; err != nil {
		tx.Rollback()
		return err
	}

	if class.Status == consts.CLASS_CANCELED {
		tx.Rollback()
		return errors.New("class is canceled")
	}

	for _, studentId := range input.StudentId {
		var studentClass models.StudentClasses
		if err := tx.Model(&models.StudentClasses{}).Where("class_id = ? AND student_id = ?", input.ClassId, studentId).First(&studentClass).Error; err != nil {
			tx.Rollback()
			return errors.New("student not found in class")
		}

		if err := tx.Where("class_id = ? AND student_id = ?", input.ClassId, studentId).Delete(&studentClass).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}

		if _, err := LoadStatusTransaction(tx, studentId, token.CenterId); err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func CountLessonLearned(classId uuid.UUID, studentId uuid.UUID) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	app.Database.DB.Model(&models.SessionAttendance{}).WithContext(ctx).Where("class_id = ? AND student_id = ?", classId, studentId).Count(&count)
	return count
}

type StudentCanBeAddedIntoClass struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
}

//func ListStudentByEnrollmentPlan(classId uuid.UUID, centerId uuid.UUID, p *consts.RequestTable, query interface{}, args []interface{}) ([]*models.Student, error) {
//	var (
//		students    []*models.Student
//		class       *models.Class
//		err         error
//		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
//	)
//	defer cancel()
//	err = app.Database.DB.WithContext(ctx).Where("id = ? AND center_id = ?", classId, centerId).
//		Preload("Subject").Preload("Branch").First(&class).Error
//	if err != nil {
//		return nil, err
//	}
//	// TODO: Xử lý lọc danh sách kế hoạch tuyển sinh của lớp học
//	db := app.Database.DB.WithContext(ctx).Model(&models.Student{}).Where("students.type = ?", consts.Official).
//		Preload("Subjects")
//	db = db.Joins("JOIN study_needs ON study_needs.student_id = students.id").Where(query, args...).
//		Where("study_needs.studying_start_date <= ? OR study_needs.studying_start_date IS NULL", class.StartAt).
//		Where("(students.id) NOT IN (SELECT student_id FROM student_classes WHERE class_id = ?)", classId) //Loại bỏ học viên đã đăng ký lớp này (student_classes).
//	if class.BranchId != uuid.Nil {
//		db = db.Where("study_needs.branch_id = ?", class.BranchId)
//	}
//	switch class.Type {
//	case 1: // Online
//		db = db.Where("study_needs.is_online_form = ?", true)
//	case 2: //Offline
//		db = db.Where("study_needs.is_offline_form = ?", true)
//	case 3: //Hybrid
//		db = db.Where("study_needs.is_online_form = ? OR study_needs.is_offline_form = ?", true, true)
//	}
//	err = db.Find(&students).Error
//	if err != nil {
//		return nil, err
//	}
//	// lọc học viên dựa trên môn học
//	var (
//		mark                    = make(map[uuid.UUID]*models.Student)
//		studentIds, hasSchedule []uuid.UUID
//	)
//	for _, student := range students {
//		_, uniqueSubjectIds, _ := CountUnclassifiedsubjects(student.ID)
//		if utils.Contains(uniqueSubjectIds, class.SubjectId) {
//			mark[student.ID] = student
//			studentIds = append(studentIds, student.ID)
//		}
//	}
//
//}

func GetClassAndSubjectByIdAndCenterId(id, centerId uuid.UUID) (models.Class, error) {
	var class models.Class
	db := app.Database.DB.Where("id = ? AND center_id = ?", id, centerId)
	db.Preload("Subject", func(db1 *gorm.DB) *gorm.DB {
		return db1.Select("id", "total_lessons")
	})
	db.First(&class)
	return class, db.Error
}

func FindStudentsCanBeAddedIntoClass(class models.Class, search string) (entries []StudentCanBeAddedIntoClass, err error) {
	var (
		ctx, cancel                         = context.WithTimeout(context.Background(), app.CTimeOut)
		searchQuery, branchQuery, formQuery string
		args                                = []interface{}{class.StartAt, class.EndAt, class.Subject.Code, consts.Official, true, class.ID}
	)
	defer cancel()
	if search != "" {
		search = fmt.Sprintf("%%%s%%", search)
		searchQuery += "\n\tAND (full_name LIKE ? OR email LIKE ? OR phone LIKE ?)"
		args = append(args, search, search, search)
	}
	if class.BranchId != nil {
		branchQuery = "AND (sn.branch_id IS NULL OR sn.branch_id = '" + class.BranchId.String() + "')"
	}
	switch class.Type {
	case consts.CLASS_TYPE_OFFLINE:
		formQuery = "AND sn.is_offline_form = true"
	case consts.CLASS_TYPE_ONLINE:
		formQuery = "AND sn.is_online_form = true"
	case consts.CLASS_TYPE_HYBRID:
		formQuery = "AND (sn.is_offline_form = true OR sn.is_online_form = true)"
	}
	// Lấy ra các học viên có môn học tương ứng với lớp học , học viên không được nằm trong lớp học
	if err = app.Database.DB.WithContext(ctx).Raw(`--
		SELECT id , fullname FROM students s
		JOIN (
			SELECT DISTINCT t1.student_id FROM (
				SELECT sc.student_id,cs.subject_id FROM student_curriculums sc
				JOIN curriculums cs ON cs.curriculum_id = sc.curriculum_id
				JOIN study_needs sn ON sn.id = sc.study_need_id
				`+branchQuery+`
				AND (sn.studying_start_date IS NULL OR sn.studying_start_date <= ?)
				`+formQuery+`
				UNION
				SELECT ss.student_id , ss.subject_id FROM student_subjects ss
				JOIN study_needs sn ON sn.id = ss.study_need_id
				`+branchQuery+`
				AND (sn.studying_start_date IS NULL OR sn.studying_start_date <= ?)
				`+formQuery+`
			) t1
			JOIN subjects s ON s.id = t1.subject_id AND s.code = ?
			) t2 ON t2.student_id = s.id
			WHERE deleted_at IS NULL AND type = ? AND is_active = ? AND id NOT IN (SELECT student_id FROM student_classes WHERE class_id = ?)`+searchQuery+`ORDER BY created_at`, args...).Scan(&entries).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	if CountScheduleClass("class_id = ? AND deleted_at IS NULL AND parent_id IS NOT NULL") == 0 {
		return
	}
	var (
		studentIds, hasSchedule []uuid.UUID
		mark                    = make(map[uuid.UUID]StudentCanBeAddedIntoClass)
	)
	for _, student := range entries {
		mark[student.ID] = student
		studentIds = append(studentIds, student.ID)
	}
	// Lọc ra ID các học viên có đăng kí lịch trống
	if err = app.Database.DB.WithContext(ctx).Raw("SELECT DISTINCT student_id FROM shifts WHERE deleted_at IS NULL AND student_id IN ?", studentIds).Scan(&hasSchedule).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	var theOthers, matchSchedules uuid.UUIDs
	// Lọc ra ID các học viên không đăng ký lịch trống
	if len(hasSchedule) < 1 {
		theOthers = studentIds
	} else {
		theOthers = utils.FindTheOtherElems(studentIds, hasSchedule)
	}
	/* Với các student có đăng kí lịch, tiến hành check lịch đăng ký và lịch học của lớp.
	matchSchedules là mảng ID các học viên đăng ký lịch thỏa mãn lịch học của lớp
	*/

	if err = app.Database.DB.WithContext(ctx).Raw(`SELECT t1.student_id
	FROM schedule_classes sc
	LEFT JOIN (
		SELECT student_id, day_of_week, work_session_id
		FROM shifts
		WHERE deleted_at IS NULL
		AND student_id IN ?	
	) t1
	ON t1.work_session_id = sc.work_session_id
	AND t1.day_of_week + 1 = DAYOFWEEK(sc.start_date)
	WHERE sc.class_id = ?
	AND sc.parent_id IS NOT NULL
	GROUP BY t1.student_id
	HAVING COUNT(*) = (
		SELECT COUNT(*)
		FROM schedule_classes sc
		WHERE sc.class_id = ?
		AND sc.parent_id IS NOT NULL
		AND sc.parent_id IS NOT NULL
	)`, hasSchedule, class.ID, class.ID).Scan(&matchSchedules).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	if len(matchSchedules) > 0 {
		theOthers = append(theOthers, matchSchedules...)
	}
	var (
		filteredStudent []StudentCanBeAddedIntoClass
		ok              bool
	)
	for _, v := range theOthers {
		if _, ok = mark[v]; ok {
			filteredStudent = append(filteredStudent, mark[v])
		}
	}
	return filteredStudent, nil
}

// DeleteClass xóa lớp học theo ID
func DeleteClass(classId uuid.UUID) error {
	db := app.Database.DB

	// Tìm lớp học trước
	var class models.Class
	if err := db.First(&class, "id = ?", classId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return err // Không tìm thấy lớp học
		}
		return err // Lỗi khác từ cơ sở dữ liệu
	}

	// Xóa lớp học
	if err := db.Delete(&class).Error; err != nil {
		return err
	}

	return nil
}

// Cập nhật trạng thái của lớp học
func UpdateClassStatus(class *models.Class) error {
	// Cập nhật trạng thái lớp học trong cơ sở dữ liệu
	return app.Database.DB.Save(class).Error
}
