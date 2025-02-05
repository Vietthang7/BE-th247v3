package models

import "github.com/google/uuid"

type CareAssignment struct {
	Model         `gorm:"embedded"`
	Type          int          `json:"type" gorm:"default:1"` // internal/consts/config.go
	IsTemplate    bool         `json:"-" gorm:"default:false"`
	OrganStructID *uuid.UUID   `json:"organ_struct_id" gorm:"default:null"`
	OrganStruct   *OrganStruct `json:"organ_struct" gorm:"foreignKey:OrganStructID"`
	CenterID      *uuid.UUID   `json:"-"`
}
