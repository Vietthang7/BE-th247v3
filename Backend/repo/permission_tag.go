package repo

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
)

func FindSubTagsByParentTag(parentTagID uuid.UUID) ([]models.PermissionTag, error) {
	var entries []models.PermissionTag
	err := app.Database.DB.Where("parent_tag_id = ?", parentTagID).Find(&entries).Error
	if err == nil && len(entries) > 0 {
		var permissions []models.Permission
		for i, v := range entries {
			entries[i].TotalPermissions = CountPermission(app.Database.DB.Where("sub_tag_id = ?", entries[i].ID))
			if permissions, err = FindPermission(app.Database.DB.Order(consts.DescCreatedAt).Where("sub_tag_id = ?", v.ID).Select("id", "name", "action")); err != nil {
				logrus.Error(err)
				return nil, err
			} else {
				if entries[i].Permissions, err = json.Marshal(permissions); err != nil {
					logrus.Error(err)
					return nil, err
				}
			}
		}
	}
	return entries, err
}
func CreatePermissionTag(DB *gorm.DB, entry *models.PermissionTag) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return DB.WithContext(ctx).Create(&entry).Error
}

func FindPermissionTags(DB *gorm.DB) ([]models.PermissionTag, error) {
	var (
		entries     []models.PermissionTag
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).Find(&entries).Error
	if err == nil && len(entries) > 0 {
		for i := range entries {
			entries[i].TotalPermissions = CountPermission(app.Database.DB.Where("tag_id = ?", entries[i].ID))
			if subTags, err := FindSubTagsByParentTag(entries[i].ID); err == nil {
				if entries[i].SubTags, err = json.Marshal(subTags); err != nil {
					logrus.Error(err)
					return nil, err
				}
			}
		}
	}
	return entries, err
}

func FirstPermissionTag(DB *gorm.DB, id string) (models.PermissionTag, error) {
	var (
		entry       models.PermissionTag
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).Where("id = ?", id).First(&entry)
	if err.Error == nil {
		if entry.ParentTagId == nil { // Load ra các sub tag
			if subTags, err := FindPermissionTags(app.Database.DB.Order(consts.DescCreatedAt).
				Where("parent_tag_id = ?", entry.ID)); err == nil {
				for i, v := range subTags {
					if permissions, err := FindPermission(DB.Order(consts.DescCreatedAt).
						Where("sub_tag_id = ?", v.ID).Select("id", "name", "action")); err != nil {
						logrus.Error(err)
					} else {
						if subTags[i].Permissions, err = json.Marshal(permissions); err != nil {
							logrus.Error(err)
						}
					}
				}
				if entry.SubTags, err = json.Marshal(subTags); err != nil {
					logrus.Error(err)
				}
			}
		}
	}
	return entry, err.Error
}
