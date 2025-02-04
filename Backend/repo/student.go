package repo

import (
	"context"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/models"
)

type (
	Student  models.Student
	Students []*models.Student
)

func (u *Student) VerifyEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Model(&models.Student{}).
		Where("email = ?", email).Update("email_verified", true).Error
}

func (u *Student) First(query interface{}, args []interface{}, preload ...string) error {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		err         error
		DB          = app.Database.DB.WithContext(ctx).Where(query, args...)
	)
	defer cancel()
	if len(preload) > 0 {
		u.PreloadStudent(DB, preload...)
		err = DB.First(&u).Error
	} else {
		err = DB.First(&u).Error
	}
	return err
}

func (u *Student) PreloadStudent(DB *gorm.DB, properties ...string) {
	for _, v := range properties {
		if v == "StudyNeeds" {
			DB.Preload("StudyNeeds", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, branch_id, student_id, studying_start_date, enrollment_id").
					Preload("Enrollment", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "name")
					})
			})
		}
		if v == "Caregiver" {
			DB.Preload("Caregiver", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "full_name")
			})
		}
		if v == "Source" {
			DB.Preload("CustomerSource", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name")
			})
		}
		if v == "Province" {
			DB.Preload("Province", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name")
			})
		}
		if v == "District" {
			DB.Preload("District", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "prefix")
			})
		}
		if v == "CustomerSource" {
			DB.Preload("CustomerSource", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name")
			})
		}
		if v == "ContactChannel" {
			DB.Preload("ContactChannel", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name")
			})
		}
	}
}
