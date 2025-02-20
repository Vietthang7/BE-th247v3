package repo

import (
	"context"
	"intern_247/app"
	"intern_247/models"

	"gorm.io/gorm"
)

type StudentSubject models.StudentSubject

func (u *StudentSubject) Create(tx *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return tx.WithContext(ctx).Model(&models.StudentSubject{}).Create(&u).Error
}

func (u *StudentSubject) Count(query string, args ...interface{}) (count int64) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	app.Database.DB.WithContext(ctx).Model(models.StudentSubject{}).Where(query, args...).Count(&count)
	return
}
