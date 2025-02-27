package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Lesson struct {
	Model         `gorm:"embedded"`
	Name          string         `gorm:"size:250" json:"name"`
	ParentId      *uuid.UUID     `gorm:"default:null" json:"-"`
	SubjectId     *uuid.UUID     `gorm:"default:null" json:"-"`
	ClassId       *uuid.UUID     `gorm:"default:null" json:"-"`
	CenterId      uuid.UUID      `json:"-"`
	CreatedBy     uuid.UUID      `json:"-"`
	FreeTrial     *bool          `gorm:"default:false" json:"free_trial"`
	Position      uint64         `gorm:"default:1" json:"position,omitempty"`
	Metadata      datatypes.JSON `json:"metadata,omitempty"`
	ScheduleId    uuid.UUID      `gorm:"default:null" json:"-"`
	Subject       *Subject       `gorm:"foreignKey:SubjectId" json:"subject,omitempty"`
	UserCreated   *User          `gorm:"foreignKey:CreatedBy" json:"created_by,omitempty"`
	IsLive        *bool          `gorm:"default:false" json:"is_live,omitempty"`
	Center        *Center        `gorm:"foreignKey:CenterId" json:"-"`
	Childrens     []*Lesson      `gorm:"foreignKey:ParentId" json:"childrens,omitempty"`
	LessonDatas   []*LessonData  `gorm:"foreignKey:LessonId" json:"lesson_datas,omitempty"`
	ScheduleClass *ScheduleClass `gorm:"foreignKey:ScheduleId" json:"schedule_class,omitempty"`
	Class         *Class         `gorm:"foreignKey:ClassId" json:"class,omitempty"`
}
