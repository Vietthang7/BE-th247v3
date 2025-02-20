package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Curriculum struct {
	Model            `gorm:"embedded"`
	Thumbnail        string         `json:"thumbnail,omitempty"`
	Name             string         `gorm:"type:varchar(250)" json:"name,omitempty"`
	CategoryId       *uuid.UUID     `json:"category_id,omitempty"`
	IsActive         *bool          `gorm:"default:true" json:"is_active,omitempty"`
	Type             int64          `json:"type,omitempty"` //1: free, 2: pay
	DocumentPath     string         `json:"document_path,omitempty"`
	Description      string         `gorm:"type:text" json:"description,omitempty"`
	LearnDescription string         `gorm:"type:text" json:"learn_description,omitempty"`
	DiscountFee      int64          `gorm:"default:0" json:"discount_fee,omitempty"`
	TotalFee         int64          `json:"total_fee,omitempty"`
	InputRequire     string         `gorm:"type:text;default:null" json:"input_require,omitempty"`
	OutputRequire    string         `gorm:"type:text;default:null" json:"output_require,omitempty"`
	Metadata         datatypes.JSON `json:"metadata,omitempty"`
	CenterId         uuid.UUID      `json:"-"`
	CreatedBy        uuid.UUID      `json:"-"`
	IsPaid           *bool          `gorm:"default:false" json:"is_paid,omitempty"`
	Code             string         `gorm:"default:null;index;size:20" json:"code,omitempty"`
	Center           Center         `gorm:"foreignKey:CenterId" json:"-"`
	User             *User          `gorm:"foreignKey:CreatedBy" json:"created_by,omitempty"`
	Subjects         []*Subject     `gorm:"many2many:curriculum_subjects" json:"subjects,omitempty"`
	Certificates     []*Certificate `gorm:"many2many:cert_curriculums" json:"certificates,omitempty"`
	Category         *Category      `gorm:"foreignKey:CategoryId" json:"category,omitempty"`
	Students         []*Student     `json:"-" gorm:"many2many:student_curriculums"`
}
