package repo

import (
	"context"
	"intern_247/app"
	"intern_247/models"

	"gorm.io/gorm"
)

type TimeSlot models.TimeSlot

func TsCreateManyTimeSlot(tx *gorm.DB, entries ...models.TimeSlot) error {
	return tx.Create(&entries).Error
}
func TsDeleteTimeSlot(tx *gorm.DB) error {
	return tx.Unscoped().Delete(&models.TimeSlot{}).Error
	//Unscoped  sẽ xóa vĩnh viễn
}

func FindTimeSlots(DB *gorm.DB) ([]models.TimeSlot, error) {
	var (
		entries     []models.TimeSlot
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).Find(&entries)
	return entries, err.Error
}
