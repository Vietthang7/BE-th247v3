package controllers

import (
	"encoding/json"
	"fmt"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"intern_247/utils"
	"sort"
	"strconv"
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
		input                 NewScheduleClassInput
		scheduleClass         models.ScheduleClass
		teacherIds            []uuid.UUID
		asistantIds           []uuid.UUID
		classroomIds          []uuid.UUID
		workSessionIds        []uuid.UUID
		totalLessonChild      int              // đếm tổng số buổi học con
		listScheduleIdChanged []uuid.UUID      // lưu trữ ID của các lịch học đã thay đổi hoặc cần được cập nhật trong quá trình xử lý
		isChangeType          bool             // ùng để kiểm tra xem loại lịch học (Type) có thay đổi so với lịch học hiện tại hay không.
		listScheduleIdSync    []uuid.UUID      //Liên kết chính xác các bài học với lịch học mới sau khi tạo hoặc cập nhật.
		listLessonsSync       []*models.Lesson // Đảm bảo các bài học được gắn đúng với lịch học tương ứng.
		scheduleIds           []uuid.UUID      // : Ngăn chặn thay đổi lịch học khi đã có thông tin điểm danh, đảm bảo tính toàn vẹn dữ liệu.
	)
	//Tạo các bản đồ (maps) để lưu trữ lịch bắt đầu và kết thúc của giáo viên, trợ giảng, và phòng học, cùng với danh sách thời gian kết thúc lịch để tính ngày kết thúc lớp.
	teacherCalendarStart := make(map[uuid.UUID][]time.Time)
	teacherCalendarEnd := make(map[uuid.UUID][]time.Time)
	asistantCalendarStart := make(map[uuid.UUID][]time.Time)
	asistantCalendarEnd := make(map[uuid.UUID][]time.Time)
	classRoomStart := make(map[uuid.UUID][]time.Time)
	classRoomEnd := make(map[uuid.UUID][]time.Time)
	var scheduleTimes []time.Time // Lưu các thời gian kết thúc (EndTime) của các lịch để tìm ngày kết thúc lớn nhất của lớp (class.EndAt).
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
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "Lớp học không hợp lệ", consts.InvalidReqInput)
	}
	if class.Status == consts.CLASS_CANCELED || class.Status == consts.CLASS_FINISHED {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid serve", consts.InvalidReqInput)
	}

	//get holiday
	// Lấy danh sách ngày lễ (holidays) dựa trên ngày bắt đầu lớp (class.StartAt), chi nhánh (class.BranchId), và trung tâm (user.CenterId).
	holidays, err := repo.GetHolidayByDateBranchIdAndCenterId(class.StartAt, *class.BranchId, user.CenterId)
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
	classrooms, err := repo.GetClassroomsByIdsAndBranchCenterId(classroomIds, *class.BranchId, user.CenterId)
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
		firstDayIndex := utils.Index(Weeks, int(checkFirstDay))
		if firstDayIndex != 0 && len(input.ScheduleDetails) == len(Weeks) {
			input.ScheduleDetails = append(input.ScheduleDetails[firstDayIndex:], input.ScheduleDetails[:firstDayIndex]...)
		}
		lessonCount := 0
		lessonIndex := 0
		//lịch trình trùng lặp nếu loại là tuần
		fmt.Println(class.TotalLessons)
		var scheduleDetailDuplicate []ScheduleDetailInput
		for lessonCount <= int(class.TotalLessons) {
			if lessonCount == int(class.TotalLessons) {
				break
			}
			for i := range input.ScheduleDetails {
				if lessonCount == int(class.TotalLessons) {
					break
				}
				schedule := input.ScheduleDetails[i]
				schedule.StartDate = &dates[lessonCount]
				lessonIndex++
				lessonCount++
				if lessonCount == int(class.TotalLessons) {
					schedule.Childrens = []ScheduleChild{}
					scheduleDetailDuplicate = append(scheduleDetailDuplicate, schedule)
					break
				}
				for j := range input.ScheduleDetails[i].Childrens {
					lessonCount++
					if lessonCount == int(class.TotalLessons) {
						schedule.Childrens = input.ScheduleDetails[j].Childrens[0 : j+1]
						scheduleDetailDuplicate = append(scheduleDetailDuplicate, schedule)
						break
					}
				}
				scheduleDetailDuplicate = append(scheduleDetailDuplicate, schedule)
			}
		}
		input.ScheduleDetails = scheduleDetailDuplicate
	}
	//kiểm tra giáo viên lịch trùng lặp
	for i := range input.ScheduleDetails {
		if input.ScheduleDetails[i].StartDate == nil {
			return ResponseError(c, fiber.StatusBadRequest, "Invalid sd", consts.InvalidReqInput)
		}
		if input.ScheduleDetails[i].AsistantId != nil {
			asistantCalendarStart[*input.ScheduleDetails[i].AsistantId] = append(asistantCalendarStart[*input.ScheduleDetails[i].AsistantId], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].StartTime))
			asistantCalendarEnd[*input.ScheduleDetails[i].AsistantId] = append(asistantCalendarEnd[*input.ScheduleDetails[i].AsistantId], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].EndTime))
		}
		// teacher
		teacherCalendarStart[input.ScheduleDetails[i].TeacherId] = append(teacherCalendarStart[input.ScheduleDetails[i].TeacherId], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].StartTime))
		teacherCalendarEnd[input.ScheduleDetails[i].TeacherId] = append(teacherCalendarEnd[input.ScheduleDetails[i].TeacherId], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].EndTime))
		//classrooms
		for c := range input.ScheduleDetails[i].ClassroomIds {
			classRoomStart[input.ScheduleDetails[i].ClassroomIds[c]] = append(classRoomStart[input.ScheduleDetails[i].ClassroomIds[c]], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].StartTime))
			classRoomEnd[input.ScheduleDetails[i].ClassroomIds[c]] = append(classRoomEnd[input.ScheduleDetails[i].ClassroomIds[c]], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].EndTime))
		}
		for j := range input.ScheduleDetails[i].Childrens {
			if input.ScheduleDetails[i].Childrens[j].AsistantId != nil {
				asistantCalendarStart[*input.ScheduleDetails[i].Childrens[j].AsistantId] = append(asistantCalendarStart[*input.ScheduleDetails[i].Childrens[j].AsistantId], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].Childrens[j].StartTime))
				asistantCalendarEnd[*input.ScheduleDetails[i].Childrens[j].AsistantId] = append(asistantCalendarEnd[*input.ScheduleDetails[i].Childrens[j].AsistantId], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].Childrens[j].EndTime))
			}
			teacherCalendarStart[input.ScheduleDetails[i].Childrens[j].TeacherId] = append(teacherCalendarStart[input.ScheduleDetails[i].Childrens[j].TeacherId], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].Childrens[j].StartTime))
			teacherCalendarEnd[input.ScheduleDetails[i].Childrens[j].TeacherId] = append(teacherCalendarEnd[input.ScheduleDetails[i].Childrens[j].TeacherId], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].Childrens[j].EndTime))
			for c := range input.ScheduleDetails[i].Childrens[j].ClassroomIds {
				classRoomStart[input.ScheduleDetails[i].Childrens[j].ClassroomIds[c]] = append(classRoomStart[input.ScheduleDetails[i].Childrens[j].ClassroomIds[c]], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].Childrens[j].StartTime))
				classRoomEnd[input.ScheduleDetails[i].Childrens[j].ClassroomIds[c]] = append(classRoomEnd[input.ScheduleDetails[i].Childrens[j].ClassroomIds[c]], *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].Childrens[j].EndTime))
			}
		}
	}
	for i := range teachers {
		if _, ok := teacherCalendarStart[teachers[i].ID]; !ok {
			continue
		}
		if len(teacherCalendarStart[teachers[i].ID]) > 0 {
			schedules, err := repo.GetScheduleClassByTeacherIdAndLimitDate(input.ClassId, teachers[i].ID, teacherCalendarStart[teachers[i].ID])
			if err != nil {
				return ResponseError(c, fiber.StatusInternalServerError, "Error", consts.ERROR_INTERNAL_SERVER_ERROR)
			}
			// Duyệt qua lịch học mới của giáo viên
			for j := range teacherCalendarStart[teachers[i].ID] {
				for k := range schedules {
					startAt := utils.MixedDateAndTime(schedules[k].StartDate, schedules[k].StartTime)
					endAt := utils.MixedDateAndTime(schedules[k].StartDate, schedules[k].EndTime)
					if utils.IsTimeRangeOverlap(teacherCalendarStart[teachers[i].ID][j], teacherCalendarEnd[teachers[i].ID][j], *startAt, *endAt) {
						return ResponseError(c, fiber.StatusBadRequest, fiber.Map{"date": startAt.Format("2006-01-02 15:04:05"), "teacher_id": teachers[i].ID}, consts.ERROR_SCHEDULE_CLASS_TEACHER_DUPLICATE_CALENDAR)
					}
				}
			}

		}
	}
	// check duplicate asistant teaching
	for i := range asisants {
		if _, ok := asistantCalendarStart[asisants[i].ID]; !ok {
			continue
		}
		if len(asistantCalendarStart[asisants[i].ID]) > 0 {
			schedules, err := repo.GetScheduleClassByAsistantIdAndLimitDate(input.ClassId, asisants[i].ID, asistantCalendarStart[asisants[i].ID])
			if err != nil {
				return ResponseError(c, fiber.StatusInternalServerError, "Error", consts.ERROR_INTERNAL_SERVER_ERROR)
			}
			for j := range asistantCalendarStart[asisants[i].ID] {
				for k := range schedules {
					startAt := utils.MixedDateAndTime(schedules[k].StartDate, schedules[k].StartTime)
					endAt := utils.MixedDateAndTime(schedules[k].StartDate, schedules[k].EndTime)
					if utils.IsTimeRangeOverlap(asistantCalendarStart[asisants[i].ID][j], asistantCalendarEnd[asisants[i].ID][j], *startAt, *endAt) {
						return ResponseError(c, fiber.StatusBadRequest, fiber.Map{"date": startAt.Format("2006-01-02 15:04:05"), "asistant_id": asisants[i].ID}, consts.ERROR_SCHEDULE_CLASS_ASISTANT_DUPLICATE_CALENDAR)
					}
				}
			}
		}
	}
	//check duplicate classroom
	for i := range classrooms {
		if _, ok := classRoomStart[classrooms[i].ID]; !ok {
			continue
		}
		if len(classRoomStart[classrooms[i].ID]) > 0 {
			schedules, err := repo.GetScheduleClassByClassroomIdAndLimitDate(input.ClassId, classrooms[i].ID, classRoomStart[classrooms[i].ID])
			if err != nil {
				return ResponseError(c, fiber.StatusInternalServerError, "Error", consts.ERROR_INTERNAL_SERVER_ERROR)
			}
			for j := range classRoomStart[classrooms[i].ID] {
				for k := range schedules {
					startAt := utils.MixedDateAndTime(schedules[k].StartDate, schedules[k].StartTime)
					endAt := utils.MixedDateAndTime(schedules[k].StartDate, schedules[k].EndTime)
					if utils.IsTimeRangeOverlap(classRoomStart[classrooms[i].ID][j], classRoomEnd[classrooms[i].ID][j], *startAt, *endAt) {
						return ResponseError(c, fiber.StatusBadRequest, fiber.Map{"date": startAt.Format("2006-01-02 15:04:05"), "classroom_id": classrooms[i].ID}, consts.ERROR_SCHEDULE_CLASS_CLASSROOM_DUPLICATE_CALENDAR)
					}
				}
			}
		}
	}
	//check if class is hybrid -> classroom
	if class.Type == consts.CLASS_TYPE_HYBRID {
		var onlineRooms []uuid.UUID
		var offlineRooms []uuid.UUID
		for i := range classrooms {
			if classrooms[i].IsOnline != nil && *classrooms[i].IsOnline {
				onlineRooms = append(onlineRooms, classrooms[i].ID)
			}
			if classrooms[i].IsOnline != nil && !*classrooms[i].IsOnline {
				offlineRooms = append(offlineRooms, classrooms[i].ID)
			}
		}
		for i := range input.ScheduleDetails {
			countOnline := 0
			countOffline := 0
			for j := range input.ScheduleDetails[i].ClassroomIds {
				if utils.Contains(onlineRooms, input.ScheduleDetails[i].ClassroomIds[j]) {
					countOnline++
				}
				if utils.Contains(offlineRooms, input.ScheduleDetails[i].ClassroomIds[j]) {
					countOffline++
				}
			}
			if countOnline > 1 || countOffline > 1 {
				return ResponseError(c, fiber.StatusBadRequest, "Online + Offline", consts.InvalidReqInput)
			}
			for j := range input.ScheduleDetails[i].Childrens {
				countOnline = 0
				countOffline = 0
				for k := range input.ScheduleDetails[i].Childrens[j].ClassroomIds {
					if utils.Contains(onlineRooms, input.ScheduleDetails[i].Childrens[j].ClassroomIds[k]) {
						countOnline++
					}
					if utils.Contains(offlineRooms, input.ScheduleDetails[i].Childrens[j].ClassroomIds[k]) {
						countOffline++
					}
				}
				if countOnline > 1 || countOffline > 1 {
					return ResponseError(c, fiber.StatusBadRequest, "Online + Offline child", consts.InvalidReqInput)
				}
			}
		}
	}
	// kiểm tra tổng số bài học trong môn học
	if (len(input.ScheduleDetails)+totalLessonChild) != int(class.TotalLessons) && input.Type == consts.SCHEDULE_CLASS_DAY_TYPE {
		return ResponseError(c, fiber.StatusBadRequest, "Total lesson invalid", consts.ERROR_TOTAL_LESSONS_NOT_EQUAL)
	}
	workSessions, err := repo.GetActiveWorkSessionByIdsAndBranchCenter(workSessionIds, *class.BranchId, user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "ws", consts.InvalidReqInput)
	}
	if len(workSessionIds) != len(workSessions) {
		return ResponseError(c, fiber.StatusBadRequest, "ws1", consts.InvalidReqInput)
	}
	//check time child record is valid
	for i := range input.ScheduleDetails {
		for j := range workSessions {
			if workSessions[j].ID == input.ScheduleDetails[i].WorkSessionId {
				if time.Duration(workSessions[j].StartTime.Hour()*3600+workSessions[j].StartTime.Minute()*60+workSessions[j].StartTime.Second())*time.Second > time.Duration(input.ScheduleDetails[i].StartTime) {
					return ResponseError(c, fiber.StatusBadRequest, "lịch trình không thể bắt đầu trước phiên làm việc", consts.ERROR_SC_START_TIME_INVALID)
				}
			}
		}
		if input.ScheduleDetails[i].StartDate != nil && class.StartAt != nil && class.StartAt.After(*input.ScheduleDetails[i].StartDate) {
			return ResponseError(c, fiber.StatusBadRequest, "Lịch trình không thể bắt đầu trước ngày khai giảng của lớp học.", consts.ERROR_SCHEDULE_DATE_IS_BIGGER_THAN_START_AT)
		}
		if input.Type == consts.SCHEDULE_CLASS_DAY_TYPE {
			for idx := range holidays {
				if input.ScheduleDetails[i].StartDate != nil {
					if utils.IsDateInRange(*input.ScheduleDetails[i].StartDate, time.Time(holidays[idx].StartDay), time.Time(holidays[idx].EndDay)) {
						return ResponseError(c, fiber.StatusBadRequest, utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].StartTime).Format("2006-01-02"), consts.ERROR_SCHEDULE_CLASS_CONFLICT_HOLIDAY)
					}
				}
			}
		}
		for k := range input.ScheduleDetails[i].Childrens {
			for j := range workSessions {
				if workSessions[j].ID == input.ScheduleDetails[i].Childrens[k].WorkSessionId {
					if time.Duration(workSessions[j].StartTime.Hour()*3600+workSessions[j].StartTime.Minute()*60+workSessions[j].StartTime.Second())*time.Second > time.Duration(input.ScheduleDetails[i].Childrens[k].StartTime) {
						return ResponseError(c, fiber.StatusBadRequest, "lịch trình không thể bắt đầu trước phiên làm việc", consts.ERROR_SC_START_TIME_INVALID)
					}
					// if time.Duration(workSessions[j].EndTime.In(location).Hour()*3600+workSessions[j].EndTime.In(location).Minute()*60+workSessions[j].EndTime.In(location).Second())*time.Second < time.Duration(input.ScheduleDetails[i].Childrens[k].EndTime) {
					// 	return ResponseError(c, fiber.StatusBadRequest, "Invalid 1", consts.ERROR_SC_END_TIME_INVALID)
					// }
				}
			}
		}
	}
	var calendar []byte
	if input.Type == consts.SCHEDULE_CLASS_WEEK_TYPE {
		calendar, err = json.Marshal(input.Weeks)
		if err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, "invalid calendar", consts.CreateFailed)
		}
	} else {
		calendar, err = json.Marshal(input.Dates)
		if err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, "invalid calendar date", consts.CreateFailed)
		}
	}
	if scheduleClass.ID != uuid.Nil {
		listScheduleIdChanged = append(listScheduleIdChanged, scheduleClass.ID)
	} else {
		scheduleClass.ID = uuid.New()
		scheduleClass.CenterId = user.CenterId
		scheduleClass.CreatedBy = user.ID
		scheduleClass.ClassId = class.ID
	}
	scheduleClass.Type = input.Type
	scheduleClass.Metadata = calendar
	newScheduleParent := scheduleClass
	scheduleIndex := 1
	if len(input.ScheduleDetails) > 0 {
		for i := range input.ScheduleDetails {
			var childSchedules []models.ScheduleClass
			var schedule models.ScheduleClass
			schedule.ID = uuid.New()
			if input.ScheduleDetails[i].Id != nil && !utils.Contains(listScheduleIdChanged, *input.ScheduleDetails[i].Id) {
				schedule.ID = *input.ScheduleDetails[i].Id
				listScheduleIdChanged = append(listScheduleIdChanged, schedule.ID)
			}
			if input.ScheduleDetails[i].StartDate == nil {
				return ResponseError(c, fiber.StatusInternalServerError, "Time not null", consts.CreateFailed)
			}
			schedule.StartDate = input.ScheduleDetails[i].StartDate
			schedule.StartTime = &input.ScheduleDetails[i].StartTime
			schedule.EndTime = &input.ScheduleDetails[i].EndTime
			scheduleTimes = append(scheduleTimes, *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].EndTime))
			for j := range input.ScheduleDetails[i].ClassroomIds {
				for k := range classrooms {
					if classrooms[k].ID == input.ScheduleDetails[i].ClassroomIds[j] {
						schedule.Classrooms = append(schedule.Classrooms, classrooms[k])
					}
				}
			}
			schedule.TeacherId = &input.ScheduleDetails[i].TeacherId
			schedule.AsistantId = input.ScheduleDetails[i].AsistantId
			schedule.ParentId = &newScheduleParent.ID
			schedule.ClassId = newScheduleParent.ClassId
			schedule.CenterId = newScheduleParent.CenterId
			schedule.CreatedBy = newScheduleParent.CreatedBy
			schedule.WorkSessionId = &input.ScheduleDetails[i].WorkSessionId
			schedule.Name = input.ScheduleDetails[i].Name
			schedule.Index = scheduleIndex
			newSchedule := schedule
			// append id
			scheduleIndex++
			listScheduleIdSync = append(listScheduleIdSync, newSchedule.ID)
			if len(input.ScheduleDetails[i].Childrens) > 0 {
				for j := range input.ScheduleDetails[i].Childrens {
					var childSchedule models.ScheduleClass
					childSchedule.ID = uuid.New()
					if input.ScheduleDetails[i].Childrens[j].Id != nil && !utils.Contains(listScheduleIdChanged, *input.ScheduleDetails[i].Childrens[j].Id) {
						childSchedule.ID = *input.ScheduleDetails[i].Childrens[j].Id
						listScheduleIdChanged = append(listScheduleIdChanged, childSchedule.ID)
					}
					listScheduleIdSync = append(listScheduleIdSync, childSchedule.ID)
					childSchedule.StartTime = &input.ScheduleDetails[i].Childrens[j].StartTime
					childSchedule.EndTime = &input.ScheduleDetails[i].Childrens[j].EndTime
					for k := range input.ScheduleDetails[i].Childrens[j].ClassroomIds {
						for h := range classrooms {
							if classrooms[h].ID == input.ScheduleDetails[i].Childrens[j].ClassroomIds[k] {
								childSchedule.Classrooms = append(childSchedule.Classrooms, classrooms[h])
							}
						}
					}
					childSchedule.TeacherId = &input.ScheduleDetails[i].Childrens[j].TeacherId
					childSchedule.AsistantId = input.ScheduleDetails[i].Childrens[j].AsistantId
					childSchedule.StartDate = schedule.StartDate
					childSchedule.ParentId = &newSchedule.ID
					childSchedule.ClassId = newSchedule.ClassId
					childSchedule.CenterId = newSchedule.CenterId
					childSchedule.CreatedBy = newSchedule.CreatedBy
					childSchedule.WorkSessionId = &input.ScheduleDetails[i].Childrens[j].WorkSessionId
					childSchedules = append(childSchedules, childSchedule)
					scheduleIndex++
					scheduleTimes = append(scheduleTimes, *utils.MixedDateAndTime(input.ScheduleDetails[i].StartDate, &input.ScheduleDetails[i].Childrens[j].EndTime))
				}
				newSchedule.Childrens = childSchedules
			}
			newSchedule.Childrens = append(newScheduleParent.Childrens, newSchedule)
		}
	}
	if len(listScheduleIdChanged) > 0 {
		schedulesOld, _ := repo.GetScheduleClassByIds(listScheduleIdChanged, user.CenterId)
		if len(schedulesOld) != len(listScheduleIdChanged) {
			return ResponseError(c, fiber.StatusBadRequest, "list record not found", consts.CreateFailed)
		}
	}
	lessons, err := repo.FilterDetailLessonByLive(true, input.ClassId, user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.DataNotFound)
	}
	for i := range lessons {
		listLessonsSync = append(listLessonsSync, lessons[i].Childrens...)
	}
	if len(listLessonsSync) != len(listScheduleIdSync) {
		return ResponseError(c, fiber.StatusInternalServerError, "total invalid", consts.ERROR_TOTAL_LESSONS_NOT_EQUAL)
	}
	for i := range listLessonsSync {
		listLessonsSync[i].ScheduleId = listScheduleIdSync[i]
	}
	//set end date class
	maxDate := lo.MaxBy(scheduleTimes, func(a, b time.Time) bool {
		return a.After(b)
	})
	class.EndAt = &maxDate
	// create schedule class and update lesson in class
	listSchedules, err := repo.CreateScheduleClass(&newScheduleParent, listScheduleIdChanged, isChangeType, user.CenterId, listLessonsSync, class)
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, "failed", consts.CreateFailed)
	}
	return ResponseSuccess(c, fiber.StatusCreated, "Success", listSchedules)
}
func handleGetDateNotInHoliday(Weeks []int, totalLessons int, startDate time.Time, holidays []models.Holiday) []time.Time {
	var weekRules []rrule.Weekday
	sort.Ints(Weeks) // Sắp xế	p mảng Weeks tăng dần để đảm bảo thứ tự ngày hợp lý.
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
func GetDetailScheduleClass(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}
	isParam := c.Params("Id")
	classId, err := uuid.Parse(isParam)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "Id class không hợp lệ", consts.InvalidReqInput)
	}
	//check lesson (table lesson) total equal than total lesson (class)
	lessons, err := repo.FilterLessonsByLive(true, classId, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "not found lessons in class", consts.DataNotFound)
	}
	class, err := repo.GetClassByIdAndCenterId(classId, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "not found class", consts.DataNotFound)
	}
	lessonLen := len(lessons)
	if lessonLen != int(class.TotalLessons) {
		lessonMissNum := strconv.Itoa(int(class.TotalLessons) - len(lessons))
		return ResponseError(c, fiber.StatusBadRequest, lessonMissNum, consts.ERROR_CAN_NOT_OPEN_WHEN_LESSON_TOTAL_NOT_EQUAL)
	}
	schedules, err := repo.GetDetailScheduleByClassId(classId, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "not found", consts.DataNotFound)
	}
	return ResponseSuccess(c, fiber.StatusOK, "Success", schedules)
}
func GetListScheduleClass(c *fiber.Ctx) error {
	user, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}
	var query consts.Query
	if err := c.QueryParser(&query); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid", consts.InvalidReqInput)
	}
	idParam := c.Params("id")
	classId, err := uuid.Parse(idParam)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "Id require1", consts.InvalidReqInput)
	}
	schedules, err := repo.GetListScheduleByClassId(classId, query, user)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, err.Error(), consts.DataNotFound)
	}
	return ResponseSuccess(c, fiber.StatusOK, "Success", schedules)
}

func GetListScheduleClassForStudent(c *fiber.Ctx) error {
	user, err := repo.GetTokenData(c)
	if err != nil {
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}
	var query consts.Query
	if err := c.QueryParser(&query); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "query không hợp lệ", consts.InvalidReqInput)
	}
	if user.RoleId == consts.Student {
		query.StudentId = user.ID.String()
	}
	studentSchedule, err := repo.GetScheduleClassForStudent(query, user)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, err.Error(), consts.InvalidReqInput)
	}
	return ResponseSuccess(c, fiber.StatusOK, "Success", studentSchedule)
}
