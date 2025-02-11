package repo

import (
	"gorm.io/gorm"
	"intern_247/models"
)

func TsDeleteShift(tx *gorm.DB) error {
	return tx.Unscoped().Delete(&models.Shift{}).Error
}
