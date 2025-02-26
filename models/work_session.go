package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkSession struct {
	Model     `gorm:"embedded"`
	Title     string     `gorm:"type:varchar(250)" json:"title" `
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	IsActive  *bool      `gorm:"default:false" json:"is_active,omitempty"`
	BranchId  *uuid.UUID `gorm:"default:null" json:"branch_id,omitempty"`
	CenterId  *uuid.UUID `gorm:"default:null" json:"-"`
	UserId    uuid.UUID  `json:"-"`
	Branch    *Branch    `json:"branch,omitempty"`
	User      *User      `json:"user,omitempty"`
	Center    *Center    `json:"center,omitempty"`
}
