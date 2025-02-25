package repo

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/utils"
	"reflect"
	"time"
)

func UpdateRoomScheduleData(tx *gorm.DB, origin models.RoomSchedule, slots []models.TimeSlot, shortShifts []models.ShortShift) error {
	var (
		err     error
		slotMap = make(map[uuid.UUID]uuid.UUID)
	)
	if !reflect.DeepEqual(origin.TimeSlots, slots) {
		if len(slots) == 0 { // If updated time slot length is zero then delete all current time slots in DB
			if err = TsDeleteTimeSlot(tx.Where("schedule_id = ?", origin.ID)); err != nil {
				logrus.Error(err)
				return err
			}
		} else {
			if len(origin.TimeSlots) > 0 {
				if err = TsDeleteTimeSlot(tx.Where("schedule_id = ?", origin.ID)); err != nil {
					logrus.Error(err)
					return err
				}
			}
			for i, v := range slots {
				slots[i].ID = uuid.New()
				slots[i].ScheduleId = origin.ID
				slots[i].CenterId = origin.CenterId
				slots[i].ClassroomId = &origin.ClassroomId
				slotMap[v.WorkSessionId] = slots[i].ID
			}
			if err = TsCreateManyTimeSlot(tx, slots...); err != nil {
				logrus.Error(err)
				return err
			}
		}
	}
	if !reflect.DeepEqual(origin.ShortShifts, shortShifts) {
		if len(shortShifts) == 0 { // If updated time slot length is zero then delete all current time slots in DB
			if err = TsDeleteShift(tx.Where("schedule_id = ?", origin.ID)); err != nil {
				logrus.Error(err)
				return err
			}
		} else {
			if len(origin.ShortShifts) > 0 {
				if err = TsDeleteShift(tx.Where("schedule_id = ?", origin.ID)); err != nil {
					logrus.Error(err)
					return err
				}
			}

			var (
				shift  models.Shift
				shifts []models.Shift
			)
			for _, v := range shortShifts {
				shift.WorkSessionId = v.WorkSessionId
				shift.CenterId = *origin.CenterId
				shift.ScheduleId = origin.ID
				shift.ClassroomId = &origin.ClassroomId
				shift.TimeSlotId = slotMap[v.WorkSessionId]
				shift.Date = time.Now()
				shift.Type = consts.Room
				for _, day := range v.DayOfWeek {
					shift.ID = uuid.New()
					shift.DayOfWeek = day
					shifts = append(shifts, shift)
				}
			}
			if err = TsCreateManyShift(tx, shifts...); err != nil {
				logrus.Error(err)
				return err
			}
		}
	}
	return nil
}
func ShortenRoomShifts(entry *models.RoomSchedule) {
	var (
		shiftMark = make(map[uuid.UUID][]int)
	)
	for _, v := range entry.RoomShifts {
		if !utils.Contains(shiftMark[v.WorkSessionId], v.DayOfWeek) {
			shiftMark[v.WorkSessionId] = append(shiftMark[v.WorkSessionId], v.DayOfWeek)
		}
	}
	for key, val := range shiftMark {
		entry.ShortShifts = append(entry.ShortShifts, models.ShortShift{WorkSessionId: key, DayOfWeek: val})
	}
}

func TsCreateRoomScheduleByClassroom(tx *gorm.DB, classroomId, centerId uuid.UUID) (id uuid.UUID, err error) {
	entry := models.RoomSchedule{
		ClassroomId: classroomId,
		CenterId:    &centerId,
	}
	err = tx.Create(&entry).Error
	return entry.ID, err
}
func CreateRoomScheduleData(tx *gorm.DB, classroom models.Classroom, scheduleId uuid.UUID, slots []models.TimeSlot, shortShifts []models.ShortShift) (err error) {
	// Classroom&TimeSlot creating data
	var (
		slotMap = make(map[uuid.UUID]uuid.UUID)
	)
	for i, v := range slots {
		slots[i].ID = uuid.New()
		slots[i].ScheduleId = scheduleId
		slots[i].ClassroomId = &classroom.ID
		slots[i].CenterId = classroom.CenterId
		slotMap[v.WorkSessionId] = slots[i].ID
	}
	if err = TsCreateManyTimeSlot(tx, slots...); err != nil {
		logrus.Error(err)
		return
	} else {
		// Classroom&Shift creating data
		var (
			shifts []models.Shift
			shift  models.Shift
		)
		for _, v := range shortShifts {
			shift.WorkSessionId = v.WorkSessionId
			shift.ClassroomId = &classroom.ID
			shift.ScheduleId = scheduleId
			shift.CenterId = *classroom.CenterId
			shift.Type = consts.Room
			shift.Date = time.Now()
			shift.TimeSlotId = slotMap[v.WorkSessionId]
			for _, day := range v.DayOfWeek {
				shift.ID, _ = uuid.NewUUID()
				shift.DayOfWeek = day
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
func TsDeleteRoomSchedule(tx *gorm.DB) error {
	return tx.Delete(&models.RoomSchedule{}).Error
}
