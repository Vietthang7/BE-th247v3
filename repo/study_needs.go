package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"intern_247/app"
	"intern_247/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StudyNeeds models.StudyNeeds
type ListStudyNeeds []models.StudyNeeds

func (u *StudyNeeds) Create() error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()

	return app.Database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Tạo StudyNeeds
		if err := tx.Create(&u).Error; err != nil {
			logrus.Error(err)
			return err
		}

		// Đảm bảo ID đã được lấy sau khi tạo
		if u.ID == uuid.Nil {
			return fmt.Errorf("failed to get ID after creating StudyNeeds")
		}

		// Tạo lịch học nếu có
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

		// ✅ Thêm StudentSubject với SubjectId
		if len(u.SubjectIds) == 1 {
			studentSubject := StudentSubject{
				StudentId: u.StudentId,
				SubjectId: u.SubjectIds[0], // ✅ Đã có SubjectId trong struct
			}

			if err := tx.Create(&studentSubject).Error; err != nil {
				logrus.Error("Failed to insert student_subjects:", err)
				return err
			}
		}

		return nil
	})
}

func GetStudyNeedsByID(studyNeedsID uuid.UUID, centerID uuid.UUID) (*StudyNeeds, error) {
	var studyNeeds StudyNeeds

	// Lấy thông tin StudyNeeds
	err := app.Database.DB.
		Where("id = ? AND center_id = ?", studyNeedsID, centerID).
		First(&studyNeeds).Error
	if err != nil {
		logrus.Error("StudyNeeds not found:", err)
		return nil, err
	}

	// ⚠️ Lấy danh sách subject_id từ student_subjects theo student_id
	var subjectIds []uuid.UUID
	err = app.Database.DB.
		Table("student_subjects").
		Select("subject_id").
		Where("student_id = ?", studyNeeds.StudentId).
		Pluck("subject_id", &subjectIds).Error
	if err != nil {
		logrus.Error("Failed to fetch subject IDs:", err)
		return nil, err
	}
	studyNeeds.SubjectIds = subjectIds

	// 🔹 Lấy danh sách tất cả student_schedules của học viên
	var studentSchedules []StudentSchedule
	err = app.Database.DB.
		Where("student_id = ?", studyNeeds.StudentId).
		Find(&studentSchedules).Error
	if err != nil {
		logrus.Warn("Student schedules not found for student:", studyNeeds.StudentId)
	}

	// Nếu có lịch học, lấy đầy đủ TimeSlots và Shifts
	var allTimeSlots []models.TimeSlot
	var allShortShifts []models.ShortShift

	for _, schedule := range studentSchedules {
		// 🔹 Lấy danh sách TimeSlots theo schedule_id
		var timeSlots []models.TimeSlot
		err = app.Database.DB.
			Where("schedule_id = ?", schedule.ID).
			Find(&timeSlots).Error
		if err != nil {
			logrus.Error("Failed to fetch time slots for schedule:", schedule.ID, err)
			return nil, err
		}
		allTimeSlots = append(allTimeSlots, timeSlots...)

		// 🔹 Lấy danh sách ShortShifts theo schedule_id
		var rawShortShifts []struct {
			WorkSessionId uuid.UUID
			DayOfWeek     string // Nhận dữ liệu JSON dưới dạng chuỗi
		}

		err = app.Database.DB.
			Table("shifts").
			Select("work_session_id, JSON_ARRAYAGG(day_of_week) AS day_of_week").
			Where("schedule_id = ?", schedule.ID).
			Group("work_session_id").
			Scan(&rawShortShifts).Error
		if err != nil {
			logrus.Error("Failed to fetch short shifts for schedule:", schedule.ID, err)
			return nil, err
		}

		// 🔹 Chuyển đổi JSON string thành slice []int
		for _, raw := range rawShortShifts {
			var days []int
			if err := json.Unmarshal([]byte(raw.DayOfWeek), &days); err != nil {
				logrus.Error("Failed to parse day_of_week JSON:", err)
				continue
			}

			allShortShifts = append(allShortShifts, models.ShortShift{
				WorkSessionId: raw.WorkSessionId,
				DayOfWeek:     days,
			})
		}
	}

	// Gán vào StudyNeeds
	studyNeeds.TimeSlots = allTimeSlots
	studyNeeds.ShortShifts = allShortShifts

	return &studyNeeds, nil
}

func GetAllStudyNeeds(centerID uuid.UUID) ([]StudyNeeds, error) {
	var studyNeedsList []StudyNeeds

	// Lấy danh sách StudyNeeds theo centerID
	err := app.Database.DB.
		Where("center_id = ?", centerID).
		Find(&studyNeedsList).Error
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Lặp qua từng studyNeeds để lấy danh sách subject_ids
	for i := range studyNeedsList {
		var subjectIds []uuid.UUID
		err = app.Database.DB.
			Table("student_subjects").
			Select("subject_id").
			Where("student_id = ?", studyNeedsList[i].StudentId). // ⚠️ Đổi `study_needs_id` thành `student_id`
			Pluck("subject_id", &subjectIds).Error
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		// Cập nhật lại struct studyNeedsList[i]
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
