package repo

import (
	"intern_247/models"

	"gorm.io/gorm"
)

type Shift models.Shift

func TsDeleteShift(tx *gorm.DB) error {
	return tx.Unscoped().Delete(&models.Shift{}).Error
}

func TsCreateManyShift(tx *gorm.DB, entries ...models.Shift) error {
	return tx.Create(&entries).Error
}
