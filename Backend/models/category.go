package models

import "github.com/google/uuid"

type Category struct {
	Model          `gorm:"embedded"`
	Name           string       `gorm:"index:,class:FULLTEXT,size:100,not null" json:"name"`
	ParentId       *uuid.UUID   `gorm:"default:null" json:"-"`
	Thumbnail      string       `gorm:"default:null" json:"thumbnail,omitempty"`
	Description    string       `gorm:"type:text;default:null" json:"description,omitempty"`
	IsActive       *bool        `gorm:"default:true" json:"is_active,omitempty"`
	CreatedBy      uuid.UUID    `json:"-"`
	CenterId       uuid.UUID    `gorm:"not null" json:"-"`
	Center         *Center      `gorm:"foreignKey:CenterId" json:"center,omitempty"`
	Created        *User        `gorm:"foreignKey:CreatedBy" json:"created_by,omitempty"`
	ChildCategory  []Category   `gorm:"foreignKey:ParentId" json:"childrens,omitempty"`
	ParentCategory *Category    `gorm:"foreignKey:ParentId" json:"parent,omitempty"`
	Curriculums    []Curriculum `gorm:"foreignKey:CategoryId" json:"curriculums,omitempty"`
	Subjects       []Subject    `gorm:"foreignKey:CategoryId" json:"subjects,omitempty"`
}
