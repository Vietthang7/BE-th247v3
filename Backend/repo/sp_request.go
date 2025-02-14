package repo

import (
	"context"
	"intern_247/app"
	"intern_247/models"

	"github.com/sirupsen/logrus"
)

type (
	SupportRequest  models.SupportRequest
	SupportRequests []models.SupportRequest
)

func (u *SupportRequest) Create() (err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = app.Database.DB.WithContext(ctx)
	)
	defer cancel()

	if len(u.SubjectIds) > 0 {
		tx := DB.Begin()
		if err = tx.Create(&u).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}

		var entries []models.SubjectSpRequest
		for _, id := range u.SubjectIds {
			entries = append(entries, models.SubjectSpRequest{SubjectId: id, SupportRequestId: u.ID})
		}
		if err = tx.Model(models.SubjectSpRequest{}).Create(&entries).Error; err != nil {
			return
		}

		return tx.Commit().Error
	}

	return DB.Create(&u).Error
}
