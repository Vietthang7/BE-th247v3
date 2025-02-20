package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type TeachingSchedule struct {
	Model     `gorm:"embedded"`
	UserId    uuid.UUID `json:"user_id"`
	CenterId  uuid.UUID `json:"center_id,omitempty"`
	IsOnline  *bool     `json:"is_online"`
	IsOffline *bool     `json:"is_offline"`
	Notes     string    `json:"notes"`

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	//TimeSlots []TimeSlot  `json:"time_slots"`
	TimeSlots datatypes.JSON `json:"time_slots" gorm:"-"`
	//UserShifts []Shift    `json:"user_shifts"`
	UserShifts datatypes.JSON `json:"user_shifts" gorm:"-"`
}
