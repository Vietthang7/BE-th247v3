package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Model             `gorm:"embedded"`
	Username          string     `json:"username,omitempty"`
	Avatar            string     `json:"avatar,omitempty"`
	SignatureImg      string     `json:"signature_img,omitempty"`
	Signature         string     `json:"signature,omitempty"`
	PasswordHash      string     `json:"-"`
	FullName          string     `json:"full_name,omitempty"`
	Gender            string     `json:"gender,omitempty"`
	DOB               *time.Time `json:"dob,omitempty"` // ngày sinh
	Phone             string     `json:"phone,omitempty"`
	Email             string     `gorm:"index:idx_name,unique" json:"email,omitempty"`
	Introduction      string     `json:"introduction,omitempty"` // Giới thiệu
	IsActive          *bool      `json:"is_active,omitempty" gorm:"default:true"`
	SalaryType        int64      `json:"salary_type,omitempty"` // Cách tính lương
	Salary            int64      `json:"salary,omitempty"`
	Position          int64      `json:"position,omitempty"` // internal/consts/role.go
	RoleId            int64      `json:"role_id,omitempty"`
	EmailVerified     bool       `json:"-" gorm:"default:false"`
	CenterId          *uuid.UUID `gorm:"default:null" json:"-"`
	BranchId          *uuid.UUID `json:"branch_id,omitempty"`
	BranchName        string     `json:"branch_name,omitempty" gorm:"-"`
	OrganStructId     *uuid.UUID `json:"organ_struct_id,omitempty"`
	OrganStructName   string     `json:"organ_struct_name,omitempty" gorm:"-"`
	ParentOrganName   string     `json:"parent_organ_name,omitempty" gorm:"-"`   // Tên Cơ cấu tổ chức cha
	PermissionGrpName string     `json:"permission_grp_name,omitempty" gorm:"-"` // tên nhóm phân quyền
	PermissionGrpId   *uuid.UUID `json:"permission_grp_id,omitempty"`            // id nhóm phân quyền
	//Center            *Center    `gorm:"foreignKey:CenterId" json:"center,omitempty"`
	Center *Center `json:"center,omitempty"`
	Branch *Branch `json:"branch,omitempty"`
	//OrganStruct       *OrganStruct `gorm:"foreignKey:OrganStructId" json:"organ_struct"`
	OrganStruct *OrganStruct `json:"organ_struct,omitempty"`
}
