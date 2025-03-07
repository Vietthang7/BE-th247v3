package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ScheduleHoliday struct {
	Model           `gorm:"embedded"`
	ScheduleClassId uuid.UUID      `json:"-"`
	ClassHolidayId  uuid.UUID      `json:"-"`
	ClassHoliday    *ClassHoliday  `gorm:"foreignKey:ClassHolidayId" json:"class_holiday,omitempty"`
	ScheduleClass   *ScheduleClass `gorm:"foreignKey:ScheduleClassId" json:"schedule_class,omitempty"`
	Metadata        datatypes.JSON `json:"metadata"`
}
