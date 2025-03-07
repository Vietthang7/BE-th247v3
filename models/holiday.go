package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Holiday struct {
	Model       `gorm:"embedded"`
	Title       string         `gorm:"type:varchar(250)" json:"title"`
	StartDay    datatypes.Date `json:"start_day"`
	EndDay      datatypes.Date `json:"end_day"`
	Description *string        `gorm:"size:1000,default" json:"description,omitempty"`
	BranchId    *uuid.UUID     `json:"-"`
	CenterId    *uuid.UUID     `gorm:"default:null" json:"-"`
	UserId      uuid.UUID      `json:"-"`
	//Branch      *Branch        `gorm:"foreignKey:BranchId" json:"branch,omitempty"`
	Branch *Branch `json:"branch,omitempty"`
	//User        *User          `gorm:"foreignKey:UserId" json:"user,omitempty"`
	User *User `json:"user,omitempty"`
	//Center      *Center        `gorm:"foreignKey:CenterId" json:"center,omitempty"`
	Center *Center `json:"center,omitempty"`
}
