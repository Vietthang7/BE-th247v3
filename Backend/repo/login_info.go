package repo

import (
	"context"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/models"
)

type LoginInfo models.LoginInfo

func (u *LoginInfo) First(query interface{}, args []interface{}, preload ...string) error {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = app.Database.DB.WithContext(ctx).Where(query, args...)
	)
	defer cancel() // đảm bảo context sẽ được hủy sau khi hoàn thành phương thức giúp giải phóng tài nguyên
	if len(preload) > 0 {
		u.Preload(DB, preload...)
	}
	return DB.First(&u).Error
}

func (u *LoginInfo) Preload(DB *gorm.DB, properties ...string) {
	for _, v := range properties {
		if v == "User" {
			DB.Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "full_name,avatar,is_active,email_verified,position,permission_grp_id,branch_id")
			})
		}
		if v == "Student" {
			DB.Preload("Student", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "full_name,avatar,email_verified,branch_id")
			})
		}
	}
}
