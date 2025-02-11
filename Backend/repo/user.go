package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/utils"
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	User models.User
)

func VerifyUserEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Update("EmailVerified", true).Error
}
func RegisterUser(entry *models.User, args map[string]interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()

	tx := app.Database.DB.WithContext(ctx).Begin()
	centerInfo := models.Center{
		IsActive: entry.IsActive,
		Email:    entry.Email,
		Phone:    entry.Phone,
	}
	// Tạo trung tâm đơn vị
	if err = tx.Model(&models.Center{}).Create(&centerInfo).Error; err != nil {
		logrus.Error(fmt.Sprint("Lỗi tạo trung tâm : %s", err.Error()))
		tx.Rollback()
		return err
	}
	entry.CenterId = &centerInfo.ID
	// Tạo user
	if err = tx.Create(entry).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Tạo thông tin đăng nhập
	pwd := args["password"].(string)
	temp, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	loginInfo := LoginInfo{
		ID:           entry.ID,
		CenterID:     *entry.CenterId,
		Username:     entry.Username,
		Phone:        entry.Phone,
		Email:        entry.Email,
		RoleId:       entry.RoleId,
		PasswordHash: string(temp),
		DeletedAt:    nil,
	}
	if err = loginInfo.Create(); err != nil {
		tx.Rollback()
		return err
	}
	// Tạo các công thức lương mặc định
	var salaryStatements []*models.SalaryStatement
	if err := tx.Model(&models.SalaryStatement{}).Where("center_id IS NULL").Where("is_default = ?", true).Find(&salaryStatements).Error; err != nil {
		logrus.Error(fmt.Sprintf("Get salary statement failed: %s", err.Error()))
		//if ok, err := app.SSOClient.DeleteUser(casUser); err != nil || !ok {
		//	logrus.Error(fmt.Sprintf("Delete casdoor user failed: %s", err.Error()))
		//	tx.Rollback()
		//	return err
		//}
		tx.Rollback()
		return err
	}
	isDefault := false
	var newSalaryStatements []*models.SalaryStatement
	for _, stmt := range salaryStatements {
		newStmt := &models.SalaryStatement{
			Title:      stmt.Title,
			SalaryType: stmt.SalaryType,
			ObjectType: stmt.ObjectType,
			IsActive:   stmt.IsActive,
			CenterId:   &centerInfo.ID,
			IsDefault:  isDefault,
		}
		newSalaryStatements = append(newSalaryStatements, newStmt)
	}
	if err := tx.Model(&models.SalaryStatement{}).Create(&newSalaryStatements).Error; err != nil {
		logrus.Error(fmt.Sprintf("Bulk create salary statements failed: %s", err.Error()))
		tx.Rollback()
		return err
	}
	now := time.Now()
	// Tạo cài dặt goi học phí mặc định
	var tuitionPkgCf TuitionFeePkgCf
	if err = tuitionPkgCf.First("is_template", []interface{}{true}); err != nil {
		logrus.Errorf("Error get tuition fee pkg cf default data: %v", err)
	} else {
		tuitionPkgCf.ID = uuid.New()
		tuitionPkgCf.CreatedAt = &now
		tuitionPkgCf.UpdatedAt = &now
		tuitionPkgCf.CenterID = &centerInfo.ID
		tuitionPkgCf.IsTemplate = false
		if err = tuitionPkgCf.Create(tx); err != nil {
			tx.Rollback()
			logrus.Errorf("Error create tuition pkg cf: %v", err)
			return err
		}
	}
	// Tạo chăm sóc mặc định
	var careAssignment CareAssignment
	if err = careAssignment.First("is_template", []interface{}{true}); err != nil {
		logrus.Errorf("Error get care assignment default data: %v", err)
	} else {
		careAssignment.ID = uuid.New()
		careAssignment.CreatedAt = &now
		careAssignment.UpdatedAt = &now
		careAssignment.CenterID = &centerInfo.ID
		careAssignment.IsTemplate = false
		if err = careAssignment.Create(tx); err != nil {
			tx.Rollback()
			logrus.Errorf("Error create care assignment: %v", err)
			return err
		}
	}

	return tx.Commit().Error
}

