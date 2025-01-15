package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Center struct {
	Model       `gorm:"embedded"`
	Name        string         `json:"name,omitempty"`
	Email       string         `json:"email,omitempty"`
	Type        string         `json:"type,omitempty"`
	Phone       string         `json:"phone,omitempty"`
	Address     string         `json:"address,omitempty"`
	ProvinceID  string         `json:"province,omitempty"` // ID tỉnh
	DistrictID  string         `json:"district,omitempty"` // ID huyện
	WardID      string         `json:"ward,omitempty"`
	Highlight   datatypes.JSON `json:"highlight,omitempty"`  // Điểm nổi bật
	Curriculum  datatypes.JSON `json:"curriculum,omitempty"` // Chương trình giảng dạy
	Description datatypes.JSON `json:"description,omitempty"`
	Age         datatypes.JSON `json:"age,omitempty"` // Độ tuổi mục tiêu của trung tâm
	IsActive    *bool          `json:"is_active,omitempty"`
	Note        string         `json:"note,omitempty"`
	Logo        string         `json:"logo,omitempty"`
	Favicon     string         `json:"favicon,omitempty"`   // Favicon đại diện cho trung tâm
	CoverImg    string         `json:"cover_img,omitempty"` // Hình ảnh bìa của trung tâm
	Subjects    datatypes.JSON `json:"subjects,omitempty"`
	UserID      *uuid.UUID     `gorm:"unique;default:null;" json:"-"`
	//User        *User          `gorm:"foreignKey:UserId" json:"user,omitempty"` //id user
	Domain string `json:"domain,omitempty"`
}
