package repo

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
)

func FirstOrganStruct(DB *gorm.DB) (models.OrganStruct, error) {
	var (
		entry       models.OrganStruct
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).First(&entry)
	return entry, err.Error
}
func PreloadOrganStruct(entry *models.OrganStruct, properties ...string) {
	for _, v := range properties {
		if v == "parentName" && entry.ParentId != nil { // Load parent name
			var (
				err    error
				parent models.OrganStruct
			)
			if parent, err = FirstOrganStruct(app.Database.DB.Where(consts.NilDeletedAt).Where("id = ?", entry.ParentId).Select("name")); err != nil {
				logrus.Error(err)
			} else {
				entry.ParentName = parent.Name
			}
		}
		if v == "permissionGrpName" {
			var (
				err           error
				permissionGrp models.PermissionGroup
			)
			if permissionGrp, err = FirstPermissionGrp(app.Database.DB.Where(consts.NilDeletedAt).
				Where("id = ?", entry.PermissionGrpId)); err != nil {
				logrus.Error(err)
			} else {
				entry.PermissionGrpName = permissionGrp.Name
			}
		}
		if v == "totalUser" {
			entry.TotalUser = CountUser(app.Database.DB.Where(consts.NilDeletedAt).
				Where("organ_struct_id = ?", entry.ID))
		}
	}
}