func GetUserByID(user_id uuid.UUID) (models.User, int64, error) {
	var user models.User
	query := app.Database.DB.Where("id = ?", user_id).First(&user)
	return user, query.RowsAffected, query.Error
}
func HasPermission(entry models.User, subject, action string) bool {
	var (
		err   error
		group models.PermissionGroup
	)
	if group, err = FirstPermissionGrp(app.Database.DB.Where("id = ?", entry.PermissionGrpId)); err != nil {
		logrus.Error(err)
		return false
	}
	if !*group.IsActive {
		return false
	}
	if group.SelectAll != nil && *group.SelectAll {
		return true
	} // nếu SelectAll tồn tại và có giá trị true
	var permissionIds []uuid.UUID
	if err = json.Unmarshal(group.PermissionIds, &permissionIds); err != nil {
		logrus.Error(err)
		return false
	}
	// Tạo đối tượng truy vấn với điều kiện
	query := app.Database.DB.Where(map[string]interface{}{
		"deleted_at": nil,
		"subject":    subject,
		"action":     action,
	}).Where("id IN (?)", permissionIds)

	// Gọi hàm GetPermissionId với truy vấn đã tạo
	if _, err = GetPermissionId(query); err == nil {
		return true
	}
	return false
}
func HasPermission2(entry models.User, action string, subject ...string) bool {
	var (
		err   error
		group models.PermissionGroup
	)
	if group, err = FirstPermissionGrp(app.Database.DB.Where("id = ?", entry.PermissionGrpId)); err != nil {
		logrus.Error(err)
		return false
	}
	if !*group.IsActive {
		return false
	}
	if group.SelectAll != nil && *group.SelectAll {
		return true
	}
	var permissionIds []uuid.UUID
	if err = json.Unmarshal(group.PermissionIds, &permissionIds); err != nil {
		logrus.Error(err)
		return false
	}
	if _, err = GetPermissionId(app.Database.DB.Where("id IN ? AND subject IN ? AND action = ?", permissionIds, subject, action)); err == nil {
		return true
	}
	return false
}

func (u *User) First(query interface{}, args []interface{}, preload ...string) error {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = app.Database.DB.WithContext(ctx).Where(query, args...)
	)
	defer cancel()
	if len(preload) > 0 {
		NewPreloadUser(DB, preload...)
	}
	return DB.First(&u).Error
}
func NewPreloadUser(DB *gorm.DB, properties ...string) {
	for _, v := range properties {
		if v == "Subjects" {
			DB.Preload("Subjects", func(db *gorm.DB) *gorm.DB {
				return db.Debug().Select("subjects.id", "subjects.name", "subjects.code").Joins("JOIN (SELECT s2.code, s2.name, MAX(updated_at) AS latest_updated_at FROM subjects as s2 WHERE s2.deleted_at IS NULL GROUP BY s2.code, s2.name) AS latest_subjects ON (subjects.code = latest_subjects.code OR subjects.name = latest_subjects.name) AND subjects.updated_at = latest_subjects.latest_updated_at")
			})
		}
		//Truy vấn con (latest_subjects) lấy danh sách môn học có updated_at mới nhất.
		//	JOIN với bảng subjects giúp lấy đúng bản ghi tương ứng, có thể là 1 hoặc nhiều bản ghi tùy vào dữ liệu trong bảng gốc.
	}
}
func CreateUser(entry *models.User, args map[string]interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	tx := app.Database.DB.WithContext(ctx).Begin()

	if err = tx.Create(entry).Error; err != nil {
		tx.Rollback()
		return err
	}
	pwd := args["password"].(string)
	temp, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	loginInfo := LoginInfo{
		ID:           entry.ID,
		CenterID:     *entry.CenterId,
		Username:     entry.Username,
		Phone:        entry.Phone,
		Email:        entry.Email,
		RoleId:       entry.RoleId,
		PasswordHash: string(temp),
		DeletedAt:    nil,
	}
	if err = loginInfo.Create(); err != nil {
		tx.Rollback()
		return err
	}
	if entry.Position == consts.Teacher || entry.Position == consts.TeachingAssistant {
		history := SalaryHistory{
			UserID:     entry.ID,
			Salary:     entry.Salary,
			SalaryType: entry.SalaryType,
			CenterID:   entry.CenterId,
			BranchID:   entry.BranchId,
			OrganID:    entry.OrganStructId,
		}
		now := time.Now()
		history.CreatedAt = &now
		if err = tx.Create(&history).Error; err != nil {
			logrus.Error("Error while create salary history: ", err)
		}
	}
	return tx.Commit().Error
}
func FindUsers(DB *gorm.DB) ([]models.User, error) {
	var (
		entries     []models.User
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := DB.WithContext(ctx).Find(&entries)
	return entries, err.Error
}
func PreloadUser(entry *models.User, properties ...string) {
	for _, v := range properties {
		if v == "organStructName" {
			var (
				err         error
				organStruct models.OrganStruct
			)
			if organStruct, err = FirstOrganStruct(app.Database.DB.Where("id = ?", entry.OrganStructId).Select("name", "parent_id")); err != nil {
				logrus.Error(err)
			} else {
				entry.OrganStructName = organStruct.Name
			}

		}
	}
}
func CountUser(DB *gorm.DB) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	DB.Model(&models.User{}).WithContext(ctx).Count(&count)
	return count
}

func FirstUser(query interface{}, args []interface{}, preload ...string) (models.User, error) {
	var (
		entry       models.User
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = app.Database.DB.WithContext(ctx).Where(query, args...)
	)
	defer cancel()
	if len(preload) > 0 {
		NewPreloadUser(DB, preload...)
	}
	err := DB.First(&entry)
	return entry, err.Error
}

