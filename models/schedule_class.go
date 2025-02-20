package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type ScheduleClass struct {
	Model         `gorm:"embedded"`
	Name          string          `gorm:"size:50;default:null" json:"name,omitempty"`
	StartDate     *time.Time      `gorm:"default:null;index" json:"start_date,omitempty"`
	StartTime     *datatypes.Time `gorm:"default:null;index" json:"start_time,omitempty"`
	EndTime       *datatypes.Time `gorm:"default:null;index" json:"end_time,omitempty"`
	Type          uint8           `gorm:"default:null" json:"type,omitempty"`
	WorkSessionId *uuid.UUID      `gorm:"default:null" json:"work_session_id,omitempty"`
	ParentId      *uuid.UUID      `gorm:"default:null" json:"-"`
	TeacherId     *uuid.UUID      `gorm:"default:null" json:"-"`
	AsistantId    *uuid.UUID      `gorm:"default:null" json:"-"`
	Classrooms    []Classroom     `gorm:"many2many:schedule_classrooms" json:"classrooms,omitempty"`
	ClassId       uuid.UUID       `json:"-"`
	CreatedBy     uuid.UUID       `json:"-"`
	CenterId      uuid.UUID       `json:"-"`
	Index         int             `gorm:"default:NULL" json:"index,omitempty"`
	Metadata      datatypes.JSON  `json:"metadata,omitempty"`
	Teacher       *User           `gorm:"foreignKey:TeacherId" json:"teacher,omitempty"`
	Asistant      *User           `gorm:"foreignKey:AsistantId" json:"asistant,omitempty"`
	Creater       *User           `gorm:"foreignKey:CreatedBy" json:"created_by,omitempty"`
	Center        *Center         `gorm:"foreignKey:CenterId" json:"-"`
	Childrens     []ScheduleClass `gorm:"foreignKey:ParentId" json:"childrens,omitempty"`
	Class         *Class          `gorm:"foreignKey:ClassId" json:"class,omitempty"`
	WorkSession   *WorkSession    `gorm:"foreignKey:WorkSessionId" json:"work_session,omitempty"`
	//not migrate
	Attendancers *[]SessionAttendance `gorm:"foreignKey:ClassSessionId" json:"attendancers,omitempty"`
	Students     *[]StudentClasses    `gorm:"foreignKey:ClassId;references:ClassId" json:"students,omitempty"`
}
