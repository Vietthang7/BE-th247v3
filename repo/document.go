package repo

import (
	"context"
	"fmt"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	Document  models.Document
	Documents []models.Document
)

func (d *Document) Create() error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()

	if d.CategoryID != uuid.Nil {
		var category models.Category
		if err := app.Database.DB.WithContext(ctx).First(&category, "id = ?", d.CategoryID).Error; err != nil {
			logrus.Error("Invalid category_id: ", d.CategoryID)
			return fmt.Errorf("invalid category_id: %s", d.CategoryID)
		}
	}

	if err := app.Database.DB.WithContext(ctx).Create(&d).Error; err != nil {
		logrus.Error("Error creating Document: ", err.Error())
		return err
	}

	return nil
}

func (u *Document) First(query interface{}, args []interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Omit("created_at, updated_at").Where(query, args...).First(&u).Error
}

func (u *Document) Find(p *consts.RequestTable, query interface{}, args []interface{}, preload ...string) (entries Documents, err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = p.CustomOptions(app.Database.DB).WithContext(ctx).Where(query, args...).Omit("updated_at")
	)
	defer cancel()

	if len(preload) > 0 {
		u.Preload(DB, preload...)
	}

	err = DB.Find(&entries).Error
	return entries, err
}

func (u *Document) Preload(DB *gorm.DB, properties ...string) {
	for _, v := range properties {
		if v == "Creator" {
			DB.Preload("Creator", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, full_name")
			})
		}

		if v == "Category" {
			DB.Preload("Category", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name")
			})
		}
	}
}

func (u *Document) Count(query interface{}, args []interface{}) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	app.Database.DB.Where(query, args...).Model(&models.Document{}).WithContext(ctx).Count(&count)
	return count
}

func (u *Document) Delete() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Delete(&u).Error
}

func (u *Document) Update() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Updates(&u).Error
}