func UpdateUser(entry *models.User, origin models.User, query interface{}, args []interface{}) (err error) {
	var (
		ctx, cancel     = context.WithTimeout(context.Background(), app.CTimeOut)
		tx              = app.Database.DB.WithContext(ctx).Begin()
		teacherPosition = []int64{consts.Teacher, consts.TeachingAssistant}
	)
	defer cancel()
	if origin.Salary != entry.Salary || origin.SalaryType != entry.SalaryType {
		if utils.Contains(teacherPosition, entry.Position) && utils.Contains(teacherPosition, origin.Position) {
			var VietNamTZ *time.Location
			VietNamTZ, err = time.LoadLocation("Asia/Ho_Chi_Minh")
			if err != nil {
				logrus.Error(err)
				tx.Rollback()
				return err
			}
			next24hour := time.Now().AddDate(0, 0, 1)
			tomorrow := time.Date(next24hour.Year(), next24hour.Month(), next24hour.Day(), 0, 5, 0, 0, VietNamTZ)
			jobrunner.In(tomorrow.Sub(time.Now()), UpdateSalary{
				UserId:     entry.ID,
				CenterId:   entry.CenterId,
				BranchId:   entry.BranchId,
				OrganID:    entry.OrganStructId,
				Salary:     entry.Salary,
				SalaryType: entry.SalaryType,
			})
		}
	} else if utils.Contains(teacherPosition, entry.Position) {
		history := SalaryHistory{
			UserID:     entry.ID,
			SalaryType: entry.SalaryType,
			Salary:     entry.Salary,
			CenterID:   entry.CenterId,
			BranchID:   entry.BranchId,
			OrganID:    entry.OrganStructId,
		}
		now := time.Now()
		history.CreatedAt = &now
		if err = tx.Create(&history).Error; err != nil {
			logrus.Error("Error while create salary history: ", err)
			tx.Rollback()
		}
	}
	if entry.OrganStructId != origin.OrganStructId && entry.OrganStructId == nil {
		if err = tx.Where(query, args...).Model(&models.User{}).Update("organ_struct_id", nil).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
	}
	if err = tx.Where(query, args...).Updates(&entry).Error; err != nil {
		logrus.Error(err)
		tx.Rollback()
		return err
	}
	loginInfo := LoginInfo{
		Username: entry.Username,
		Phone:    entry.Phone,
		Email:    entry.Email,
	}
	if err = loginInfo.Update("id = ?", []interface{}{entry.ID}); err != nil {
		logrus.Error(err)
		tx.Rollback()
		return err
	}
	if err = tx.Model(entry).Association("Subjects").Replace(entry.Subjects); err != nil {
		logrus.Error(err)
		tx.Rollback()
		return err
	}
	if entry.BranchId == nil {
		if err = tx.Model(&models.User{}).Where(query, args...).Update("BranchId", nil).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
	}

	if !entry.EmailVerified {
		if err = tx.Model(&models.User{}).Where(query, args...).
			Update("EmailVerified", false).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
	}
	// Nếu xóa vị trí => xóa lịch dạy của người dùng
	if utils.Contains(teacherPosition, entry.Position) && !utils.Contains(teacherPosition, entry.Position) {
		if err = app.Database.DB.WithContext(ctx).Where("user_id = ?", entry.ID).
			Delete(&models.TeachingSchedule{}).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}

		if err = TsDeleteTimeSlot(app.Database.DB.WithContext(ctx).Where("user_id = ?", entry.ID)); err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}

		if err = TsDeleteShift(app.Database.DB.WithContext(ctx).Where("user_id = ?", entry.ID)); err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func DeleteUser(query interface{}, args []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.CTimeOut)
	defer cancel()
	return app.Database.DB.WithContext(ctx).Where(query, args...).Delete(&models.User{}).Error
}

func (u *User) Delete(teacherIds, otherPosIds []uuid.UUID) (err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = app.Database.DB.WithContext(ctx)
	)
	defer cancel()

	if len(teacherIds) < 1 {
		tx := DB.Begin()
		if err = tx.Error; err != nil {
			logrus.Error(err)
			return err
		}

		if err = tx.Where("id IN (?)", otherPosIds).Delete(&models.User{}).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
		if err = tx.Where("id IN (?)", otherPosIds).Delete(&models.LoginInfo{}).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}

		return tx.Commit().Error
	} else {
		tx := DB.Begin()
		if err = tx.Error; err != nil {
			logrus.Error(err)
			return err
		}

		if err = tx.Where("id IN (?)", append(otherPosIds, teacherIds...)).Delete(&models.User{}).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
		if err = tx.Where("id IN (?)", append(otherPosIds, teacherIds...)).Delete(&models.LoginInfo{}).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
		if err = tx.Where("user_id IN (?)", teacherIds).Delete(&models.TeachingSchedule{}).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
		if err = tx.Where("user_id IN (?)", teacherIds).Delete(&models.TimeSlot{}).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
		if err = tx.Where("user_id IN (?)", teacherIds).Delete(&models.Shift{}).Error; err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}

		return tx.Commit().Error
	}
}

func NewFindUsers(query interface{}, args []interface{}) ([]models.User, error) {
	var (
		entries     []models.User
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	err := app.Database.DB.WithContext(ctx).Where(query, args...).Find(&entries)
	return entries, err.Error
}
