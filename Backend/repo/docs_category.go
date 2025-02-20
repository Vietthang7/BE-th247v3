package repo

import (
	"context"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"

	"gorm.io/gorm"
)

type (
	DocsCategory   models.DocsCategory
	DocsCategories []models.DocsCategory
)

func (u *DocsCategory) Create() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Create(&u).Error
}

func (u *DocsCategory) Count(query interface{}, args []interface{}) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	app.Database.DB.Where(query, args...).Model(&models.DocsCategory{}).WithContext(ctx).Count(&count)
	return count
}

func (u *DocsCategory) First(query interface{}, args []interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Where(query, args...).First(&u).Error
}

func (u *DocsCategory) Find(p *consts.RequestTable, query interface{}, args []interface{}, preload ...string) (entries DocsCategories, err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = p.CustomOptions(app.Database.DB).WithContext(ctx).Where(query, args...)
	)
	defer cancel()

	if p.Search != "" {
		DB = DB.Where("name LIKE ?", "%"+p.Search+"%")
	}
	if len(preload) > 0 {
		u.Preload(DB, preload...)
	}

	err = DB.Find(&entries).Error
	return entries, err
}

func (u *DocsCategory) Preload(DB *gorm.DB, properties ...string) {
	for _, v := range properties {
		if v == "Creator" {
			DB.Preload("Creator", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, full_name")
			})
		}
	}
}

func (u *DocsCategory) Update() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Updates(&u).Update("is_active", u.IsActive).Error
}

func (u *DocsCategory) Delete() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Delete(&u).Error
}
