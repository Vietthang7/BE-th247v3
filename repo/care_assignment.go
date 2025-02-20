package repo

import (
	"context"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/models"
)

type (
	CareAssignment  models.CareAssignment
	CareAssignments []models.CareAssignment
)

func (u *CareAssignment) Create(tx *gorm.DB) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Create(&u).Error
}
func (u *CareAssignment) First(query interface{}, args []interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Omit("created_at", "updated_at").Where(query, args...).Preload("OrganStruct", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).First(&u).Error
}
