package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TeachingSchedule struct {
	Model     `gorm:"embedded"`
	UserId    uuid.UUID `json:"user_id"`
	CenterId  uuid.UUID `json:"center_id,omitempty"`
	SubjectId uuid.UUID `json:"subject_id"` // Thêm SubjectId
	IsOnline  *bool     `json:"is_online"`
	IsOffline *bool     `json:"is_offline"`
	Notes     string    `json:"notes"`

	StartDate  time.Time      `json:"start_date"`
	EndDate    time.Time      `json:"end_date"`
	TimeSlots  datatypes.JSON `json:"time_slots" gorm:"-"`
	UserShifts datatypes.JSON `json:"user_shifts" gorm:"-"`
}

type CreateTeachScheForm struct {
	UserId     uuid.UUID      `json:"user_id"`
	SubjectId  uuid.UUID      `json:"subject_id"` // Thêm SubjectId
	IsOnline   *bool          `json:"is_online"`
	IsOffline  *bool          `json:"is_offline"`
	Notes      string         `json:"notes"`
	StartDate  string         `json:"start_date"`
	EndDate    string         `json:"end_date"`
	TimeSlots  datatypes.JSON `json:"time_slots" gorm:"-"`
	UserShifts datatypes.JSON `json:"user_shifts" gorm:"-"`
}
