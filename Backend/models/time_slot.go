package models

import (
	"github.com/google/uuid"
	"time"
)

// TimeSlot Khung giờ các ca do user config
type TimeSlot struct {
	Model         `gorm:"embedded"`
	UserId        *uuid.UUID `json:"user_id,omitempty"`
	StudentId     *uuid.UUID `json:"-"`
	ClassroomId   *uuid.UUID `json:"classroom_id,omitempty"`
	CenterId      *uuid.UUID `json:"center_id,omitempty"`
	ScheduleId    uuid.UUID  `json:"schedule_id,omitempty"`
	WorkSessionId uuid.UUID  `json:"work_session_id"` // phiên làm việc
	StartTime     *time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
}
