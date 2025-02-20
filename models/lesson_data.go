package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type LessonData struct {
	Model       `gorm:"embedded"`
	Name        string         `gorm:"size:250" json:"name"`
	Type        int            `json:"type,omitempty"`
	LessonId    uuid.UUID      `json:"-"`
	Metadata    datatypes.JSON `json:"metadata,omitempty"`
	Lesson      *Lesson        `gorm:"foreignKey:LessonId" json:"lesson,omitempty"`
	Position    uint64         `json:"position,omitempty"`
	CreatedBy   uuid.UUID      `json:"-"`
	CenterId    uuid.UUID      `json:"-"`
	Center      *Center        `gorm:"foreignKey:CenterId" json:"-"`
	UserCreated *User          `gorm:"foreignKey:CreatedBy" json:"created_by,omitempty"`
	Progress    *StudyProgress `json:"progress,omitempty" gorm:"foreignKey:LessonDataId"`
	Exam        []*ExamResult  `json:"exam_result,omitempty" gorm:"foreignKey:LessonDataId"`
	DocumentId  *uuid.UUID     `json:"document_id,omitempty"`
}
