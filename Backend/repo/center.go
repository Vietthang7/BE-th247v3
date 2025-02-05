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
	return DB.WithContext(ctx).Create(page).Error
}
func GetCenterIDByUserID(user_id uuid.UUID) (models.Center, error) {
	var center models.Center
	query := app.Database.DB.Select("id,domain").Where("user_id = ?", user_id).First(&center)
	return center, query.Error
}
