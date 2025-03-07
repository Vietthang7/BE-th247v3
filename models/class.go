package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"intern_247/consts"
	"time"
)

type Class struct {
	Model        `gorm:"embedded"`
	Name         string     `gorm:"size:404;index" json:"name"`
	Code         string     `gorm:"size:30;index" json:"code,omitempty"`
	StartAt      *time.Time `gorm:"index" json:"start_at,omitempty"`
	EndAt        *time.Time `gorm:"index" json:"end_at,omitempty"`
	Type         int        `gorm:"index" json:"type,omitempty"`
	Description  string     `gorm:"type:text" json:"description,omitempty"` //number and character
	BranchId     uuid.UUID  `json:"-"`
	ClassroomId  uuid.UUID  `json:"-"`
	Classroom    *Classroom `gorm:"foreignKey:ClassroomId" json:"classroom,omitempty"`
	PlanId       uuid.UUID  `gorm:"default:null" json:"plan_id"` //ke hoach tuyen sinh
	CurriculumId *uuid.UUID `gorm:"default:null" json:"curriculum_id,omitempty"`
	// CategoryId   uuid.UUID      `gorm:"default:null" json:"-"`
	SubjectId     uuid.UUID `json:"subject_id"`
	CreatedBy     uuid.UUID `json:"-"`
	CenterId      uuid.UUID `json:"-"`
	GroupUrl      string    `gorm:"default:null" json:"group_url,omitempty"`
	TotalLessons  uint64    `json:"total_lessons"`
	TotalStudent  int64     `json:"total_students"`
	LessonLearned int64     `json:"lesson_learned,omitempty" gorm:"-"`
	TotalChild    int       `gorm:"default:0" json:"-"`
	IsAdded       *bool     `json:"is_added,omitempty" gorm:"-"`
	// ClassroomId uuid.UUID      `json:"-"`
	CurrentLessonData uuid.UUID      `json:"current,omitempty" gorm:"-"`
	CuratorId         uuid.UUID      `gorm:"default:NULL" json:"-"`
	Metadata          datatypes.JSON `json:"metadata,omitempty"`
	Branch            *Branch        `gorm:"foreignKey:BranchId" json:"branch,omitempty"`
	Curriculum        *Curriculum    `gorm:"foreignKey:CurriculumId" json:"curriculum,omitempty"`
	// Category   *Category      `gorm:"foreignKey:CategoryId" json:"category,omitempty"`
	Subject *Subject `gorm:"foreignKey:SubjectId" json:"subject,omitempty"`
	// Classroom *Classroom `gorm:"foreignKey:ClassroomId" json:"classroom,omitempty"`
	Curator             *User                  `gorm:"foreignKey:CuratorId" json:"curator,omitempty"`
	Center              *Center                `gorm:"foreignKey:CenterId" json:"-"`
	Creater             *User                  `gorm:"foreignKey:CreatedBy" json:"created_by,omitempty"`
	Students            []*Student             `gorm:"many2many:student_classes" json:"students,omitempty"`
	ScheduleClass       []*ScheduleClass       `gorm:"foreignKey:ClassId" json:"schedule_class,omitempty"`
	StudentClasses      *StudentClasses        `gorm:"foreignKey:ClassId" json:"student_classes,omitempty"`
	StudentsClasses     []StudentClasses       `gorm:"foreignKey:ClassId" json:"students_classes,omitempty"`
	Lessons             []Lesson               `gorm:"foreignKey:ClassId" json:"lessons,omitempty"`
	Status              int                    `gorm:"default:1;index" json:"status"`
	Enrollment          *EnrollmentPlan        `json:"enrollment,omitempty" gorm:"foreignKey:PlanId"`
	Exams               *[]ExamResult          `gorm:"foreignKey:ClassId" json:"exams,omitempty"`
	SessionAttendancers []SessionAttendance    `gorm:"foreignKey:ClassId" json:"-"`
	StudentCertificates *[]StudentCertificates `gorm:"foreignKey:ClassId" json:"student_certificates,omitempty"`
}

type StudyProgress struct {
	StudentId    string     `json:"student_id"`
	ClassId      uuid.UUID  `json:"-"`
	LessonDataId uuid.UUID  `json:"-"`
	Progress     float64    `json:"progress" gorm:"type:float;default:0"`
	Completed    bool       `json:"completed" gorm:"default:false"`
	CompletedAt  *time.Time `json:"completed_at"`
}
type StudentClasses struct {
	ClassId    uuid.UUID  `json:"-"`
	StudentId  uuid.UUID  `json:"-"`
	Status     *int64     `json:"status" gorm:"default:null"`
	ReservedAt *time.Time `json:"reserved_at,omitempty" gorm:"default:null"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	Progress   int        `json:"progress" gorm:"default:0"`
	Result     int        `json:"result" gorm:"default:0"`
}
type ClassOverview struct {
	ComingSoon int64 `json:"coming_soon"`
	InProgress int64 `json:"in_progress"`
	Finished   int64 `json:"finished"`
	Canceled   int64 `json:"canceled"`
}

func (c *Class) AfterFind(*gorm.DB) error {
	if c.Status == consts.CLASS_CANCELED || c.Status == consts.CLASS_FINISHED {
		return nil
	}
	c.Status = consts.CLASS_COMING_SOON
	if c.StartAt != nil && time.Now().After(*c.StartAt) {
		c.Status = consts.CLASS_IN_PROGRESS
	}
	if c.EndAt != nil && time.Now().After(*c.EndAt) {
		c.Status = consts.CLASS_FINISHED
	}
	return nil
}
