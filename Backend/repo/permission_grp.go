package repo

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
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
func CreatePermissionGrp(DB *gorm.DB, entry *models.PermissionGroup) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return DB.WithContext(ctx).Create(&entry).Error
}
func FindPermissionGrp(DB *gorm.DB) ([]models.PermissionGroup, error) {
	var (
		entries     []models.PermissionGroup
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).Find(&entries).Error
	return entries, err
}
func CountPermissionGrp(DB *gorm.DB) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	DB.Model(&models.PermissionGroup{}).WithContext(ctx).Count(&count)
	return count
}

func PreloadPermissionGrp(entry *models.PermissionGroup, properties ...string) {
	for _, val := range properties {
		if val == "tag" {
			var (
				tags        []models.CustomPermissionTag
				shortenTags []models.CustomPermissionTag
			)
			if err := json.Unmarshal(entry.Tags, &tags); err != nil {
				logrus.Error(err)
			}
			for _, v := range tags {
				if v.CountSelected > 0 {
					shortenTags = append(shortenTags, v)
				}
			}
			if tagBytes, err := json.Marshal(shortenTags); err != nil {
				logrus.Error(err)
			} else {
				entry.Tags = tagBytes
			}
		}
	}
}
func UpdatePermissionGroup(DB *gorm.DB, entry *models.PermissionGroup) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return DB.WithContext(ctx).Updates(&entry).Error
}
func DeletePermissionGroup(DB *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return DB.WithContext(ctx).Delete(&models.PermissionGroup{}).Error
}
