package repo

import (
	"context"
	"fmt"
	"intern_247/app"
	"intern_247/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StudyNeeds models.StudyNeeds
type ListStudyNeeds []models.StudyNeeds

func CheckStudentExists(studentID uuid.UUID) error {
	var student Student
	if err := app.Database.DB.Where("id = ?", studentID).First(&student).Error; err != nil {
		return err
	}
	return nil
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

func CheckStudentHasBranch(studentID uuid.UUID) error {
	var existingStudyNeeds StudyNeeds
	if err := app.Database.DB.Where("student_id = ?", studentID).First(&existingStudyNeeds).Error; err == nil {
		return fmt.Errorf("%s", "học viên đã được gán chi nhánh trước đó")
	}
	return nil
}

func (u *StudyNeeds) Create() error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()

	return app.Database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&u).Error; err != nil {
			logrus.Error(err)
			return err
		}

		return nil
	})
}

func GetStudyNeedsByID(studyNeedsID uuid.UUID, centerID uuid.UUID) (*StudyNeeds, error) {
	var studyNeeds StudyNeeds
	err := app.Database.DB.
		Where("id = ? AND center_id = ?", studyNeedsID, centerID).
		First(&studyNeeds).Error // Lấy 1 bản ghi duy nhất

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &studyNeeds, nil
}

func GetAllStudyNeeds(centerID uuid.UUID) ([]StudyNeeds, error) {
	var studyNeeds []StudyNeeds
	if err := app.Database.DB.
		Where("center_id = ?", centerID).
		Find(&studyNeeds).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return studyNeeds, nil
}

func (sn *StudyNeeds) Update() error {
	return app.Database.DB.Model(&StudyNeeds{}).Where("id = ?", sn.ID).Updates(sn).Error
}
