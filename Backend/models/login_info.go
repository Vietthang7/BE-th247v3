package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type LoginInfo struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey"`
	//SsoID        string    `json:"-"`
	CenterID     uuid.UUID `json:"-"`
	Username     string    `json:"username,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Email        string    `json:"email,omitempty"`
	PasswordHash string    `json:"-"`
	RoleId       int64     `json:"role_id,omitempty"`

	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	User    *User    `json:"user,omitempty" gorm:"foreignKey:ID;references:ID"`
	Student *Student `json:"student,omitempty" gorm:"foreignKey:ID;references:ID"`
}
