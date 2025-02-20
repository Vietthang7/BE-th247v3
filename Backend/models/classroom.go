package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Classroom struct {
	Model    `gorm:"embedded"`
	BranchId *uuid.UUID `json:"branch_id,omitempty" gorm:"not null"`
	CenterId *uuid.UUID `json:"center_id,omitempty"`
	Name     string     `json:"name" gorm:"not null"`
	IsOnline *bool      `json:"is_online,omitempty" gorm:"default:true"` // Loại phòng: true (Phòng online), false (Phòng offline)
	RoomType string     `json:"room_type,omitempty"`                     // Kiểu phòng học: Google Meet, Zoom, Live class
	//Metadata struct {
	//	Notes string `json:"notes"`
	//	Link  string `json:"link"`   // Link học trực tuyến
	//}
	Metadata datatypes.JSON `json:"metadata,omitempty"`
	Slots    *int64         `json:"slots,omitempty" gorm:"default:null"` // Sức chứa
	IsActive *bool          `json:"is_active,omitempty" gorm:"default:true"`

	Branch   *Branch       `json:"branch,omitempty" gorm:"foreignKey:BranchId"`
	Schedule *RoomSchedule `json:"schedule,omitempty" gorm:"foreignKey:ClassroomId"`
}
