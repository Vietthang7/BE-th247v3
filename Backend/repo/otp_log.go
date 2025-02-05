package repo

import (
	"intern_247/app"
	"intern_247/models"
)

func CreateOTPLog(otp *models.OTPLog) error {
	query := app.Database.DB.Create(&otp)
	return query.Error
}

func GetNewestOTPLogByReceiver(email string) (models.OTPLog, error) {
	var otpInfo models.OTPLog
	query := app.Database.DB.Where("receiver = ?", email).Limit(1).First(&otpInfo)
	return otpInfo, query.Error
}
func UpdateOTPLogById(otp *models.OTPLog) error {
	query := app.Database.DB.Where("id = ?", otp.ID).Updates(&otp)
	return query.Error
}
