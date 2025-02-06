package repo

import (
	"context"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/utils"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (u *Student) Create() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()

	username := utils.GenerateUniqueUsername()
	u.Username = username
	if u.Type == consts.Official || u.Type == consts.Trial {
		tx := app.Database.DB.WithContext(ctx).Begin() // transaction
		if err = tx.Create(&u).Error; err != nil {
			tx.Rollback()
			return err
		} // lưu thông tin vào bảng student
		var (
			pwd = app.Config("DEFAULT_PASSWORD")
		)
		temp, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		loginInfo := LoginInfo{
			ID:           u.ID,
			CenterID:     u.CenterId,
			Username:     username,
			Phone:        u.Phone,
			Email:        u.Email,
			PasswordHash: string(temp),
			RoleId:       consts.Student,
			DeletedAt:    nil,
		}
		if err = loginInfo.Create(); err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit().Error
	}
	return app.Database.DB.WithContext(ctx).Model(&models.Student{}).Create(&u).Error
}
func (u *Student) Update(origin Student, query interface{}, args []interface{}) (err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		tx          = app.Database.DB.WithContext(ctx).Begin()
	)
	defer cancel()

	if err = tx.Model(&models.Student{}).Where(query, args...).Updates(&u).Error; err != nil {
		logrus.Error(err)
		tx.Rollback()
		return err
	}

	if err = app.Database.DB.WithContext(ctx).Model(&models.LoginInfo{}).Where("id = ?", u.ID).
		Update("phone", u.Phone).Update("email", u.Email).Error; err != nil {
		logrus.Error(err)
		tx.Rollback()
		return err
	}

	if !u.EmailVerified {
		if err = tx.Model(&models.Student{}).Where(query, args...).
			Update("EmailVerified", false).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
