package repo

import (
	"context"
	"intern_247/app"
	"intern_247/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StudyNeeds models.StudyNeeds
type ListStudyNeeds []models.StudyNeeds

func (u *StudyNeeds) Create() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(&u).Error; err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	})
}

func GetStudyNeedsByStudentID(studentID uuid.UUID, centerID uuid.UUID) ([]StudyNeeds, error) {
	var studyNeeds []StudyNeeds
	err := app.Database.DB.
		Where("student_id = ? AND center_id = ?", studentID, centerID).
		Find(&studyNeeds).Error // Dùng Find thay vì First

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return studyNeeds, nil
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
