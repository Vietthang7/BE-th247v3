package repo

import (
	"encoding/json"
	"fmt"
	"intern_247/app"
	"intern_247/models"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// CreateTeachingSchedule xử lý logic tạo lịch giảng dạy
func CreateTeachingSchedule(form models.CreateTeachScheForm) (*models.TeachingSchedule, error) {
	// Kiểm tra UserId có hợp lệ không
	if form.UserId == uuid.Nil {
		return nil, fmt.Errorf("UserId is required")
	}

	// Kiểm tra User có tồn tại không
	var user models.User
	if err := app.Database.DB.First(&user, "id = ?", form.UserId).Error; err != nil {
		return nil, fmt.Errorf("User not found")
	}

	// Chuyển đổi StartDate và EndDate sang time.Time
	startDate, err := time.Parse("2006-01-02", form.StartDate)
	if err != nil {
		return nil, fmt.Errorf("%s", "Invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", form.EndDate)
	if err != nil {
		return nil, fmt.Errorf("%s", "Invalid end date format")
	}

	// Tạo TeachingSchedule mới
	newSchedule := models.TeachingSchedule{
		UserId:     form.UserId,
		StartDate:  startDate,
		EndDate:    endDate,
		IsOnline:   form.IsOnline,
		IsOffline:  form.IsOffline,
		Notes:      form.Notes,
		TimeSlots:  form.TimeSlots,
		UserShifts: form.UserShifts,
	}

	// Lưu TeachingSchedule vào database
	if err := app.Database.DB.Create(&newSchedule).Error; err != nil {
		return nil, fmt.Errorf("%s", "Failed to create teaching schedule")
	}

	// Parse danh sách TimeSlots
	var timeSlots []models.TimeSlot
	if err := json.Unmarshal(form.TimeSlots, &timeSlots); err != nil {
		return nil, fmt.Errorf("%s", "Invalid time_slots format")
	}

	// Kiểm tra từng WorkSessionId trước khi lưu TimeSlots
	for i := range timeSlots {
		var workSession models.WorkSession
		if err := app.Database.DB.First(&workSession, "id = ? AND is_active = ?", timeSlots[i].WorkSessionId, true).Error; err != nil {
			return nil, fmt.Errorf("%s", "Invalid or inactive work session")
		}

		timeSlots[i].ScheduleId = newSchedule.ID
		timeSlots[i].UserId = &form.UserId
	}

	// Lưu TimeSlots vào database
	if len(timeSlots) > 0 {
		if err := app.Database.DB.Create(&timeSlots).Error; err != nil {
			return nil, fmt.Errorf("%s", "Failed to create time slots")
		}
	}

	// Kiểm tra TeachingSchedule đã tồn tại chưa
	var existingSchedule models.TeachingSchedule
	if err := app.Database.DB.First(&existingSchedule, "user_id = ?", form.UserId).Error; err != nil {
		return nil, fmt.Errorf("%s", "Teaching schedule not found")
	}

	// Parse danh sách UserShifts
	var rawUserShifts []models.Shift
	if err := json.Unmarshal(form.UserShifts, &rawUserShifts); err != nil {
		return nil, fmt.Errorf("%s", "Invalid user_shifts format")
	}

	// Chuyển đổi weekday của Go sang hệ thống
	convertWeekday := func(w time.Weekday) int {
		if w == 0 {
			return 1 // Chủ Nhật = 1
		}
		return int(w) + 1 // Thứ Hai = 2, Thứ Ba = 3, ...
	}

	// Lưu danh sách các Shift
	var userShifts []models.Shift

	for _, rawShift := range rawUserShifts {
		// Kiểm tra WorkSessionId có tồn tại và active không
		var workSession models.WorkSession
		if err := app.Database.DB.First(&workSession, "id = ? AND is_active = ?", rawShift.WorkSessionId, true).Error; err != nil {
			return nil, fmt.Errorf("%s", "Invalid or inactive work session")
		}

		// Lặp qua các ngày trong khoảng StartDate - EndDate
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			if convertWeekday(d.Weekday()) == rawShift.DayOfWeek {
				// Tìm TimeSlotId phù hợp với WorkSessionId của newSchedule
				var timeSlotId uuid.UUID
				for _, ts := range timeSlots {
					if ts.WorkSessionId == rawShift.WorkSessionId {
						timeSlotId = ts.ID
						break
					}
				}

				userShifts = append(userShifts, models.Shift{
					ScheduleId:    newSchedule.ID, // ✅ Dùng schedule mới tạo
					UserId:        &form.UserId,
					WorkSessionId: rawShift.WorkSessionId,
					DayOfWeek:     rawShift.DayOfWeek,
					Date:          d,
					Type:          "user",
					TimeSlotId:    timeSlotId, // ✅ Liên kết với TimeSlot từ schedule mới
				})
			}
		}
	}

	// Lưu danh sách Shift vào database
	if len(userShifts) > 0 {
		if err := app.Database.DB.Create(&userShifts).Error; err != nil {
			return nil, fmt.Errorf("%s", "Failed to create user shifts")
		}
	}

	return &newSchedule, nil
}

// DeleteTeachingSchedule xóa lịch giảng dạy và các liên kết liên quan
func DeleteTeachingSchedule(scheduleID uuid.UUID) error {
	// Kiểm tra xem schedule có tồn tại không
	var schedule models.TeachingSchedule
	if err := app.Database.DB.First(&schedule, "id = ?", scheduleID).Error; err != nil {
		return err
	}

	// Xóa tất cả TimeSlots liên quan
	if err := app.Database.DB.Where("schedule_id = ?", scheduleID).Delete(&models.TimeSlot{}).Error; err != nil {
		return err
	}

	// Xóa tất cả UserShifts liên quan
	if err := app.Database.DB.Where("schedule_id = ?", scheduleID).Delete(&models.Shift{}).Error; err != nil {
		return err
	}

	// Xóa TeachingSchedule
	if err := app.Database.DB.Delete(&schedule).Error; err != nil {
		return err
	}

	return nil
}

func ReadTeachSchedule(scheduleID uuid.UUID) (*models.TeachingSchedule, error) {
	var schedule models.TeachingSchedule

	// Truy vấn TeachingSchedule
	result := app.Database.DB.First(&schedule, "id = ?", scheduleID)
	if result.Error != nil {
		return nil, result.Error
	}

	// Struct chỉ chứa các field cần thiết (bỏ created_at, updated_at)
	type TimeSlotResponse struct {
		ID            uuid.UUID `json:"id"`
		ScheduleId    uuid.UUID `json:"schedule_id"`
		WorkSessionId uuid.UUID `json:"work_session_id"`
		StartTime     string    `json:"start_time"`
		EndTime       string    `json:"end_time"`
	}

	type ShiftResponse struct {
		ID            uuid.UUID `json:"id"`
		ScheduleId    uuid.UUID `json:"schedule_id"`
		WorkSessionId uuid.UUID `json:"work_session_id"`
		DayOfWeek     int       `json:"day_of_week"`
	}

	// Lấy TimeSlots từ bảng TimeSlot
	var timeSlots []TimeSlotResponse
	if err := app.Database.DB.Model(&models.TimeSlot{}).
		Select("id, schedule_id, work_session_id, start_time, end_time").
		Where("schedule_id = ?", schedule.ID).
		Find(&timeSlots).Error; err != nil {
		return nil, err
	}

	// Lấy UserShifts từ bảng Shift
	var userShifts []ShiftResponse
	if err := app.Database.DB.Model(&models.Shift{}).
		Select("id, schedule_id, work_session_id, day_of_week").
		Where("schedule_id = ?", schedule.ID).
		Find(&userShifts).Error; err != nil {
		return nil, err
	}

	// Chuyển danh sách về JSON
	timeSlotsJSON, err := json.Marshal(timeSlots)
	if err != nil {
		return nil, err
	}

	userShiftsJSON, err := json.Marshal(userShifts)
	if err != nil {
		return nil, err
	}

	// Gán JSON vào struct
	schedule.TimeSlots = datatypes.JSON(timeSlotsJSON)
	schedule.UserShifts = datatypes.JSON(userShiftsJSON)

	return &schedule, nil
}

func ListTeachSchedule() ([]models.TeachingSchedule, error) {
	var schedules []models.TeachingSchedule

	// Lấy danh sách TeachingSchedule từ database
	if err := app.Database.DB.Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

// UpdateTeachSchedule cập nhật lịch giảng dạy
func UpdateTeachSchedule(scheduleID uuid.UUID, form models.CreateTeachScheForm) (*models.TeachingSchedule, error) {
	// Kiểm tra lịch giảng dạy có tồn tại không
	var schedule models.TeachingSchedule
	if err := app.Database.DB.First(&schedule, "id = ?", scheduleID).Error; err != nil {
		log.Printf("Error: Schedule not found for ID %v: %v", scheduleID, err)
		return nil, fmt.Errorf("%s", "Teaching schedule not found")
	}

	// Kiểm tra UserId có hợp lệ không
	if form.UserId == uuid.Nil {
		log.Println("Error: UserId is required")
		return nil, fmt.Errorf("UserId is required")
	}

	// Kiểm tra User có tồn tại không
	var user models.User
	if err := app.Database.DB.First(&user, "id = ?", form.UserId).Error; err != nil {
		log.Printf("Error: User not found with ID %v: %v", form.UserId, err)
		return nil, fmt.Errorf("User not found")
	}

	// Chuyển đổi StartDate và EndDate sang time.Time
	startDate, err := time.Parse("2006-01-02", form.StartDate)
	if err != nil {
		return nil, fmt.Errorf("%s", "Invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", form.EndDate)
	if err != nil {
		return nil, fmt.Errorf("%s", "Invalid end date format")
	}

	// Cập nhật thông tin lịch giảng dạy
	schedule.UserId = form.UserId
	schedule.StartDate = startDate
	schedule.EndDate = endDate
	schedule.IsOnline = form.IsOnline
	schedule.IsOffline = form.IsOffline
	schedule.Notes = form.Notes

	// Lưu cập nhật vào database
	if err := app.Database.DB.Save(&schedule).Error; err != nil {
		return nil, fmt.Errorf("%s", "Failed to update teaching schedule")
	}

	// Xóa TimeSlots và UserShifts cũ
	app.Database.DB.Where("schedule_id = ?", schedule.ID).Delete(&models.TimeSlot{})
	app.Database.DB.Where("schedule_id = ?", schedule.ID).Delete(&models.Shift{})

	// Parse danh sách TimeSlots
	var timeSlots []models.TimeSlot
	if err := json.Unmarshal(form.TimeSlots, &timeSlots); err != nil {
		return nil, fmt.Errorf("%s", "Invalid time_slots format")
	}

	// Kiểm tra từng WorkSessionId trước khi lưu TimeSlots
	for i := range timeSlots {
		var workSession models.WorkSession
		if err := app.Database.DB.First(&workSession, "id = ? AND is_active = ?", timeSlots[i].WorkSessionId, true).Error; err != nil {
			return nil, fmt.Errorf("%s", "Invalid or inactive work session")
		}

		timeSlots[i].ScheduleId = schedule.ID
		timeSlots[i].UserId = &form.UserId
	}

	// Lưu TimeSlots vào database
	if len(timeSlots) > 0 {
		if err := app.Database.DB.Create(&timeSlots).Error; err != nil {
			return nil, fmt.Errorf("%s", "Failed to create time slots")
		}
	}

	// Parse danh sách UserShifts
	var rawUserShifts []models.Shift
	if err := json.Unmarshal(form.UserShifts, &rawUserShifts); err != nil {
		return nil, fmt.Errorf("%s", "Invalid user_shifts format")
	}

	// Chuyển đổi weekday của Go sang hệ thống
	convertWeekday := func(w time.Weekday) int {
		if w == 0 {
			return 1 // Chủ Nhật = 1
		}
		return int(w) + 1 // Thứ Hai = 2, Thứ Ba = 3, ...
	}

	// Lưu danh sách Shift
	var userShifts []models.Shift

	for _, rawShift := range rawUserShifts {
		// Kiểm tra WorkSessionId có tồn tại và active không
		var workSession models.WorkSession
		if err := app.Database.DB.First(&workSession, "id = ? AND is_active = ?", rawShift.WorkSessionId, true).Error; err != nil {
			return nil, fmt.Errorf("%s", "Invalid or inactive work session")
		}

		// Lặp qua các ngày trong khoảng StartDate - EndDate
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			if convertWeekday(d.Weekday()) == rawShift.DayOfWeek {
				// Tìm TimeSlotId phù hợp với WorkSessionId của lịch giảng dạy
				var timeSlotId uuid.UUID
				for _, ts := range timeSlots {
					if ts.WorkSessionId == rawShift.WorkSessionId {
						timeSlotId = ts.ID
						break
					}
				}

				userShifts = append(userShifts, models.Shift{
					ScheduleId:    schedule.ID, // ✅ Dùng schedule hiện có
					UserId:        &form.UserId,
					WorkSessionId: rawShift.WorkSessionId,
					DayOfWeek:     rawShift.DayOfWeek,
					Date:          d,
					Type:          "user",
					TimeSlotId:    timeSlotId, // ✅ Liên kết với TimeSlot
				})
			}
		}
	}

	// Lưu danh sách Shift vào database
	if len(userShifts) > 0 {
		if err := app.Database.DB.Create(&userShifts).Error; err != nil {
			return nil, fmt.Errorf("%s", "Failed to create user shifts")
		}
	}

	// 🔹 **Lấy danh sách TimeSlots sau khi cập nhật chỉ với các trường cần thiết**
	var updatedTimeSlots []struct {
		ID            uuid.UUID `json:"id"`
		ScheduleID    uuid.UUID `json:"schedule_id"`
		WorkSessionID uuid.UUID `json:"work_session_id"`
		StartTime     string    `json:"start_time"`
		EndTime       string    `json:"end_time"`
	}

	app.Database.DB.Model(&models.TimeSlot{}).
		Where("schedule_id = ?", schedule.ID).
		Select("id, schedule_id, work_session_id, start_time, end_time").
		Find(&updatedTimeSlots)

	// 🔹 **Lấy danh sách UserShifts sau khi cập nhật chỉ với các trường cần thiết**
	var updatedUserShifts []struct {
		ID            uuid.UUID `json:"id"`
		WorkSessionID uuid.UUID `json:"work_session_id"`
		DayOfWeek     int       `json:"day_of_week"`
	}

	app.Database.DB.Model(&models.Shift{}).
		Where("schedule_id = ?", schedule.ID).
		Select("id, work_session_id, day_of_week").
		Find(&updatedUserShifts)

	// Chuyển đổi thành JSON để gán vào schedule
	timeSlotsJSON, err := json.Marshal(updatedTimeSlots)
	if err != nil {
		return nil, fmt.Errorf("%s", "Failed to marshal time_slots")
	}

	userShiftsJSON, err := json.Marshal(updatedUserShifts)
	if err != nil {
		return nil, fmt.Errorf("%s", "Failed to marshal user_shifts")
	}

	schedule.TimeSlots = datatypes.JSON(timeSlotsJSON)
	schedule.UserShifts = datatypes.JSON(userShiftsJSON)

	return &schedule, nil
}
