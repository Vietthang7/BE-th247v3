package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Document struct {
	Model      `gorm:"embedded"`
	Name       string         `json:"name" gorm:"size:250"`
	Type       int            `json:"type"` // internal/consts/subject.go
	CategoryID uuid.UUID      `json:"category_id" gorm:"default:null"`
	Metadata   datatypes.JSON `json:"metadata"`
	CenterID   *uuid.UUID     `json:"-"`
	CreatedBy  uuid.UUID      `json:"-"`
	Creator    *User          `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Category   *DocsCategory  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}
