package models

import (
	"gorm.io/datatypes"
)

type Center struct {
	Model       `gorm:"embedded"`
	Name        string         `json:"name,omitempty"`
	Email       string         `json:"email,omitempty"`
	Type        string         `json:"type,omitempty"`
	Phone       string         `json:"phone,omitempty"`
	Address     string         `json:"address,omitempty"`
	ProvinceId  *int64         `json:"province_id,omitempty"`
	DistrictId  *int64         `json:"district_id,omitempty"`
	WardId      *int64         `json:"ward_id,omitempty"`
	Highlight   datatypes.JSON `json:"highlight,omitempty"`  // Điểm nổi bật
	Curriculum  datatypes.JSON `json:"curriculum,omitempty"` // Chương trình giảng dạy
	Description datatypes.JSON `json:"description,omitempty"`
	Age         datatypes.JSON `json:"age,omitempty"`
	IsActive    *bool          `json:"is_active,omitempty"`
	Note        string         `json:"note,omitempty"`
	Logo        string         `json:"logo,omitempty"`
	Favicon     string         `json:"favicon,omitempty"`
	CoverImg    string         `json:"cover_img,omitempty"`
	Subjects    datatypes.JSON `json:"subjects,omitempty"` // Môn học
	// UserId      *uuid.UUID     `gorm:"unique;default:null" json:"-"`            //id user owner
	// User        *User          `gorm:"foreignKey:UserId" json:"user,omitempty"` //id user owner
	Domain string `json:"domain,omitempty"`
}
