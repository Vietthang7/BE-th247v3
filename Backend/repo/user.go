package repo

import (
	"context"
	"intern_247/app"
	"intern_247/models"
)

func VerifyUserEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Update("EmailVerified", true).Error
}
