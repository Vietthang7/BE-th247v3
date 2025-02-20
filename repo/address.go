package repo

import (
	"context"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
)

type Province models.Province
type Provinces []models.Province

func (u *Province) Find(p *consts.RequestTable, query interface{}, args []interface{}) (Provinces, error) {
	var (
		entries     Provinces
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = p.CustomOptions(app.Database.DB)
	)
	defer cancel()

	if p.Search != "" {
		DB = DB.Where("name LIKE ?", "%"+p.Search+"%")
	}

	err := DB.Model(&models.Province{}).WithContext(ctx).Where(query, args...).Find(&entries)
	return entries, err.Error
}

type Ward models.Ward
type Wards []models.Ward

func (u *Ward) Find(p *consts.RequestTable, query interface{}, args []interface{}) (Wards, error) {
	var (
		entries     Wards
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = p.CustomOptions(app.Database.DB)
	)
	defer cancel()

	if p.Search != "" {
		DB = DB.Where("name LIKE ?", "%"+p.Search+"%")
	}

	err := DB.Model(&models.Ward{}).WithContext(ctx).Where(query, args...).Find(&entries)
	return entries, err.Error
}

type District models.District
type Districts []models.District

func (u *District) Find(p *consts.RequestTable, query interface{}, args []interface{}) (Districts, error) {
	var (
		entries     Districts
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = p.CustomOptions(app.Database.DB)
	)
	defer cancel()

	if p.Search != "" {
		DB = DB.Where("name LIKE ?", "%"+p.Search+"%")
	}

	err := DB.Model(&models.District{}).WithContext(ctx).Where(query, args...).Find(&entries)
	return entries, err.Error
}
