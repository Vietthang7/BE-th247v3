package repo

import (
	"context"
	"fmt"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	SupportRequest  models.SupportRequest
	SupportRequests []models.SupportRequest
)

func (u *SupportRequest) Create() (err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = app.Database.DB.WithContext(ctx)
	)
	defer cancel()

	if len(u.SubjectIds) > 0 {
		tx := DB.Begin()
		if err = tx.Create(&u).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}

		var entries []models.SubjectSpRequest
		for _, id := range u.SubjectIds {
			entries = append(entries, models.SubjectSpRequest{SubjectId: id, SupportRequestId: u.ID})
		}
		if err = tx.Model(models.SubjectSpRequest{}).Create(&entries).Error; err != nil {
			return
		}

		return tx.Commit().Error
	}

	return DB.Create(&u).Error
}

func (u *SupportRequest) Preload(DB *gorm.DB, properties ...string) {
	for _, v := range properties {
		if v == "Creator" {
			DB.Preload("Creator", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, full_name, avatar, phone")
			})
		}
		if v == "Responder" {
			DB.Preload("Responder", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, full_name")
			})
		}
		if v == "Subjects" {
			DB.Preload("Subjects", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name")
			})
		}
	}
}

func (u *SupportRequest) Find(p *consts.RequestTable, query interface{}, args []interface{}, searchStudent string, preload ...string) (entries SupportRequests, err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = p.CustomOptions(app.Database.DB).WithContext(ctx)
	)
	defer cancel()

	if searchStudent != "" {
		searchStudent = "%" + searchStudent + "%"
		DB = DB.Joins("JOIN students s ON s.id = support_requests.created_by").
			Where("s.full_name LIKE ? OR s.phone LIKE ?", searchStudent, searchStudent)
	}
	DB = DB.Where(query, args...).
		Omit("leave_from_date, leave_until_date, file, subject_id, content, updated_at, make_up_class")

	if len(preload) > 0 {
		u.Preload(DB, preload...)
	}

	err = DB.Find(&entries).Error
	return
}

func (u *SupportRequest) Count(query interface{}, args []interface{}) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	app.Database.DB.Where(query, args...).Model(&models.SupportRequest{}).WithContext(ctx).Count(&count)
	return count
}

func (u *SupportRequest) First(query interface{}, args []interface{}, preload ...string) (err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = app.Database.DB.WithContext(ctx).Where(query, args...)
	)
	defer cancel()

	if len(preload) > 0 {
		u.Preload(DB, preload...)
	}

	if err = DB.First(&u).Error; err != nil {
		return
	}

	if u.Type == consts.SRStopStudying {
		var student Student
		fmt.Sscanf(student.PreloadCompletedSubject(u.CreatedBy), "%d/%d", &u.FinishedSubject, &u.RegisteredSubject)
	}

	return
}

func (u *SupportRequest) Update() (err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = app.Database.DB.WithContext(ctx)
	)
	defer cancel()

	if len(u.SubjectIds) > 0 {
		tx := DB.Begin()
		if err = tx.Updates(&u).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}

		var count int64
		if tx.Model(&models.SubjectSpRequest{}).Where("support_request_id = ? AND subject_id IN ?", u.ID,
			u.SubjectIds).Count(&count); count != int64(len(u.SubjectIds)) {
			if err = tx.Where("support_request_id = ?", u.ID).
				Unscoped().Delete(&models.SubjectSpRequest{}).Error; err != nil {
				logrus.Error(err)
				tx.Rollback()
				return err
			}
			var entries []models.SubjectSpRequest
			for _, id := range u.SubjectIds {
				entries = append(entries, models.SubjectSpRequest{SubjectId: id, SupportRequestId: u.ID})
			}
			if err = tx.Model(models.SubjectSpRequest{}).Create(&entries).Error; err != nil {
				return
			}
		}

		return tx.Commit().Error
	}

	return DB.Updates(&u).Error
}

func (u *SupportRequest) Delete() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Delete(&u).Error
}

func (u *SupportRequest) Respond(ids uuid.UUIDs) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Where("id IN ?", ids).Updates(&u).Error
}
