package consts

import "github.com/google/uuid"

type Pagination struct {
	CurrentPage  int     `json:"current_page"`
	TotalPages   float64 `json:"total_pages"`
	TotalResults int64   `json:"total"`
}
type Query struct {
	ID             uuid.UUID `json:"id"`
	Search         string    `json:"search"`
	Active         string    `query:"active"`
	Page           int       `json:"page"`
	Length         int       `json:"length" form:"length"`
	Order          string    `json:"order"`
	Sort           string    `json:"sort"`
	Relation       string    `json:"relation" query:"relation"` // id relation (child id or parent id ...)
	Children       string    `query:"children"`                 // get child element ??
	Curriculum     string    `query:"curriculum"`
	StudentId      string    `query:"studentId"`
	EnrollmentId   string    `query:"enroll_id"`
	Branch         string    `query:"branch"`
	Certificate    string    `query:"certificate"`
	Class          string    `query:"class"`
	StartAt        string    `query:"start_at"`
	EndAt          string    `query:"end_at"`
	Subject        string    `query:"subject"`
	Type           uint      `query:"type"`
	ScheduleLength int       `query:"schedule_length"`
	Classroom      string    `query:"classroom"`
	Status         int       `query:"status"`
	Teacher        string    `query:"teacher"`
	Result         int       `query:"result"`
	DisableId      string    `query:"disable_id"`
}
