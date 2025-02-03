package repo

import (
	"context"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/models"
)

func FirstPermissionGrp(DB *gorm.DB) (models.PermissionGroup, error) {
	var (
		entry       models.PermissionGroup
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).First(&entry)
	return entry, err.Error
}
