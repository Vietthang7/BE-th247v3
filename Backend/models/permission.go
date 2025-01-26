package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PermissionGroup struct {
	Model         `gorm:"embedded"`
	Name          string         `json:"name"`
	IsActive      *bool          `json:"is_active" gorm:"default:true"`
	SelectAll     *bool          `json:"select_all" gorm:"default:false"`
	PermissionIds datatypes.JSON `json:"permission_ids"`
	CenterId      *uuid.UUID     `json:"-"`
	Tags          datatypes.JSON `json:"tags"`
}
