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

		if len(u.TimeSlots) > 0 && len(u.ShortShifts) > 0 {
			var schedule StudentSchedule
			if err := schedule.CreateByClassroom(tx, u.StudentId, u.CenterId); err != nil {
				logrus.Error(err)
				return err
			}

			if err := CreateStudentScheduleData(tx, *u, schedule.ID, u.TimeSlots, u.ShortShifts); err != nil {
				logrus.Error(err)
				return err
			}
		}

		if len(u.SubjectIds) > 0 {
			for _, subjectID := range u.SubjectIds {
				studentSubject := StudentSubject{
					StudentId: u.StudentId,
					SubjectId: subjectID,
				}

				if err := studentSubject.Create(tx); err != nil {
					logrus.Error(err)
					return err
				}
			}
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

func (sn *StudyNeeds) Update(studentID uuid.UUID, centerID uuid.UUID) error {
	var existingStudyNeeds StudyNeeds
	if err := app.Database.DB.Where("student_id = ? AND center_id = ?", studentID, centerID).First(&existingStudyNeeds).Error; err != nil {
		logrus.Error(err)
		return fmt.Errorf("%s", "Study needs not found")
	}

	if sn.StudyGoals != "" {
		existingStudyNeeds.StudyGoals = sn.StudyGoals
	}
	if sn.TeacherRequirements != "" {
		existingStudyNeeds.TeacherRequirements = sn.TeacherRequirements
	}
	if sn.IsOnlineForm != nil {
		existingStudyNeeds.IsOnlineForm = sn.IsOnlineForm
	}
	if sn.IsOfflineForm != nil {
		existingStudyNeeds.IsOfflineForm = sn.IsOfflineForm
	}
	if sn.StudyingStartDate != nil {
		existingStudyNeeds.StudyingStartDate = sn.StudyingStartDate
	}
	if sn.BranchId != nil {
		existingStudyNeeds.BranchId = sn.BranchId
	}

	return app.Database.DB.Save(&existingStudyNeeds).Error
}
