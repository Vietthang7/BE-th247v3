package repo

import (
	"context"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/models"
)

type (
	TuitionFeePkgCf  models.TuitionFeePkgCf
	TuitionFeePkgCfs []models.TuitionFeePkgCf
)

func (u *TuitionFeePkgCf) First(query interface{}, args []interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Omit("created_at, updated_at").Where(query, args...).First(&u).Error
}
func (u *TuitionFeePkgCf) Create(tx *gorm.DB) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return tx.WithContext(ctx).Create(&u).Error
}
