package models

import "github.com/google/uuid"

type RoomSchedule struct {
	Model       `gorm:"embedded"`
	ClassroomId uuid.UUID  `json:"classroom_id"`
	CenterId    *uuid.UUID `json:"center_id,omitempty"`

	TimeSlots   []TimeSlot   `json:"time_slots" gorm:"foreignKey:ScheduleId"`
	RoomShifts  []Shift      `json:"-" gorm:"foreignKey:ScheduleId"`
	ShortShifts []ShortShift `json:"short_shifts" gorm:"-"`
}
