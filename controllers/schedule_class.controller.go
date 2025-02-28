package controllers

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type NewScheduleClassInput struct {
	Id              *uuid.UUID            `json:"id"`
	ClassId         uuid.UUID             `json:"class_id"`
	Type            uint8                 `json:"type"`
	ScheduleDetails []ScheduleDetailInput `json:"schedule_details"`
	Weeks           []int                 `json:"weeks"`
	Dates           []time.Time           `json:"dates"`
}
type ScheduleDetailInput struct {
	Id            *uuid.UUID      `json:"id"`
	Name          string          `json:"name"`
	StartDate     *time.Time      `json:"start_date"`
	StartTime     datatypes.Time  `json:"start_time"`
	EndTime       datatypes.Time  `json:"end_time"`
	TeacherId     uuid.UUID       `json:"teacher_id"`
	AsistantId    *uuid.UUID      `json:"asistant_id"`
	ClassroomIds  []uuid.UUID     `json:"classroom_ids"`
	WorkSessionId uuid.UUID       `json:"session_id"`
	Childrens     []ScheduleChild `json:"childrens"`
}
type ScheduleChild struct {
	Id            *uuid.UUID     `json:"id"`
	StartTime     datatypes.Time `json:"start_time"`
	EndTime       datatypes.Time `json:"end_time"`
	TeacherId     uuid.UUID      `json:"teacher_id"`
	AsistantId    *uuid.UUID     `json:"asistant_id"`
	ClassroomIds  []uuid.UUID    `json:"classroom_ids"`
	WorkSessionId uuid.UUID      `json:"session_id"`
}

//func CreateScheduleClass(c *fiber.Ctx) error {
//	user, err := repo.GetTokenData(c)
//	if err != nil {
//		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
//	}
//	var (
//		input                 NewScheduleClassInput
//		scheduleClass         models.ScheduleClass
//		teacherIds            []uuid.UUID
//		asistantIds           []uuid.UUID
//		classroomIds          []uuid.UUID
//		workSessionIds        []uuid.UUID
//		totalLessonChild      int              // đếm tổng số buổi học con
//		listScheduleIdChanged []uuid.UUID      // lưu trữ ID của các lịch học đã thay đổi hoặc cần được cập nhật trong quá trình xử lý
//		isChangeType          bool             // ùng để kiểm tra xem loại lịch học (Type) có thay đổi so với lịch học hiện tại hay không.
//		listScheduleIdSync    []uuid.UUID      //Liên kết chính xác các bài học với lịch học mới sau khi tạo hoặc cập nhật.
//		listLessonsSync       []*models.Lesson // Đảm bảo các bài học được gắn đúng với lịch học tương ứng.
//		scheduleIds           []uuid.UUID      // : Ngăn chặn thay đổi lịch học khi đã có thông tin điểm danh, đảm bảo tính toàn vẹn dữ liệu.
//	)
//	teacherCalendarStart := make(map[uuid.UUID][]time.Time)
//	teacherCalendarEnd := make(map[uuid.UUID][]time.Time)
//	asistantCalendarStart := make(map[uuid.UUID][]time.Time)
//	asistantCalendarEnd := make(map[uuid.UUID][]time.Time)
//	classRoomStart := make(map[uuid.UUID][]time.Time)
//	classRoomEnd := make(map[uuid.UUID][]time.Time)
//	var scheduleTimes []time.Time
//	totalLessonChild = 0
//	if err := c.BodyParser(&input); err != nil {
//		return ResponseError(c, fiber.StatusBadRequest, "Invalid", err.Error())
//	}
//
//}
