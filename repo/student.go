package repo

import (
	"context"
	"errors"
	"fmt"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	Student  models.Student
	Students []*models.Student
)

func (u *Student) VerifyEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Model(&models.Student{}).
		Where("email = ?", email).Update("email_verified", true).Error
}

func (u *Student) First(query interface{}, args []interface{}, preload ...string) error {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		err         error
		DB          = app.Database.DB.WithContext(ctx).Where(query, args...)
	)
	defer cancel()
	if len(preload) > 0 {
		u.PreloadStudent(DB, preload...)
		err = DB.First(&u).Error
	} else {
		err = DB.First(&u).Error
	}
	return err
}

func (u *Student) PreloadStudent(DB *gorm.DB, properties ...string) {
	for _, v := range properties {
		if v == "StudyNeeds" {
			DB.Preload("StudyNeeds", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, branch_id, student_id, studying_start_date, enrollment_id").
					Preload("Enrollment", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "name")
					})
			})
		}
		if v == "Caregiver" {
			DB.Preload("Caregiver", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "full_name")
			})
		}
		if v == "Source" {
			DB.Preload("CustomerSource", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name")
			})
		}
		if v == "Province" {
			DB.Preload("Province", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name")
			})
		}
		if v == "District" {
			DB.Preload("District", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "prefix")
			})
		}
		if v == "CustomerSource" {
			DB.Preload("CustomerSource", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name")
			})
		}
		if v == "ContactChannel" {
			DB.Preload("ContactChannel", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name")
			})
		}
	}
}

func (u *Student) Create() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()

	username := utils.GenerateUniqueUsername()
	u.Username = username
	if u.Type == consts.Official || u.Type == consts.Trial {
		tx := app.Database.DB.WithContext(ctx).Begin() // transaction
		if err = tx.Create(&u).Error; err != nil {
			tx.Rollback()
			return err
		} // lưu thông tin vào bảng student
		var (
			pwd = app.Config("DEFAULT_PASSWORD")
		)
		temp, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		loginInfo := LoginInfo{
			ID:           u.ID,
			CenterID:     u.CenterId,
			Username:     username,
			Phone:        u.Phone,
			Email:        u.Email,
			PasswordHash: string(temp),
			RoleId:       consts.Student,
			DeletedAt:    nil,
		}
		if err = loginInfo.Create(); err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit().Error
	}
	return app.Database.DB.WithContext(ctx).Model(&models.Student{}).Create(&u).Error
}
func (u *Student) Update(origin Student, query interface{}, args []interface{}) (err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		tx          = app.Database.DB.WithContext(ctx).Begin()
	)
	defer cancel()

	if err = tx.Model(&models.Student{}).Where(query, args...).Updates(&u).Error; err != nil {
		logrus.Error(err)
		tx.Rollback()
		return err
	}

	if err = app.Database.DB.WithContext(ctx).Model(&models.LoginInfo{}).Where("id = ?", u.ID).
		Update("phone", u.Phone).Update("email", u.Email).Error; err != nil {
		logrus.Error(err)
		tx.Rollback()
		return err
	}

	if !u.EmailVerified {
		if err = tx.Model(&models.Student{}).Where(query, args...).
			Update("EmailVerified", false).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (u *Student) Count(DB *gorm.DB) (count int64) {
	var ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	DB.Model(&models.Student{}).WithContext(ctx).Count(&count)
	return
}

func (u *Student) Find(DB *gorm.DB, preload ...string) (Students, error) {
	var (
		entries     Students
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()

	if len(preload) > 0 {
		u.PreloadStudent(DB, preload...)
	}

	err := DB.WithContext(ctx).Model(&models.Student{}).Find(&entries)
	return entries, err.Error
}

func PreloadTotalTrialSession(entry *models.Student) {
	var ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	app.Database.DB.WithContext(ctx).Where("student_id = ?", entry.ID).
		Model(&models.StudentSession{}).Count(&entry.TotalTrialSession)
}

func LoadCareResult(studentId uuid.UUID) (result string) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	app.Database.DB.WithContext(ctx).Raw("SELECT cs.`name` FROM `care_infos` ci\n"+
		"JOIN care_results cs ON cs.id = ci.result_id\n"+
		"WHERE ci.student_id = ? AND ci.deleted_at IS NULL\n"+
		"ORDER BY ci.created_at DESC\nLIMIT 1", studentId).Scan(&result)
	return
}

func (u *Student) PreloadCompletedSubject(studentId uuid.UUID) string {
	var ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	var result struct {
		Completed string `json:"completed"`
		Total     string `json:"total"`
	}

	if err := app.Database.DB.WithContext(ctx).Raw("WITH t1 AS (\n"+
		"SELECT COUNT(*) FROM `student_subjects`\n"+
		"WHERE student_id = ?), t2 AS (\n"+
		"SELECT COUNT(c.id) FROM `student_classes` sc\n"+
		"JOIN classes c ON c.id = sc.class_id\n"+
		"WHERE sc.student_id = ? AND c.end_at <= NOW() AND c.deleted_at IS NULL)\n\n"+
		"SELECT (SELECT * FROM t2) completed, (SELECT * FROM t1) total", studentId, studentId).
		Scan(&result).Error; err != nil {
		logrus.Error(err)
	}

	return fmt.Sprintf("%s/%s", result.Completed, result.Total)
}
func (u *Student) Delete(studentIds []uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//jobrunner.Now(DeleteCasdoorUsers{IDs: studentIds})

		if err = tx.Where("id IN ?", studentIds).Delete(&models.Student{}).Error; err != nil {
			logrus.Error(err)
			return err
		}
		if err = tx.Where("id IN ?", studentIds).Delete(&models.LoginInfo{}).Error; err != nil {
			logrus.Error(err)
			return err
		}
		if err = tx.Where("student_id IN ?", studentIds).Delete(&models.StudyNeeds{}).Error; err != nil {
			logrus.Error(err)
			return err
		}

		return nil
	})
}
func GetStudentsBySubjectAndCenterId(subjectId, centerId uuid.UUID) ([]models.Student, error) {
	var students []models.Student
	db := app.Database.DB.Model(&models.Student{}).Select("students.`id`", "ss.*").Where("center_id = ?", centerId)
	db.Joins("INNER JOIN student_subjects as ss ON students.`id` = ss.`student_id`")
	db.Where("ss.`subject_id` = ?", subjectId)
	db.Find(&students)
	return students, db.Error
}

