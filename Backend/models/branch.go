package models

import "github.com/google/uuid"

type Branch struct {
	Model     `gorm:"embedded"`
	Name      string     `gorm:"name"`
	Address   string     `json:"address,omitempty"`
	IsActive  *bool      `gorm:"default:false" json:"is_active,omitempty"`
	TotalUser int64      `json:"total_user,omitempty" gorm:"-"`
	CenterId  *uuid.UUID `json:"-"`
	UserId    uuid.UUID  `json:"-"` //id user created
	Center    *Center    `json:"center,omitempty"`
}
