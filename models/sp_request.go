package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SupportRequest struct {
	Model             `gorm:"embedded"`
	Title             string         `json:"title" gorm:"not null;size:250"`
	Content           string         `json:"content,omitempty" gorm:"not null;type:longtext"`
	Type              int            `json:"type" gorm:"default:1"`      // Loại yêu cầu; internal/consts/config.go
	LeaveFromDate     *time.Time     `json:"leave_from_date,omitempty"`  // Ngày bắt đầu nghỉ
	LeaveUntilDate    *time.Time     `json:"leave_until_date,omitempty"` // Ngày kết thúc nghỉ
	MakeUpClass       *bool          `json:"make_up_class,omitempty"`    // Yêu cầu học bù
	SubjectIds        []uuid.UUID    `json:"subject_ids,omitempty" gorm:"-"`
	Subjects          []*Subject     `json:"subjects,omitempty" gorm:"many2many:subject_sp_requests"`
	Agree             *bool          `json:"agree"`                  // Kết quả xử lý
	ResolveDate       *time.Time     `json:"resolve_date,omitempty"` // Ngày xử lý
	Response          string         `json:"response,omitempty"`     // Nội dung phản hồi
	Metadata          datatypes.JSON `json:"metadata,omitempty"`
	CreatedBy         uuid.UUID      `json:"-"`
	Creator           *Student       `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	RespondedBy       *uuid.UUID     `json:"-"`
	Responder         *User          `json:"responder,omitempty" gorm:"foreignKey:RespondedBy"` // Người xử lý
	CenterID          uuid.UUID      `json:"-" gorm:"not null"`
	CanUD             bool           `json:"can_ud" gorm:"-"`
	RegisteredSubject int64          `json:"registered_subject,omitempty" gorm:"-"` // Số môn đã đăng ký
	FinishedSubject   int64          `json:"finished_subject,omitempty" gorm:"-"`   // Số môn hoàn thành
}

func (u *SupportRequest) AfterFind(*gorm.DB) error {
	if u.Agree == nil {
		u.CanUD = true
	} else {
		u.CanUD = false
	}
	return nil
}

type SubjectSpRequest struct {
	SupportRequestId uuid.UUID `json:"support_request_id"`
	SubjectId        uuid.UUID `json:"subject_id"`
}
