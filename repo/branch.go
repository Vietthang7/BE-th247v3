package repo

import (
	"context"
	"fmt"
	"intern_247/app"
	"intern_247/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Branch  models.Branch
	Branchs []models.Branch
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

func GetBranchByIdAndCenterId(id, centerId uuid.UUID) (models.Branch, error) {
	var branch models.Branch
	query := app.Database.DB.Where("id = ? AND center_id = ?", id, centerId).First(&branch)
	return branch, query.Error
}

func CheckBranchIsActive(branchID uuid.UUID) error {
	var branch Branch
	if err := app.Database.DB.Where("id = ?", branchID).First(&branch).Error; err != nil {
		return fmt.Errorf("%s", "Không tìm thấy chi nhánh")
	}
	if branch.IsActive == nil || !*branch.IsActive {
		return fmt.Errorf("%s", "Chi nhánh không hoạt động")
	}
	return nil
}
func IsExistBranchInCenter(id, centerId uuid.UUID, isActive *bool) bool {
	var branch models.Branch
	query := app.Database.DB.Where("id = ? AND center_id = ?", id, centerId)
	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}
	query.First(&branch)
	return query.RowsAffected > 0
}
