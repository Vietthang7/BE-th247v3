package consts

import (
	"github.com/google/uuid"
	"intern_247/utils"
	"math"
	"strconv"
)

type Pagination struct {
	CurrentPage  int     `json:"current_page"`
	TotalPages   float64 `json:"total_pages"`
	TotalResults int64   `json:"total"`
}

func (p *Pagination) GetTotalPages(len int) float64 {
	return math.Ceil(float64(p.TotalResults) / float64(len))
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

func (q *Query) GetActive() *bool {
	if q.Active != "" {
		if status, err := strconv.ParseBool(q.Active); err == nil {
			return &status
		}
	}
	return nil
}
func (q *Query) GetOffset() int {
	return (q.GetPage() - 1) * q.GetPageSize()
}
func (q *Query) GetPage() int {
	if q.Page < 1 {
		q.Page = 1
	}
	return q.Page
}
func (q *Query) GetPageSize() int {
	if q.Length > 200 {
		q.Length = 200
	}
	if q.Length < 1 {
		q.Length = 12
	}
	return q.Length
}
func (q *Query) GetField(Orders []string, d string) string {
	if !utils.Contains(Orders, q.Order) {
		q.Order = d
	}
	return q.Order
}
func (q *Query) GetSort() string {
	if q.Sort != "asc" {
		q.Sort = "desc"
	}
	return q.Sort
}
