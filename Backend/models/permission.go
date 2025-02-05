package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Permission struct {
	Model    `gorm:"embedded"`
	TagId    uuid.UUID `json:"tag_id"`
	SubTagId uuid.UUID `json:"sub_tag_id"`
	Name     string    `json:"name,omitempty"`
	Action   string    `json:"action,omitempty"`
	Subject  string    `json:"subject,omitempty"`
}
type PermissionTag struct {
	Model            `gorm:"embedded"`
	ParentTagId      *uuid.UUID     `json:"parent_tag_id,omitempty"`
	Name             string         `json:"name"`
	Key              string         `json:"key"`
	Permissions      datatypes.JSON `json:"permissions,omitempty" gorm:"-"`
	PermissionIds    []uuid.UUID    `json:"-" gorm:"-"`
	SubTags          datatypes.JSON `json:"sub_tags,omitempty" gorm:"-"`
	TotalPermissions int64          `json:"total_permissions" gorm:"-"`
}
type PermissionGroup struct {
	Model         `gorm:"embedded"`
	Name          string         `json:"name"`
	IsActive      *bool          `json:"is_active" gorm:"default:true"`
	SelectAll     *bool          `json:"select_all" gorm:"default:false"`
	PermissionIds datatypes.JSON `json:"permission_ids"`
	CenterId      *uuid.UUID     `json:"-"`
	Tags          datatypes.JSON `json:"tags"`
}
type CustomPermissionTag struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	CountSelected int64     `json:"countSelected"`
}
