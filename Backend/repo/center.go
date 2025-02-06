package repo

import (
	"context"
	"github.com/google/uuid"
	"intern_247/app"
	"intern_247/models"

	"gorm.io/gorm"
)

func CreateCenter(DB *gorm.DB, page *models.Center) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return DB.WithContext(ctx).Create(&page).Error
}
func GetCenterIDByUserID(user_id uuid.UUID) (models.Center, error) {
	var center models.Center
	query := app.Database.DB.Select("id,domain").Where("user_id = ?", user_id).First(&center)
	return center, query.Error
}
func ReadCenter(DB *gorm.DB) (models.Center, error) {
	var (
		entry       models.Center
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).First(&entry)
	return entry, err.Error
}
func UpdateCenter(DB *gorm.DB, center *models.Center) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return DB.WithContext(ctx).Updates(&center).Error
}
