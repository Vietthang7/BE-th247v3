package models

import "github.com/google/uuid"

type ContactChannel struct {
	Model       `gorm:"embedded"`
	Name        string     `json:"name" gorm:"size:250"`
	Description string     `json:"description,omitempty" gorm:"type:longtext"`
	IsActive    bool       `json:"is_active"`
	CenterID    *uuid.UUID `json:"-"`
}
