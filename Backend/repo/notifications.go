package repo

import (
	"context"
	"intern_247/app"
	"intern_247/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification models.Notification

func (notify *Notification) Create() error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Model(&models.Notification{}).Create(&notify).Error
}

func CheckStudentID(studentID uuid.UUID) (bool, error) {
	var student models.Student
	// Truy vấn cơ sở dữ liệu để kiểm tra sinh viên
	err := app.Database.Where("id = ?", studentID).First(&student).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (notify *Notification) Find(DB *gorm.DB) ([]*models.Notification, error) {
	var (
		entries     []*models.Notification
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()

	err := DB.WithContext(ctx).Model(&models.Notification{}).Find(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (notify *Notification) Count(DB *gorm.DB) (count int64) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	DB.Model(&models.Notification{}).WithContext(ctx).Count(&count)
	return
}
