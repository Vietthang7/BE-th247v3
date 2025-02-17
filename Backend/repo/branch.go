package repo

import (
	"context"
	"intern_247/app"
	"intern_247/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateBranch(DB *gorm.DB, entry *models.Branch) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return DB.WithContext(ctx).Create(&entry).Error
}

func FirstBranch(DB *gorm.DB) (models.Branch, error) {
	var (
		entry       models.Branch
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).First(&entry)
	return entry, err.Error
}

func ListBranches(DB *gorm.DB, centerId uuid.UUID) ([]models.Branch, error) {
	var (
		branches    []models.Branch
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()

	// Truy vấn tất cả chi nhánh theo centerId
	err := DB.WithContext(ctx).
		Where("center_id = ?", centerId).
		Find(&branches).Error

	return branches, err
}

func UpdateBranch(DB *gorm.DB, branchId uuid.UUID, updatedBranch *models.Branch) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()

	// Cập nhật thông tin chi nhánh
	err := DB.WithContext(ctx).Model(&models.Branch{}).
		Where("id = ?", branchId).
		Updates(updatedBranch).Error

	return err
}

func DeleteBranch(DB *gorm.DB, branchId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()

	// Thực hiện xóa chi nhánh
	err := DB.WithContext(ctx).Where("id = ?", branchId).Delete(&models.Branch{}).Error
	return err
}
