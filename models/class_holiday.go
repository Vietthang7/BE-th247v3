package models

import (
	"time"

	"github.com/google/uuid"
)

type ClassHoliday struct {
	Model            `gorm:"embedded"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	StartAt          time.Time         `json:"start_at"`
	EndAt            time.Time         `json:"end_at"`
	ClassID          uuid.UUID         `json:"class_id"`
	IsAuto           bool              `gorm:"default:false" json:"is_auto"`
	IsChanged        *bool             `gorm:"default:false" json:"is_changed"`
	CenterId         uuid.UUID         `json:"-"`
	Center           *Center           `gorm:"foreignKey:CenterId" json:"-"`
	ScheduleHolidays []ScheduleHoliday `json:"schedule_holidays,omitempty"`
}
