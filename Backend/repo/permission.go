package repo

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/models"
)

func GetPermissionId(DB *gorm.DB) (uuid.UUID, error) {
	var (
		entry       models.Permission
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).First(&entry).Select("id")
	return entry.ID, err.Error
}

//func GetPermissionId(DB *gorm.DB, conditions map[string]interface{}) (uuid.UUID, error) {
//	var (
//		entry       models.Permission
//		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
//	)
//	defer cancel()
//	err := DB.WithContext(ctx).Where(conditions).First(&entry).Select("id")
//	return entry.ID, err.Error
//}
