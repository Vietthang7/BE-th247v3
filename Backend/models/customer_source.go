package models

import "github.com/google/uuid"

type CustomerSource struct {
	Model       `gorm:"embedded"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty" gorm:"type:longtext"`
	IsActive    bool       `json:"is_active,omitempty"`
	CenterID    *uuid.UUID `json:"-"`
}
