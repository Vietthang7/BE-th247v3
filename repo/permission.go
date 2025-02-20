package repo

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/models"
)

func FindPermission(DB *gorm.DB) ([]models.Permission, error) {
	var (
		entries     []models.Permission
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).Find(&entries).Error
	return entries, err
}
func GetPermissionId(DB *gorm.DB) (uuid.UUID, error) {
	var (
		entry       models.Permission
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).First(&entry).Select("id")
	return entry.ID, err.Error
}
func CreatePermission(DB *gorm.DB, entry *models.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return DB.WithContext(ctx).Create(&entry).Error
}

func CountPermission(DB *gorm.DB) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	DB.Model(&models.Permission{}).WithContext(ctx).Count(&count)
	return count
}
