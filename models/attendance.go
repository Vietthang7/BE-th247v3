package models

import "github.com/google/uuid"

type SessionAttendance struct {
	ClassId        uuid.UUID      `json:"class_id"`
	ClassSessionId uuid.UUID      `json:"class_session_id"`
	StudentId      uuid.UUID      `json:"student_id"`
	Note           string         `json:"note,omitempty"`
	UserId         *uuid.UUID     `json:"user_id"`
	CreatedBy      *User          `json:"created_by,omitempty" gorm:"foreignKey:UserId"`
	Class          *Class         `gorm:"foreignKey:ClassId" json:"class,omitempty"`
	Session        *ScheduleClass `gorm:"foreignKey:ClassSessionId" json:"session,omitempty"`
	Student        *Student       `gorm:"foreignKey:StudentId" json:"student,omitempty"`
}
