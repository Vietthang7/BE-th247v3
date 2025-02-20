package repo

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
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

//func GetListSubjectsByCenterId(q consts.Query, user TokenData) ([]models.SubjectMoreInfo, consts.Pagination, error) {
//	var (
//		subjects   []models.SubjectMoreInfo
//		pagination consts.Pagination
//		isActive   *bool
//	)
//	db := app.Database.DB.Model(&models.Subject{}).Where("subjects.center_id = ?", user.CenterId).Select("subjects.id", "subjects.code", "subjects.name", "subjects.thumbnail", "subjects.fee_type", "subjects.category_id", "subjects.origin_fee", "subject.discount_fee", "subjects.created_at", "subjects.updated_at", "subjects.is_active", "subjects.total_lessons", "(SELECT COUNT(classes.id) FROM classes WHERE subject_id = subjects.id AND classes.deleted_at IS NULL) AS class_total"
//}
