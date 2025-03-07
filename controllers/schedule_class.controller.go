package controllers

import (
	"fmt"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"intern_247/utils"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/teambition/rrule-go"
	"gorm.io/datatypes"
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

func CreateScheduleClass(c *fiber.Ctx) error {
	user, err := repo.GetTokenData(c)
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}
	var (
		input            NewScheduleClassInput
		scheduleClass    models.ScheduleClass
		teacherIds       []uuid.UUID
		asistantIds      []uuid.UUID
		classroomIds     []uuid.UUID
		workSessionIds   []uuid.UUID
		totalLessonChild int // đếm tổng số buổi học con
		//listScheduleIdChanged []uuid.UUID      // lưu trữ ID của các lịch học đã thay đổi hoặc cần được cập nhật trong quá trình xử lý
		isChangeType bool // ùng để kiểm tra xem loại lịch học (Type) có thay đổi so với lịch học hiện tại hay không.
		//listScheduleIdSync    []uuid.UUID      //Liên kết chính xác các bài học với lịch học mới sau khi tạo hoặc cập nhật.
		//listLessonsSync       []*models.Lesson // Đảm bảo các bài học được gắn đúng với lịch học tương ứng.
		scheduleIds []uuid.UUID // : Ngăn chặn thay đổi lịch học khi đã có thông tin điểm danh, đảm bảo tính toàn vẹn dữ liệu.
	)
	//Tạo các bản đồ (maps) để lưu trữ lịch bắt đầu và kết thúc của giáo viên, trợ giảng, và phòng học, cùng với danh sách thời gian kết thúc lịch để tính ngày kết thúc lớp.
	//teacherCalendarStart := make(map[uuid.UUID][]time.Time)
	//teacherCalendarEnd := make(map[uuid.UUID][]time.Time)
	//asistantCalendarStart := make(map[uuid.UUID][]time.Time)
	//asistantCalendarEnd := make(map[uuid.UUID][]time.Time)
	//classRoomStart := make(map[uuid.UUID][]time.Time)
	//classRoomEnd := make(map[uuid.UUID][]time.Time)
	//var scheduleTimes []time.Time // Lưu các thời gian kết thúc (EndTime) của các lịch để tìm ngày kết thúc lớn nhất của lớp (class.EndAt).
	totalLessonChild = 0
	if err := c.BodyParser(&input); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", err.Error())
	}
	if input.Type != consts.SCHEDULE_CLASS_DAY_TYPE && input.Type != consts.SCHEDULE_CLASS_WEEK_TYPE {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "Type not support", consts.InvalidReqInput)
	}
	class, err := repo.GetClassAndSubjectByIdAndCenterId(input.ClassId, user.CenterId)
	fmt.Println(user.CenterId)
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "Invalid 4a", consts.InvalidReqInput)
	}
	if class.Status == consts.CLASS_CANCELED || class.Status == consts.CLASS_FINISHED {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid serve", consts.InvalidReqInput)
	}

	//get holiday
	// Lấy danh sách ngày lễ (holidays) dựa trên ngày bắt đầu lớp (class.StartAt), chi nhánh (class.BranchId), và trung tâm (user.CenterId).
	holidays, err := repo.GetHolidayByDateBranchIdAndCenterId(class.StartAt, class.BranchId, user.CenterId)
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "get holiday failed", consts.InvalidReqInput)
	}
	// Lấy danh sách lịch học hiện có của lớp và kiểm tra xem có điểm danh (attendance) nào chưa, để đảm bảo không thể thay đổi lịch đã có điểm danh.
	schedules, err := repo.GetListScheduleByClassId(class.ID, consts.Query{}, user)
	for i := range schedules {
		scheduleIds = append(scheduleIds, schedules[i].ID)
	}
	if len(scheduleIds) > 0 {
		attendancers, err := repo.GetAttendanceByScheduleIdsAndClassId(scheduleIds, class.ID)
		if err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, "Invalid serve", consts.InvalidReqInput)
		}
		if len(attendancers) > 0 {
			return ResponseError(c, fiber.StatusInternalServerError, "Invalid serve", consts.ERROR_SCHEDULE_CLASS_ATTENDANCED)
		}
	}
	//validate and get uuid list classroom, worksession, teacher
	for i := range input.ScheduleDetails {
		if input.ScheduleDetails[i].EndTime < input.ScheduleDetails[i].StartTime {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_START_TIME_MUST_SMALLER_THAN_END_TIME)
		}
		teacherIds = append(teacherIds, input.ScheduleDetails[i].TeacherId)
		if input.ScheduleDetails[i].AsistantId != nil {
			asistantIds = append(asistantIds, *input.ScheduleDetails[i].AsistantId)
		}
		classroomLen := len(input.ScheduleDetails[i].ClassroomIds)
		if classroomLen > 2 {
			return ResponseError(c, fiber.StatusBadRequest, "số lượng lớp không hỗ trợ", consts.InvalidReqInput)
		}
		if classroomLen > 1 && class.Type != consts.CLASS_TYPE_HYBRID {
			return ResponseError(c, fiber.StatusBadRequest, "type not hybrid", consts.InvalidReqInput)
		}
		classroomIds = append(classroomIds, input.ScheduleDetails[i].ClassroomIds...)
		workSessionIds = append(workSessionIds, input.ScheduleDetails[i].WorkSessionId)
		for j := range input.ScheduleDetails[i].Childrens {
			if input.ScheduleDetails[i].Childrens[j].EndTime < input.ScheduleDetails[i].Childrens[j].StartTime {
				return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.ERROR_START_TIME_MUST_SMALLER_THAN_END_TIME)
			}
			classroomLen = len(input.ScheduleDetails[i].Childrens[j].ClassroomIds)
			if classroomLen > 2 {
				return ResponseError(c, fiber.StatusBadRequest, "len not support", consts.InvalidReqInput)
			}
			if classroomLen > 1 && class.Type != consts.CLASS_TYPE_HYBRID {
				return ResponseError(c, fiber.StatusBadRequest, "child type not hybrid", consts.InvalidReqInput)
			}
			teacherIds = append(teacherIds, input.ScheduleDetails[i].Childrens[j].TeacherId)
			if input.ScheduleDetails[i].Childrens[j].AsistantId != nil {
				asistantIds = append(asistantIds, *input.ScheduleDetails[i].Childrens[j].AsistantId)
			}
			classroomIds = append(classroomIds, input.ScheduleDetails[i].ClassroomIds...)
			workSessionIds = append(workSessionIds, input.ScheduleDetails[i].WorkSessionId)
			totalLessonChild = totalLessonChild + 1
		}
	}
	// Loại bỏ trùng lặp ID
	teacherIds = lo.Uniq(teacherIds)
	asistantIds = lo.Uniq(asistantIds)
	classroomIds = lo.Uniq(classroomIds)
	workSessionIds = lo.Uniq(workSessionIds)
	// Kiểm tra thực thể
	teachers, err := repo.GetTeachersByIdsAndCenterId(teacherIds, user.CenterId)
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "Invalid teacher", consts.InvalidReqInput)
	}
	if len(teachers) != len(teacherIds) {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "Invalid teacher1", consts.InvalidReqInput)
	}
	asisants, err := repo.GetAsistantsByIdsAndCenterId(asistantIds, user.CenterId)
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, "Invalid asistant find", consts.InvalidReqInput)
	}
	if len(asisants) != len(asistantIds) {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "Invalid asistant find current", consts.InvalidReqInput)
	}
	classrooms, err := repo.GetClassroomsByIdsAndBranchCenterId(classroomIds, class.BranchId, user.CenterId)
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "Invalid classroom", consts.InvalidReqInput)
	}
	if len(classrooms) != len(classroomIds) {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "Invalid classroom", consts.InvalidReqInput)
	}
	// Kiểm tra lịch học cũ
	scheduleClassOld, err := repo.GetSingleScheduleClassByClassId(input.ClassId, user.CenterId)
	if err == nil {
		scheduleClass = scheduleClassOld
		if scheduleClass.Type != input.Type {
			isChangeType = true
		}
	}
	// Xử lí lịch theo tuần
	if input.Type == consts.SCHEDULE_CLASS_WEEK_TYPE {
		Weeks := lo.Uniq(input.Weeks)
		input.Weeks = Weeks
		dates := handleGetDateNotInHoliday(Weeks, int(class.TotalLessons), *class.StartAt, holidays)
		checkFirstDay := dates[0].Weekday()
		fmt.Println("day la checkFirstDay")
		fmt.Println(checkFirstDay)
		//firstDayIndex := utils.Index(Weeks, int(checkFirstDay))
	}
	if isChangeType {
		fmt.Println("okhehe")
	}
	return ResponseSuccess(c, fiber.StatusCreated, "Success", nil)
}
func handleGetDateNotInHoliday(Weeks []int, totalLessons int, startDate time.Time, holidays []models.Holiday) []time.Time {
	var weekRules []rrule.Weekday
	sort.Ints(Weeks) // Sắp xếp mảng Weeks tăng dần để đảm bảo thứ tự ngày hợp lý.
	for _, week := range Weeks {
		var day rrule.Weekday
		switch week {
		case 1:
			day = rrule.MO // Monday 1
		case 2:
			day = rrule.TU // Tuesday 2
		case 3:
			day = rrule.WE // Wednesday 3
		case 4:
			day = rrule.TH // Thursday 4
		case 5:
			day = rrule.FR // Friday 5
		case 6:
			day = rrule.SA // Saturday 6
		case 0:
			day = rrule.SU // Sunday 0
		default:
			return []time.Time{} // Trả về mảng rỗng nếu tuần không hợp lệ
		}
		weekRules = append(weekRules, day)
	}
	// Tạo rule lặp lại hàng tuần
	rule, err := rrule.NewRRule(rrule.ROption{
		Freq:      rrule.WEEKLY,
		Count:     totalLessons,
		Byweekday: weekRules,
		Dtstart:   startDate,
		Interval:  1, // Lặp lại hàng tuần
	})
	if err != nil {
		return []time.Time{}
	}
	// Lấy tất cả các ngày theo rule
	dates := rule.All() // Trả về mảng []time.Time chứa tất cả các ngày theo quy tắc lặp.
	// Hàm đệ quy để loại bỏ ngày trùng và bổ sung ngày mới nếu cần
	return getNonHolidayDates(dates, totalLessons, holidays, weekRules)
}

