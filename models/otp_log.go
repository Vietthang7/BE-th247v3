package models

import (
	"github.com/google/uuid"
	"time"
)

type OTPLog struct {
	Model       `gorm:"embedded"`
	Receiver    string    `gorm:"index" json:"receiver"`
	Code        string    `gorm:"index" json:"code"`
	ExpiredAt   time.Time `json:"expired_at"`
	IsConfirmed bool      `gorm:"default:false;index" json:"is_confirmed"`
	CreatedBy   uuid.UUID `json:"created_by"`
	//User        *User     `gorm:"foreignKey:CreatedBy" json:"created,omitempty"`
}
