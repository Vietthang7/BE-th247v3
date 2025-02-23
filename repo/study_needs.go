package repo

import (
	"context"
	"errors"
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

		if len(u.SubjectIds) == 1 {
			studentSubject := StudentSubject{
				StudentId: u.StudentId,
				SubjectId: u.SubjectIds[0], // Lấy subject_id duy nhất
			}

			if err := studentSubject.Create(tx); err != nil {
				logrus.Error(err)
				return err
			}
		}

		return nil
	})
}

func CheckSubjectsExist(subjectIDs []uuid.UUID) error {
	if len(subjectIDs) == 0 {
		return nil
	}

	var count int64
	if err := app.Database.DB.Model(&models.Subject{}).
		Where("id IN (?)", subjectIDs).
		Count(&count).Error; err != nil {
		return errors.New("lỗi khi kiểm tra môn học")
	}

	if count != int64(len(subjectIDs)) {
		return errors.New("một hoặc nhiều môn học không tồn tại")
	}

	return nil
}

func GetStudyNeedsByID(studyNeedsID uuid.UUID, centerID uuid.UUID) (*StudyNeeds, error) {
	var studyNeeds StudyNeeds

	// Lấy thông tin StudyNeeds
	err := app.Database.DB.
		Where("id = ? AND center_id = ?", studyNeedsID, centerID).
		First(&studyNeeds).Error
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Lấy danh sách subject_id từ bảng student_subjects
	var subjectIds []uuid.UUID
	err = app.Database.DB.
		Table("student_subjects").
		Select("subject_id").
		Where("student_id = ?", studyNeeds.StudentId).
		Pluck("subject_id", &subjectIds).Error
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Gán subject_ids vào StudyNeeds
	studyNeeds.SubjectIds = subjectIds

	return &studyNeeds, nil
}

func GetAllStudyNeeds(centerID uuid.UUID) ([]StudyNeeds, error) {
	var studyNeedsList []StudyNeeds

	// Lấy danh sách StudyNeeds
	err := app.Database.DB.
		Where("center_id = ?", centerID).
		Find(&studyNeedsList).Error
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Lấy subject_id cho từng StudyNeeds
	for i, studyNeeds := range studyNeedsList {
		var subjectIds []uuid.UUID
		err = app.Database.DB.
			Table("student_subjects").
			Select("subject_id").
			Where("student_id = ?", studyNeeds.StudentId).
			Pluck("subject_id", &subjectIds).Error
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		studyNeedsList[i].SubjectIds = subjectIds
	}

	return studyNeedsList, nil
}

func (sn *StudyNeeds) Update(studyNeedsID uuid.UUID, centerID uuid.UUID) error {
	var existingStudyNeeds StudyNeeds
	if err := app.Database.DB.Where("id = ? AND center_id = ?", studyNeedsID, centerID).First(&existingStudyNeeds).Error; err != nil {
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

	if len(sn.SubjectIds) == 1 {
		subjectID := sn.SubjectIds[0]

		if err := app.Database.DB.Where("student_id = ?", existingStudyNeeds.StudentId).Delete(&StudentSubject{}).Error; err != nil {
			logrus.Error("Failed to delete old StudentSubjects:", err)
			return err
		}

		newStudentSubject := StudentSubject{
			StudentId: existingStudyNeeds.StudentId,
			SubjectId: subjectID,
		}
		if err := app.Database.DB.Create(&newStudentSubject).Error; err != nil {
			logrus.Error("Failed to insert new StudentSubject:", err)
			return err
		}
	}

	return app.Database.DB.Save(&existingStudyNeeds).Error
}
