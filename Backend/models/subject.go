package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Subject struct {
	Model         `gorm:"embedded"`
	Name          string         `gorm:"index;size:250" json:"name"` //s_c : subject -center
	Thumbnail     string         `json:"thumbnail,omitempty"`
	IsActive      *bool          `gorm:"default:true" json:"is_active,omitempty"`
	FeeType       uint8          `gorm:"default:2" json:"fee_type,omitempty"` //1 - free, 2 - payment
	OriginFee     uint64         `gorm:"default:0" json:"origin_fee,omitempty"`
	DiscountFee   uint64         `gorm:"default:0" json:"discount_fee,omitempty"`
	Description   string         `gorm:"type:text" json:"description,omitempty"`
	TotalLessons  uint64         `gorm:"default:null" json:"total_lessons,omitempty"`
	InputRequire  string         `gorm:"type:text" json:"input_require,omitempty"`
	OutputRequire string         `gorm:"type:text" json:"output_require,omitempty"`
	CertificateId *uuid.UUID     `gorm:"default:null" json:"certificate_id,omitempty"`
	CertIssuance  *bool          `gorm:"default:null" json:"cert_issuance"`
	Metadata      datatypes.JSON `json:"metadata"`
	Position      int            `json:"position,omitempty" gorm:"default:1"`
	CreatedBy     uuid.UUID      `json:"-"`
	CenterId      uuid.UUID      `json:"-"`
	CategoryId    uuid.UUID      `json:"-"`
	Code          string         `gorm:"default:null;index;size:20" json:"code,omitempty"`
	Center        Center         `gorm:"foreignKey:CenterId" json:"-"`
	Curriculums   []*Curriculum  `gorm:"many2many:curriculum_subjects" json:"curriculums,omitempty"`
	Lessons       []*Lesson      `json:"lessons,omitempty"`
	Category      *Category      `gorm:"foreignKey:CategoryId" json:"category,omitempty"`
	UserCreated   *User          `gorm:"foreignKey:CreatedBy" json:"created_by,omitempty"`
	Teachers      []*User        `gorm:"many2many:subject_teachers;joinForeignKey:SubjectID;" json:"teachers,omitempty"`

	Students           []*Student   `json:"-" gorm:"many2many:student_subjects"`
	Certificate        *Certificate `gorm:"foreignKey:CertificateId" json:"certificate,omitempty"`
	TeacherNames       string       `json:"teacher_names,omitempty" gorm:"-"` // Tên các Giảng viên
	Stars              float64      `json:"stars" gorm:"-"`                   // Điểm đánh giá
	TotalReview        int64        `json:"total_review" gorm:"-"`            // Số Đánh giá
	TotalUpcomingClass int64        `json:"total_upcoming_class" gorm:"-"`    // Số lớp sắp diễn ra
	TotalStudent       int64        `json:"total_student" gorm:"-"`           // Số học viên đăng ký
	LessonCount        int64        `json:"lesson_count" gorm:"-"`            // Số lượng bài học
}

type SubjectMoreInfo struct {
	Subject
	StudentPendings int `json:"student_pendings"`
	ClassTotal      int `json:"class_total"`
}
