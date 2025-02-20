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
	CenterId          *uuid.UUID `gorm:"default:null" json:"center_id"`
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
	SubjectIds  []*uuid.UUID `json:"subject_ids,omitempty" gorm:"-"`
	Subjects    []*Subject   `json:"subjects,omitempty" gorm:"many2many:subject_teachers"`

	ProvinceId *int64    `json:"province_id,omitempty"`
	DistrictId *int64    `json:"district_id,omitempty"`
	Province   *Province `json:"province,omitempty" gorm:"foreignKey:ProvinceId"`
	District   *District `json:"district,omitempty" gorm:"foreignKey:DistrictId"`
	Address    string    `json:"address,omitempty"`
}

type DataUserReturn struct {
	ID             uuid.UUID        `json:"id"`
	FirstTimeLogin bool             `json:"first_time_login"`
	RoleId         int64            `json:"role_id"`
	PermissionGrp  *PermissionGroup `json:"permission_grp,omitempty"`
	Position       int64            `json:"position,omitempty"`
	FullName       string           `json:"full_name"`
	Avatar         string           `json:"avatar"`
	Email          string           `json:"email"`
	Domain         string           `json:"domain,omitempty"`
	Phone          string           `json:"phone"`
	Gender         string           `json:"gender"`
	Role           string           `json:"role"`
	BranchID       *uuid.UUID       `json:"branch_id,omitempty"`
	Token          string           `json:"token"`
	//SSO_ID         *string          `json:"sso_id"`
	GoogleId     string `json:"google_id"`
	RefreshToken string `json:"refresh_token"`
}
type CreateUserForm struct {
	ExcelInd        int64        `json:"excel_ind"` //STT trong file excel
	Avatar          string       `json:"avatar"`
	FullName        string       `json:"full_name"`
	Position        int64        `json:"position"` // Vai trò (1: Giảng Viên , 2 : Trợ giảng)
	Email           string       `json:"email"`
	Phone           string       `json:"phone"`
	BranchId        *uuid.UUID   `json:"branch_id"`
	OrganStructId   *uuid.UUID   `json:"organ_struct_id"`   // ID cơ cấu phân quyền tổ chức
	PermissionGrpId *uuid.UUID   `json:"permission_grp_id"` // Id nhóm phân quyền
	Username        string       `json:"username"`
	Password        string       `json:"password"`
	Introduction    string       `json:"introduction"` // Giới thiệu
	IsActive        *bool        `json:"is_active"`
	SalaryType      int64        `json:"salary_type"`
	Salary          int64        `json:"salary"`
	SubjectIds      []*uuid.UUID `json:"subject_ids"`
}
