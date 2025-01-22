package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type Student struct {
	Model               `gorm:"embedded"`
	CenterId            uuid.UUID              `json:"-"`
	BranchId            *uuid.UUID             `json:"-"`
	Type                int64                  `json:"type,omitempty"`
	Avatar              string                 `json:"avatar,omitempty"`
	SignatureImg        string                 `json:"signature_img,omitempty"`
	Signature           string                 `json:"signature,omitempty"`
	FullName            string                 `gorm:"index" json:"full_name"`
	Username            string                 `json:"-"`
	Status              int64                  `json:"status,omitempty" gorm:"default:1"`
	Gender              string                 `json:"gender,omitempty"`
	DOB                 *time.Time             `json:"dob,omitempty"`
	Email               string                 `json:"email,omitempty"`
	EmailVerified       bool                   `json:"-" gorm:"default:false"`
	Phone               string                 `json:"phone,omitempty"`
	ProvinceId          *int64                 `json:"province_id,omitempty"`
	DistrictId          *int64                 `json:"district_id,omitempty"`
	Address             string                 `json:"address,omitempty"`
	CustomerSourceID    *uuid.UUID             `json:"customer_source_id,omitempty"` // ID Nguồn khách hàng
	ContactChannelID    *uuid.UUID             `json:"contact_channel_id,omitempty"` // ID Kênh liên hệ
	CustomerSource      *CustomerSource        `json:"customer_source,omitempty" gorm:"foreignKey:CustomerSourceID"`
	ContactChannel      *ContactChannel        `json:"contact_channel,omitempty" gorm:"foreignKey:ContactChannelID"`
	AddedAt             *time.Time             `json:"AddedAt,omitempty" gorm:"-"`                   //Ngày thêm vào lớp học
	IsOfficialAt        time.Time              `json:"is_official_at,omitempty" gorm:"default:null"` //Ngày trở thành học viên chính thức
	IsTrialAt           *time.Time             `json:"is_trial_at,omitempty" gorm:"default:null"`    //Ngày bắt đầu học thử
	TotalTrialSession   int64                  `json:"total_trial_session,omitempty" gorm:"-"`       // Số buổi học thử
	Province            *Province              `json:"province,omitempty" gorm:"foreignKey:ProvinceId"`
	District            *District              `json:"district,omitempty" gorm:"foreignKey:DistrictId"`
	StudyNeeds          *StudyNeeds            `json:"study_needs,omitempty" gorm:"foreignKey:StudentId"` // nhu cầu học của từng sinh viên
	Subjects            []*Subject             `json:"-" gorm:"many2many:student_subjects"`
	Logs                []*StudentLog          `json:"action_logs,omitempty" gorm:"foreignKey:StudentId"`
	Curriculums         []*Curriculum          `json:"-" gorm:"many2many:student_curriculums"`
	Classes             []*Class               `json:"classes,omitempty" gorm:"many2many:student_classes"`
	SessionAttendances  *[]SessionAttendance   `json:"session_attendances,omitempty" gorm:"foreignKey:StudentId;references:ID"` // thông tin điểm danh của sinh viên
	Exams               []*ExamResult          `json:"exam,omitempty" gorm:"foreignKey:StudentId;references:ID"`
	CompletedSubject    string                 `json:"completed_subject,omitempty" gorm:"-"`       // thông tin môn học đã hoàn thành
	CaregiverId         *uuid.UUID             `json:"caregiver_id,omitempty" gorm:"default:null"` // ID người chăm sóc
	Caregiver           *User                  `json:"caregiver,omitempty" gorm:"foreignKey:CaregiverId"`
	CareResult          string                 `json:"care_result,omitempty" gorm:"-"`
	Branch              *Branch                `json:"branch,omitempty" gorm:"foreignKey:BranchId"`
	StudentCertificates *[]StudentCertificates `gorm:"foreignKey:StudentId" json:"student_certificates,omitempty"`
}

// StudentNeeds Nhu cầu học tập
type StudyNeeds struct {
	Model               `gorm:"embedded"`
	StudentId           uuid.UUID      `json:"student_id" gorm:"not null"` // ID học viên
	CenterId            uuid.UUID      `json:"-"`
	BranchId            *uuid.UUID     `json:"branch_id"`
	EnrollmentId        uuid.UUID      `json:"enrollment_id" gorm:"default:null"`                   // ID Kế hoạch tuyển sinh
	StudyGoals          string         `json:"study_goals,omitempty" gorm:"type:longtext"`          // mục tiêu học tập của sinh viên
	TeacherRequirements string         `json:"teacher_requirements,omitempty" gorm:"type:longtext"` // yêu cầu cụ thể từ giáo viên về học viên
	IsOnlineForm        *bool          `json:"is_online_form,omitempty"`
	IsOfflineForm       *bool          `json:"is_offline_form,omitempty"`
	StudyingStartDate   *time.Time     `json:"studying_start_date,omitempty"`
	TimeSlots           []TimeSlot     `json:"time_slots,omitempty" gorm:"-"`   // Khung giờ các ca trong lịch
	ShortShifts         []ShortShift   `json:"short_shifts,omitempty" gorm:"-"` // lưu trữ các ca ngắn trong lịch học
	SubjectIds          []uuid.UUID    `json:"subject_ids,omitempty" gorm:"-"`
	CurriculumIds       []uuid.UUID    `json:"curriculum_ids,omitempty" gorm:"-"` // lưu giữ chương trình học mà học viên muốn theo học
	Metadata            datatypes.JSON `json:"metadata,omitempty"`

	Branch      *Branch          `json:"branch,omitempty" gorm:"foreignKey:BranchId"`
	Schedule    *StudentSchedule `json:"schedule,omitempty" gorm:"foreignKey:StudentId;references:StudentId"`
	Subjects    []*Subject       `json:"subjects,omitempty" gorm:"-"`
	Curriculums []*Curriculum    `json:"curriculums,omitempty" gorm:"-"`
	Student     *Student         `json:"-" gorm:"foreignKey:StudentId"`
	Enrollment  *EnrollmentPlan  `json:"enrollment,omitempty" gorm:"foreignKey:EnrollmentId"`
}

// StudentSchedule Lịch học viên
type StudentSchedule struct {
	Model     `gorm:"embedded"`
	StudentId uuid.UUID  `json:"student_id"`
	CenterId  *uuid.UUID `json:"center_id,omitempty"`

	TimeSlots   []TimeSlot   `json:"time_slots" gorm:"foreignKey:ScheduleId"`
	StudyShifts []Shift      `json:"-" gorm:"foreignKey:ScheduleId"`
	ShortShifts []ShortShift `json:"short_shifts" gorm:"-"`
}