func CheckStudentExists(studentID uuid.UUID) error {
	var student Student
	if err := app.Database.DB.Where("id = ?", studentID).First(&student).Error; err != nil {
		return err
	}
	return nil
}
func AddStudentToClass(input []models.StudentToClass, token TokenData, c *fiber.Ctx) error {
	tx := app.Database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var err error
	for _, dataInput := range input {
		var class models.Class
		if err := tx.Model(&models.Class{}).Where("id = ? AND center_id = ?", dataInput.ClassId, token.CenterId).First(&class).Error; err != nil {
			tx.Rollback()
			return err
		}
		if class.Status == consts.CLASS_CANCELED {
			tx.Rollback()
			return errors.New("class is canceled")
		}
		var student models.Student
		if err := tx.Model(&models.Student{}).Where("id = ?", dataInput.StudentId).First(&student).Error; err != nil {
			tx.Rollback()
			return err
		}
		now := time.Now()
		studentClassInfo := models.StudentClasses{
			ClassId:   class.ID,
			StudentId: student.ID,
			CreatedAt: &now,
		}
		if err := tx.Model(&models.StudentClasses{}).Create(&studentClassInfo).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
		if student.Status, err = LoadStatusTransaction(tx, student.ID, token.CenterId); err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func LoadStatusTransaction(tx *gorm.DB, studentId, centerId uuid.UUID) (status int64, err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		count       int64
	)
	defer cancel()

	if err = tx.WithContext(ctx).Raw(`--
    SELECT COUNT(*) FROM (
        SELECT s.code FROM student_subjects
        JOIN subjects s ON s.id = subject_id AND s.deleted_at IS NULL AND s.center_id = ?
        WHERE student_id = ?
        UNION
        SELECT t1.code FROM student_curriculums sc
        JOIN (
            SELECT cs.curriculum_id, s.code FROM curriculum_subjects cs
            JOIN subjects s ON s.id = cs.subject_id AND s.deleted_at IS NULL AND s.center_id = ?
        ) t1 ON t1.curriculum_id = sc.curriculum_id
        WHERE student_id = ?) temp`, centerId, studentId, centerId, studentId).Scan(&count).Error; err != nil {
		logrus.Error(err)
		return 0, err
	}
	if count == 0 {
		return consts.Pending, nil
	}

	if err = tx.WithContext(ctx).Raw(`--
    SELECT COUNT(*) FROM (
        SELECT s.code FROM student_subjects
        JOIN subjects s ON s.id = subject_id AND s.deleted_at IS NULL AND s.center_id = ?
        WHERE student_id = ?
        UNION
        SELECT t1.code FROM student_curriculums sc
        JOIN (
            SELECT cs.curriculum_id, s.code FROM curriculum_subjects cs
            JOIN subjects s ON s.id = cs.subject_id AND s.deleted_at IS NULL AND s.center_id = ?
        ) t1 ON t1.curriculum_id = sc.curriculum_id
        WHERE student_id = ?
    ) s
    LEFT JOIN (
        SELECT c.subject_code FROM student_classes sc
        JOIN (
            SELECT c.id, s.code subject_code FROM classes c
            JOIN subjects s ON s.id = c.subject_id AND s.deleted_at IS NULL AND s.center_id = ?
            WHERE c.deleted_at IS NULL
        ) c ON c.id = sc.class_id
        WHERE student_id = ?
    ) c ON c.subject_code = s.code
    WHERE c.subject_code IS NULL
    `, centerId, studentId, centerId, studentId, centerId, studentId).Scan(&count).Error; err != nil {
		logrus.Error(err)
		return 0, err
	}
	if count > 0 {
		return consts.Pending, nil
	}

	if err = tx.WithContext(ctx).Raw(`--
    SELECT COUNT(*) FROM student_classes sc
    JOIN classes c ON c.id = class_id AND c.deleted_at IS NULL AND c.center_id = ?
    WHERE student_id = ? AND sc.status IS NULL`, centerId, studentId).Scan(&count).Error; err != nil {
		logrus.Error(err)
		return 0, err
	}
	if count == 0 {
		return consts.Reserved, nil
	}

	return consts.Studying, nil
}
