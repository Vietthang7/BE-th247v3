package models

import "github.com/google/uuid"

type SalaryStatement struct {
	Model      `gorm:"embedded"`
	SalaryType uint8      `json:"salary_type"` // 1. hour , 2. session, 3. full course
	ObjectType uint8      `json:"object_type"` // 1. teaching asistant, 2. teacher, 3. teacher and teaching asistant
	IsActive   *bool      `gorm:"default:true" json:"is_active"`
	CenterId   *uuid.UUID `json:"-"`
	Title      *string    `gorm:"default:null" json:"title"`
	IsDefault  bool       `json:"-" gorm:"default:false"`
}
