package repo

import (
	"errors"
	"fmt"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
func GetSubjectByNameAndCenterId(name string, centerId uuid.UUID) (models.Subject, error) {
	var subject models.Subject
	db := app.Database.DB.Select("id", "name").Where("name = ? and center_id = ?", name, centerId).First(&subject)
	return subject, db.Error
}

func CreateSubject(subject *models.Subject) (*models.Subject, error) {
	query := app.Database.DB.Create(&subject)
	return subject, query.Error
}
func GetSubjectByNameAndIdAndCenterId(name string, id uuid.UUID, centerId uuid.UUID) (models.Subject, error) {
	var subject models.Subject
	db := app.Database.DB.Select("id", "name").Where("name = ? and center_id = ? and id != ?", name, centerId, id).First(&subject)
	return subject, db.Error
}
func GetSubjectByIdAndCenterId(id, centerId uuid.UUID) (models.Subject, error) {
	var subject models.Subject
	db := app.Database.DB.Select("id", "name").Where("id = ? AND center_id = ?", id, centerId).First(&subject)
	return subject, db.Error
}
func UpdateSubject(subject *models.Subject) (*models.Subject, error) {
	tx := app.Database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return subject, err
	}
	if err := tx.Where("id = ?", subject.ID).Updates(&subject).Error; err != nil {
		tx.Rollback()
		return subject, err
	}
	if len(subject.Teachers) > 0 {
		if err := tx.Model(&subject).Association("Teachers").Replace(subject.Teachers); err != nil {
			tx.Rollback()
			return subject, err
		}
	}
	return subject, tx.Commit().Error
}

func DeleteSubjectByIdAndCenterId(id, centerId uuid.UUID) (int64, error) {
	db := app.Database.DB.Where("id = ? AND center_id = ?", id, centerId).Delete(&models.Subject{})
	return db.RowsAffected, db.Error
}

func GetDetailSubjectByIdAndCenterId(id, centerId uuid.UUID) (models.Subject, error) {
	fmt.Println("ok")
	var subject models.Subject
	query := app.Database.DB.Where("id = ? AND center_id = ?", id, centerId)
	query.Preload("Category", func(db *gorm.DB) *gorm.DB { return db.Select("id", "name") })
	query.Preload("Teachers", func(db *gorm.DB) *gorm.DB { return db.Where("position = ?", consts.Teacher).Select("id", "full_name") })
	query.Preload("Lessons", func(db *gorm.DB) *gorm.DB {
		return db.Where("parent_id IS NULL").Omit("created_at", "updated_at").Order("position ASC, updated_at DESC").Preload("Childrens", func(db1 *gorm.DB) *gorm.DB {
			return db1.Omit("created_at", "updated_at").Order("position ASC, updated_at DESC").Preload("LessonDatas", func(db2 *gorm.DB) *gorm.DB {
				return db2.Omit("created_at", "updated_at").Order("position ASC, updated_at DESC")
			})
		})
	})
	query.Omit("created_at", "updated_at").First(&subject)
	return subject, query.Error
}

func GetListSubjectsByCenterId(q consts.Query, user TokenData) ([]models.SubjectMoreInfo, consts.Pagination, error) {
	var (
		subjects   []models.SubjectMoreInfo
		pagination consts.Pagination
		isActive   *bool
	)
	db := app.Database.DB.Model(&models.Subject{}).Where("subjects.center_id = ?", user.CenterId).Select("subjects.id", "subjects.code", "subjects.name", "subjects.thumbnail", "subjects.fee_type", "subjects.category_id", "subjects.origin_fee", "subjects.discount_fee", "subjects.created_at", "subjects.updated_at", "subjects.is_active", "subjects.total_lessons", "(SELECT COUNT(classes.id) FROM classes WHERE subject_id = subjects.id AND classes.deleted_at IS NULL) AS class_total", "(SELECT COUNT(ss.`student_id`) FROM student_subjects as ss LEFT JOIN student_classes AS sc ON ss.student_id = sc.student_id JOIN students as s ON s.id = ss.student_id WHERE sc.student_id IS NULL AND ss.subject_id = subjects.id AND s.deleted_at IS NULL) as student_pendings")
	db.Joins(`JOIN (SELECT code, name, MAX(updated_at) AS latest_updated_at FROM subjects WHERE center_id = ? AND deleted_at IS NULL GROUP BY code, name) AS latest_subjects ON (subjects.code = latest_subjects.code OR subjects.name = latest_subjects.name) AND subjects.updated_at = latest_subjects.latest_updated_at`, user.CenterId)
	isActive = q.GetActive()
	db.Preload("Teachers", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "full_name", "center_id")
	})
	if isActive != nil {
		db = db.Where("subjects.is_active = ?", *isActive)
	}
	if q.Relation != "" {
		db.Where("subjects.category_id = ?", q.Relation)
	}
	if q.Search != "" {
		searchStr := "%" + q.Search + "%"
		db.Where("(subjects.`name` LIKE ? OR subjects.`code` LIKE ?)", searchStr, searchStr)
	}

	db.Preload("Category", func(db1 *gorm.DB) *gorm.DB {
		return db1.Select("id", "name")
	})
	if q.StudentId != "" {
		db.Joins("INNER JOIN student_subjects ON student_subjects.subject_id = subjects.id").
			Where("student_subjects.student_id = ?", q.StudentId)
	}
	db.Count(&pagination.TotalResults)
	db.Order(fmt.Sprintf("%s %s, subjects.updated_at DESC", q.GetField(consts.SubjectField, "subjects.created_at"), q.GetSort())).Offset(q.GetOffset()).Limit(q.GetPageSize()).Find(&subjects)

	pagination.CurrentPage = q.GetPage()
	pagination.TotalPages = pagination.GetTotalPages(q.GetPageSize())
	return subjects, pagination, db.Error
}
func GetAllSubjectByCenterId(q consts.Query, centerId uuid.UUID) ([]models.Subject, error) {
	var (
		subjects []models.Subject
		isActive *bool
	)
	db := app.Database.DB.Select("subjects.`id`, subjects.`name`").Where("subjects.center_id", centerId)
	db.Joins(`INNER JOIN (SELECT name, MAX(updated_at) AS latest_updated_at FROM subjects WHERE center_id = ? AND deleted_at IS NULL AND name = subjects.name GROUP BY name) AS latest_subjects ON subjects.name = latest_subjects.name AND subjects.updated_at = latest_subjects.latest_updated_at`, centerId)
	//if q.Curriculum != "" {
	//	db.Joins("INNER JOIN curriculum_subjects as cs ON cs.`subject_id` = subjects.`id`")
	//	db.Where("cs.`curriculum_id` = ?", q.Curriculum)
	//}
	isActive = q.GetActive()
	if isActive != nil {
		db = db.Where("subjects.`is_active` = ?", *isActive)
	}
	db.Order("subjects.`created_at` DESC").Find(&subjects)
	return subjects, db.Error
}

func CheckSubjectsExist(subjectIDs []uuid.UUID) error {
	if len(subjectIDs) == 0 {
		return nil
	}

	var count int64
	if err := app.Database.DB.Model(&models.Subject{}).
		Where("id IN (?)", subjectIDs).
		Count(&count).Error; err != nil {
		return errors.New("lỗi khi kiểm tra môn học")
	}

	if count != int64(len(subjectIDs)) {
		return errors.New("một hoặc nhiều môn học không tồn tại")
	}

	return nil
}
