package repo

import (
	"context"
	"intern_247/app"
	"intern_247/models"

	"gorm.io/gorm"
)

//	func CreateCenter(DB *gorm.DB, page *models.Center) error {
//		ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
//		defer cancel()
//		return DB.WithContext(ctx).Create(page).Error
//	}
func CreateCenter(DB *gorm.DB, page *models.Center) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return DB.WithContext(ctx).Create(page).Error
}
