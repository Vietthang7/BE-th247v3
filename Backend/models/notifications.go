package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Notification struct {
	Model    `gorm:"embedded"`
	Title    string         `json:"title"`
	Content  string         `gorm:"type:text;default:null" json:"content,omitempty"`
	To       uuid.UUID      `json:"to"`
	From     *uuid.UUID     `json:"from"`
	Metadata datatypes.JSON `json:"metadata"`
	IsRead   bool           `json:"is_read" gorm:"default:false"`
}
