package models

import (
	"github.com/google/uuid"
)

type DocsCategory struct {
	Model       `gorm:"embedded"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty" gorm:"type:longtext"`
	IsActive    bool       `json:"is_active,omitempty"`
	CenterID    *uuid.UUID `json:"-"`
	CreatedBy   uuid.UUID  `json:"-"`
	TotalDocs   int64      `json:"total_docs,omitempty" gorm:"-"`
	Creator     *User      `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}
