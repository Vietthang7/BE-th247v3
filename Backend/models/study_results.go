package models

import (
	"github.com/google/uuid"
	"time"
)

type ExamResult struct {
	Model        `gorm:"embedded"` // created_at là ngày nộp bài
	Name         string            `json:"name,omitempty"`
	LessonData   *LessonData       `json:"lesson_data,omitempty" gorm:"foreignKey:LessonDataId"`
	Point        float64           `json:"point,omitempty" gorm:"type:float;default:0"`       // Điểm bài kiểm tra
	TotalPoint   float64           `json:"total_point,omitempty" gorm:"type:float;default:0"` // Tổng số điểm trong bài kiểm tra
	Result       bool              `json:"result" gorm:"default:false"`
	Status       int               `json:"status" gorm:"default:1"`       // Kết quả: true: Đạt, false: ko đạt
	IsGrade      bool              `json:"is_grade" gorm:"default:false"` // true : đã chấm , false : chưa chấm
	LessonDataId uuid.UUID         `json:"lesson_data_id"`
	TestId       uuid.UUID         `json:"test_id,omitempty"`
	ClassId      uuid.UUID         `json:"-"`
	StudentId    uuid.UUID         `json:"-"`
	UserId       *uuid.UUID        `json:"-"`
	CenterId     uuid.UUID         `json:"-"`
	UserAnswerId *uuid.UUID        `json:"user_answer_id,omitempty"`                // lưu giữ câu trả lời của người dùng
	Test         *TestService      `json:"exam,omitempty" gorm:"foreignKey:TestId"` // trở tới bảng bài kiểm tra
	GradeDate    *time.Time        `json:"grade_date,omitempty"`                    // lưu giữ ngày chấm điểm
	Class        *Class            `json:"class,omitempty" gorm:"foreignKey:ClassId"`
	Student      *Student          `json:"student,omitempty" gorm:"foreignKey:StudentId"`
	Teacher      *User             `json:"teacher,omitempty" gorm:"foreignKey:UserId"` //Giáo viên chấm điểm
	Center       *Center           `gorm:"foreignKey:CenterId" json:"-"`
}
type TestService struct {
	Model
	Name       string  `json:"name,omitempty"`
	Duration   int     `json:"duration"` // time
	Status     string  `json:"status"`
	IsSetScore bool    `json:"is_set_score"` // kiểm tra xem bài kiểm tra có được thiết lập điểm số hay không
	MaxAnswers int64   `json:"max_answers"`
	PassPoint  float64 `json:"pass_point" gorm:"type:float"` // điểm tối thiểu mà sinh viên cần đạt để pass
}
