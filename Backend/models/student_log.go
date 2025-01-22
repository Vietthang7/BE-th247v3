package models

import "github.com/google/uuid"

type StudentLog struct {
	Model
	StudentId uuid.UUID `json:"student_id"`
	Action    string    `json:"action"`
	UserId    uuid.UUID `json:"-"`
	CreatedBy *User     `gorm:"foreignKey:UserId" json:"created_by"`
}
