package repo

import (
	"intern_247/consts"
	"intern_247/models"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StudentSchedule models.StudentSchedule

func CreateStudentScheduleData(tx *gorm.DB, studyNeeds StudyNeeds, scheduleId uuid.UUID,
	slots []models.TimeSlot, shortShifts []models.ShortShift) (err error) {

	// Student & TimeSlot creating data
	var (
		slotMap = make(map[uuid.UUID]uuid.UUID)
	)

	// Lấy ngày bắt đầu học
	startDate := time.Now() // Mặc định là ngày hiện tại
	if studyNeeds.StudyingStartDate != nil {
		startDate = *studyNeeds.StudyingStartDate
	}

	for i, v := range slots {
		slots[i].ID = uuid.New()
		slots[i].ScheduleId = scheduleId
		slots[i].StudentId = &studyNeeds.StudentId
		slots[i].CenterId = &studyNeeds.CenterId
		slotMap[v.WorkSessionId] = slots[i].ID
	}

	if err = TsCreateManyTimeSlot(tx, slots...); err != nil {
		logrus.Error(err)
		return
	} else {
		// Student & Shift creating data
		var (
			shifts []models.Shift
			shift  models.Shift
		)

		for _, v := range shortShifts {
			shift.WorkSessionId = v.WorkSessionId
			shift.ScheduleId = scheduleId
			shift.StudentId = &studyNeeds.StudentId
			shift.CenterId = studyNeeds.CenterId
			shift.Type = consts.StudentShift
			shift.TimeSlotId = slotMap[v.WorkSessionId]

			for _, day := range v.DayOfWeek {
				shift.ID, _ = uuid.NewUUID()
				shift.DayOfWeek = day
				shift.Date = getNextDayOfWeek(startDate, day)
				shifts = append(shifts, shift)
			}
		}

		if err = TsCreateManyShift(tx, shifts...); err != nil {
			logrus.Error(err)
			return
		}
	}

	return nil
}

func getNextDayOfWeek(startDate time.Time, targetDay int) time.Time {
	startDay := int(startDate.Weekday()) + 1
	if startDay == 8 {
		startDay = 1
	}

	daysToAdd := targetDay - startDay
	if daysToAdd < 0 {
		daysToAdd += 7
	}

	return startDate.AddDate(0, 0, daysToAdd)
}

func (u *StudentSchedule) CreateByClassroom(tx *gorm.DB, studentId, centerId uuid.UUID) (err error) {
	u.StudentId = studentId
	u.CenterId = &centerId

	err = tx.Create(&u).Error
	return err
}
