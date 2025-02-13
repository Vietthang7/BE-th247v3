package repo

import (
	"context"
	"intern_247/app"
	"intern_247/models"
)

type Notification models.Notification

func (notify *Notification) Create() error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Model(&models.Notification{}).Create(&notify).Error
}
