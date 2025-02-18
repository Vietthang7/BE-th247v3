package repo

import (
	"intern_247/app"
	"intern_247/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateCategory(category *models.Category) (*models.Category, error) {
	query := app.Database.DB.Create(&category)
	return category, query.Error
}

func GetCategoryByNameAndCenterId(name string, centerId uuid.UUID) (models.Category, error) {
	var category models.Category
	db := app.Database.DB.Debug().Select("id", "name").Where("name = ? AND center_id = ?", name, centerId).First(&category)
	return category, db.Error
}

func GetCategoryByIdAndCenterId(id uuid.UUID, centerId uuid.UUID) (models.Category, error) {
	var category models.Category
	query := app.Database.DB.
		Preload("Created", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "full_name")
		}).
		Where("id = ? AND center_id = ?", id, centerId).
		First(&category)
	return category, query.Error
}

func GetCategoriesByCenterIdAndActive(centerId uuid.UUID, isActive *bool) ([]models.Category, error) {
	var categories []models.Category
	query := app.Database.DB.Where("center_id = ?", centerId)

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Preload("Created", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "full_name")
	}).Find(&categories).Error

	return categories, err
}

// Kiểm tra số lượng danh mục con
func CountChildCategories(parentId uuid.UUID) (int64, error) {
	var count int64
	err := app.Database.DB.Model(&models.Category{}).
		Where("parent_id = ?", parentId).
		Count(&count).Error
	return count, err
}

// Kiểm tra xem danh mục có phụ thuộc trong `Curriculums` hoặc `Subjects`
func HasCategoryDependencies(categoryId uuid.UUID) (bool, error) {
	var count int64
	err := app.Database.DB.
		Table("curriculums").
		Where("category_id = ?", categoryId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	err = app.Database.DB.
		Table("subjects").
		Where("category_id = ?", categoryId).
		Count(&count).Error
	return count > 0, err
}

// Xóa danh mục theo ID
func DeleteCategoryById(categoryId uuid.UUID) error {
	return app.Database.DB.Delete(&models.Category{}, "id = ?", categoryId).Error
}

// Cập nhật danh mục
func UpdateCategory(category *models.Category) (*models.Category, error) {
	query := app.Database.DB.Save(category)
	return category, query.Error
}