// Hàm đệ quy để kiểm tra và thay thế ngày trùng với ngày nghỉ
func getNonHolidayDates(dates []time.Time, totalLessons int, holidays []models.Holiday, weekRules []rrule.Weekday) []time.Time {
	// Kiểm tra xem có ngày nào trùng với ngày nghỉ không
	for i := range dates {
		for _, holiday := range holidays {
			if utils.IsDateInRange(dates[i], time.Time(holiday.StartDay), time.Time(holiday.EndDay)) {
				// Loại bỏ ngày trùng
				dates = append(dates[:i], dates[i+1:]...)
				// Tạo rule mới để thêm ngày thay thế
				missingLessons := totalLessons - len(dates)
				if missingLessons > 0 {
					// Tạo rule mới để lấy thêm ngày
					newRule, err := rrule.NewRRule(rrule.ROption{
						Freq:      rrule.WEEKLY,
						Count:     missingLessons,
						Byweekday: weekRules,
						Dtstart:   dates[len(dates)-1].AddDate(0, 0, 1), // Bắt đầu sau ngày cuối cùng trong dates
						Interval:  1,
					})
					if err != nil {
						return dates
					}
					// Lấy thêm các ngày mới
					newDates := newRule.All()
					// Gọi đệ quy để kiểm tra các ngày mới và thay thế tiếp nếu cần
					dates = append(dates, newDates...)
					return getNonHolidayDates(dates, totalLessons, holidays, weekRules)
				}
			}
		}
	}
	// Nếu số lượng ngày trong dates đạt yêu cầu thì trả về kết quả
	if len(dates) == totalLessons {
		return dates
	}
	return dates
}
