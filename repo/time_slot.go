package repo

import (
	"intern_247/models"

	"gorm.io/gorm"
)

type TimeSlot models.TimeSlot

func TsDeleteTimeSlot(tx *gorm.DB) error {
	return tx.Unscoped().Delete(&models.TimeSlot{}).Error
}

func TsCreateManyTimeSlot(tx *gorm.DB, entries ...models.TimeSlot) error {
	return tx.Create(&entries).Error
}
