package models

import (
	"github.com/google/uuid"
	"time"
)

// Kế hoạch tuyển sinh
type EnrollmentPlan struct {
	Model
	Name        string        `json:"name,omitempty"`
	StartDate   *time.Time    `json:"start_date,omitempty"`
	EndDate     *time.Time    `json:"end_date,omitempty"`
	BranchId    *uuid.UUID    `json:"branch_id"`
	Branch      *Branch       `json:"branch,omitempty" gorm:"foreignKey:BranchId"`
	TotalMem    int           `json:"total_members,omitempty"`
	CenterId    *uuid.UUID    `json:"-"`
	Center      *Center       `gorm:"foreignKey:CenterId" json:"-"`
	Subjects    []*Subject    `json:"subjects,omitempty" gorm:"many2many:enroll_subject"`
	Curriculums []*Curriculum `json:"curriculums,omitempty" gorm:"many2many:enroll_curriculum"`
	LearnType   int           `json:"learn_type,omitempty"` // phương thức học áp dụng 1: online, 2: offline, 3: cả 2
	Description string        `gorm:"type:text;default:null" json:"description,omitempty"`
	PaymentType int           `json:"payment_type,omitempty"` // hình thức thanh toán 1: trọn gói 1 lần, 2: theo từng lớp, 3: theo lần
}
