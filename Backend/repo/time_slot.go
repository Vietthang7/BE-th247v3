package repo

import (
	"gorm.io/gorm"
	"intern_247/models"
)

func TsDeleteTimeSlot(tx *gorm.DB) error {
	return tx.Unscoped().Delete(&models.TimeSlot{}).Error
}
