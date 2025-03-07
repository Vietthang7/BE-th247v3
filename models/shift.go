package models

import (
	"time"

	"github.com/google/uuid"
)

type Shift struct {
	Model         `gorm:"embedded"`
	ScheduleId    uuid.UUID  `json:"-"`
	UserId        *uuid.UUID `json:"-"`
	StudentId     *uuid.UUID `json:"-"`
	ClassroomId   *uuid.UUID `json:"-"`
	CenterId      uuid.UUID  `json:"-"`
	Type          string     `json:"-" gorm:"not null"`
	TimeSlotId    uuid.UUID  `json:"-"`
	WorkSessionId uuid.UUID  `json:"work_session_id"`
	DayOfWeek     int        `json:"day_of_week"`
	Date          time.Time  `json:"-"`
}
type ShortShift struct {
	WorkSessionId uuid.UUID `json:"work_session_id"`
	DayOfWeek     []int     `json:"day_of_week"`
	Date          time.Time `json:"date"`
}
