package repo

import (
	"intern_247/app"
	"intern_247/models"

	"github.com/google/uuid"
)

func CreateClassHoliday(classHoliday *models.ClassHoliday) error {
	return app.Database.DB.Create(classHoliday).Error
}

func IsClassHolidayExist(name string) bool {
	var count int64
	app.Database.DB.Model(&models.ClassHoliday{}).Where("name = ?", name).Count(&count)
	return count > 0
}

func IsClassIDExist(classID uuid.UUID) bool {
	var count int64
	app.Database.DB.Model(&models.Class{}).Where("id = ?", classID).Count(&count)
	return count > 0
}

func GetListClassHoliday() ([]models.ClassHoliday, error) {
	var classHolidays []models.ClassHoliday
	err := app.Database.DB.Find(&classHolidays).Error
	return classHolidays, err
}

func GetDetailClassHoliday(id uuid.UUID) (*models.ClassHoliday, error) {
	var classHoliday models.ClassHoliday
	err := app.Database.DB.Where("id = ?", id).First(&classHoliday).Error
	if err != nil {
		return nil, err
	}
	return &classHoliday, nil
}

func IsClassHolidayExistByID(id uuid.UUID) (bool, error) {
	var count int64
	err := app.Database.DB.Model(&models.ClassHoliday{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DeleteClassHoliday(id uuid.UUID) error {
	return app.Database.DB.Where("id = ?", id).Delete(&models.ClassHoliday{}).Error
}

// Cập nhật thông tin ngày nghỉ trong database
func UpdateClassHoliday(classHoliday *models.ClassHoliday) error {
	return app.Database.DB.Save(classHoliday).Error
}
