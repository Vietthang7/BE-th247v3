package consts

import (
	"github.com/google/uuid"
	"time"
)

const (
	FREE_SUBJECT                          = 1
	PAYMENT_SUBJECT                       = 2
	YOUTUBE_LINK_TYPE                     = 1 //youtube
	S3_LINK_TYPE                          = 2 // link s3
	PARAGRAPH_TYPE                        = 3 //doan van
	TEST_TYPE                             = 4 // bai kiem tra
	DOCUMENT_TYPE                         = 5 //tai lieu
	SUBJECT_CODE_PREFIX                   = "MH"
	LESSON_DATA_POINT_FOLLOW_TEST         = 1
	LESSON_DATA_POINT_NOT_CHECK           = 2
	LESSON_DATA_DONE_TYPE_ANSWER_QUESTION = 1 // trả lời câu hỏi
	LESSON_DATA_DONE_TYPE_WATCHED_CONTENT = 2 // xem nội dung
)

var LESSON_DATAS_TYPE = []int{YOUTUBE_LINK_TYPE, S3_LINK_TYPE, PARAGRAPH_TYPE, TEST_TYPE, DOCUMENT_TYPE}

type LessonDataMetadata struct {
	Url      string `json:"url,omitempty"`
	FileName string `json:"fileName,omitempty"`
	Content  string `json:"content,omitempty"`
	//test
	TestId       *uuid.UUID `json:"test_id,omitempty"`
	TestName     string     `json:"test_name,omitempty"`
	Duration     int        `json:"duration,omitempty"`
	MaxAnswers   int64      `json:"max_answers,omitempty"`
	PointTest    uint       `json:"point_test,omitempty"`
	ExpiredAt    *time.Time `json:"expired_at,omitempty"`
	AllowExpired bool       `json:"allow_expired,omitempty"`
	FileSize     uint64     `json:"file_size,omitempty"`
	ContentType  string     `json:"content_type,omitempty"` // Loại nội dung file (VD: application/pdf, image/png).
	// class lesson config
	DoneType  int    `json:"done_type,omitempty"`  // Kiểu điều kiện để đánh dấu bài học đã hoàn thành (VD: 0 = Xem xong, 1 = Hoàn thành bài kiểm tra, v.v.).
	DoneValue string `json:"done_value,omitempty"` // Giá trị liên quan đến điều kiện hoàn thành (VD: số điểm tối thiểu cần đạt).
}
