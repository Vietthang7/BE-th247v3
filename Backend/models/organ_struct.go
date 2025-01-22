package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OrganStruct struct {
	Model             `gorm:"embedded"`
	Name              string         `json:"name"`
	ParentId          string         `json:"parent_id,omitempty" gorm:"default:null"`
	ParentName        string         `json:"parent_name,omitempty" gorm:"-"`
	IsActive          *bool          `json:"is_active,omitempty"`
	PermissionGrpName string         `json:"permission_grp_name,omitempty" gorm:"-"` // tên nhóm phân quyền
	PermissionGrpId   *uuid.UUID     `json:"permission_grp_id,omitempty"`            // id nhóm phân quyền
	PermissionIds     datatypes.JSON `json:"-"  gorm:"-"`
	TotalUser         int64          `json:"total_user,omitempty" gorm:"-"`
	UserId            uuid.UUID      `json:"-"` //user created
	CenterId          *uuid.UUID     `json:"-"`
	Center            *Center        `json:"center,omitempty"`
}
