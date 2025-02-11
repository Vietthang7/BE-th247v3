package controllers

import (
	"errors"
	"fmt"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/helpers"
	"intern_247/models"
	"intern_247/repo"
	"intern_247/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewLoginGetToken(c *fiber.Ctx) error {
	type Input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var (
		input Input
		err   error
	)
	if err := c.BodyParser(&input); err != nil {
		return c.JSON(fiber.Map{"status": false, "message": "Review your input", "error": err.Error()})
	}
	var loginInfo repo.LoginInfo
	if err = loginInfo.First("email = ? OR phone = ? OR username = ?",
		[]interface{}{input.Email, input.Email, input.Email}, "Student", "User"); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFailed, err.Error()), consts.GetFailed)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(loginInfo.PasswordHash), []byte(input.Password)); err != nil {
		return c.JSON(fiber.Map{"status": false, "message": "Review your input password", "error": err, "user": nil})
	}
	verifiedEmail := false
	if loginInfo.RoleId == consts.Student {
		verifiedEmail = loginInfo.Student.EmailVerified
	} else {
		verifiedEmail = loginInfo.User.EmailVerified
	}
	if !utils.IsVerifiedEmail(&verifiedEmail) {
		// send mail
		var emailInfo helpers.EmailSchema
		emailInfo.Title = helpers.EmailVerifyTitle
		emailInfo.Content = helpers.EmailVerifyContent
		emailInfo.Receivers = loginInfo.Email
		emailInfo.Sender = app.Config("SEND_MAIL_EMAIL")
		emailInfo.Provider = app.Config("PROVIDER")
		emailInfo.Password = app.Config("SEND_MAIL_PASSWORD")
		emailInfo.Port = app.Config("PORT_EMAIL")
		code, err := helpers.SendEmailOTP(emailInfo)
		if err != nil {
			return ResponseError(c, fiber.StatusInternalServerError, "Error send OTP!", consts.SendOTPFailed)
		}
		var otpInfo models.OTPLog
		otpInfo.Code = code
		otpInfo.CreatedBy = loginInfo.ID
		otpInfo.Receiver = loginInfo.Email
		otpInfo.ExpiredAt = time.Now().Add(time.Minute * 5)
		err3 := repo.CreateOTPLog(&otpInfo)
		if err3 != nil {
			return ResponseError(c, fiber.StatusInternalServerError, "Error OTP log!", consts.LoginFailed)
		}
		return ResponseError(c, fiber.StatusForbidden, loginInfo.Email, consts.EmailIsNotVerified)
	}
	roleData := "owner"
	if loginInfo.RoleId == consts.Student {
		roleData = "student"
	} else if loginInfo.User.Position == consts.Teacher {
		roleData = "teacher"
	}
	claims_new := jwt.MapClaims{
		"user_id": loginInfo.ID,
		"role_id": loginInfo.RoleId,
		"site_id": loginInfo.CenterID,
		//"sso_id":  &loginInfo.SsoID,
		"role":   roleData,
		"status": true,
		"exp":    time.Now().Add(time.Hour * 200).Unix(),
	}
	claims_refresh := jwt.MapClaims{
		"user_id": loginInfo.ID,
		"role_id": loginInfo.RoleId,
		"site_id": loginInfo.CenterID,
		"role":    roleData,
		//"sso_id":  &loginInfo.SsoID,
		"status": true,
		"exp":    time.Now().Add(time.Hour * 168).Unix(),
	}
	token_new := jwt.NewWithClaims(jwt.SigningMethodHS256, claims_new)
	t, errs := token_new.SignedString([]byte(app.Config("SECRET_KEY")))
	if errs != nil {
		return c.JSON(fiber.Map{"status": false, "message": "Token generation failed", "error": errs.Error(), "user": nil})
	}
	token_refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, claims_refresh)
	refresh_token, errs1 := token_refresh.SignedString([]byte(app.Config("SECRET_KEY")))
	if errs1 != nil {
		return c.JSON(fiber.Map{"status": false, "message": "Token generation failed", "error": errs1.Error(), "user": nil})
	}
	userReturn := models.DataUserReturn{
		ID:           loginInfo.ID,
		RoleId:       loginInfo.RoleId,
		Email:        loginInfo.Email,
		Phone:        loginInfo.Phone,
		Token:        t,
		RefreshToken: refresh_token,
	}
	return ResponseSuccess(c, fiber.StatusOK, "Login success", userReturn)

}

func VerifyEmailOTP(c *fiber.Ctx) error {
	type VerifySchema struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	var input VerifySchema
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Failed 1", consts.InvalidReqInput)
	}
	otpInfo, err := repo.GetNewestOTPLogByReceiver(input.Email)
	if err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadGateway, "Failed 2", consts.GetFailed)
	}
	if strings.Compare(otpInfo.Code, input.Code) != 0 {
		logrus.Println("Invalid OTP!")
		return ResponseError(c, fiber.StatusBadRequest, "Failed 2", consts.InvalidReqInput)
	}
	if time.Now().After(otpInfo.ExpiredAt) {
		logrus.Println("Expired OTP!")
		return ResponseError(c, fiber.StatusBadRequest, "Failed 3", consts.ERROR_EXPIRED_TIME)
	}
	if otpInfo.IsConfirmed {
		logrus.Println("OTP Confirmed")
		return ResponseError(c, fiber.StatusBadRequest, "Failed 4", consts.InvalidReqInput)
	}
	otpInfo.IsConfirmed = true
	err2 := repo.UpdateOTPLogById(&otpInfo)
	if err2 != nil {
		logrus.Error(err2)
		return ResponseError(c, fiber.StatusInternalServerError, "Failed 5", consts.UpdateFailed)
	}

	var loginInfo repo.LoginInfo
	if err = loginInfo.First("email = ?", []interface{}{input.Email}, "User", "Student"); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, "Failed 6", consts.GetFailed)
	}

	if loginInfo.User == nil && loginInfo.Student == nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed 7", consts.GetFailed)
	}

	var (
		fullName, avatar string
		position         int64
		permissionGrpID  uuid.UUID
	)
	if loginInfo.RoleId == consts.Student {
		if loginInfo.Student.EmailVerified {
			return ResponseError(c, fiber.StatusBadRequest, "Failed", consts.InvalidReqInput)
		}

		var student repo.Student
		fullName = loginInfo.Student.FullName
		avatar = loginInfo.Student.Avatar

		if err = student.VerifyEmail(input.Email); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, "Failed", consts.UpdateFailed)
		}
	} else {
		if loginInfo.User.EmailVerified {
			return ResponseError(c, fiber.StatusBadRequest, "Failed 8", consts.InvalidReqInput)
		}

		if err = repo.VerifyUserEmail(input.Email); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, "Failed", consts.UpdateFailed)
		}

		if !utils.IsActiveData(loginInfo.User.IsActive) {
			return c.JSON(fiber.Map{"status": false, "message": "Login failed", "error": consts.UserIsInactive})
		}

		if loginInfo.User.PermissionGrpId != nil &&
			!utils.Contains([]int64{consts.Root, consts.CenterOwner, consts.Student}, loginInfo.User.RoleId) {
			permissionGrpID = *loginInfo.User.PermissionGrpId
		}

		fullName = loginInfo.User.FullName
		avatar = loginInfo.User.Avatar
		position = loginInfo.User.Position
	}
	roleData := "owner"
	if loginInfo.RoleId == consts.Student {
		roleData = "student"
	} else if loginInfo.User.Position == consts.Teacher {
		roleData = "teacher"
	}
	newClaims := jwt.MapClaims{
		"user_id":   loginInfo.ID,
		"role_id":   loginInfo.RoleId,
		"full_name": fullName,
		"role":      roleData,
		"site_id":   loginInfo.CenterID,
		//"sso_id":    loginInfo.SsoID,
		"status": true,
		"exp":    time.Now().Add(time.Hour * 200).Unix(),
	}
	refreshClaims := jwt.MapClaims{
		"user_id":   loginInfo.ID,
		"role_id":   loginInfo.RoleId,
		"full_name": fullName,
		"role":      roleData,
		"site_id":   loginInfo.CenterID,
		//"sso_id":    loginInfo.SsoID,
		"status": true,
		"exp":    time.Now().Add(time.Hour * 168).Unix(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	t, errs := newToken.SignedString([]byte(app.Config("SECRET_KEY")))
	if errs != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed errs", consts.UpdateFailed)
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refresh_token, errs1 := refreshToken.SignedString([]byte(app.Config("SECRET_KEY")))
	if errs1 != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed errs", consts.UpdateFailed)
	}
	userReturn := models.DataUserReturn{
		ID:             loginInfo.ID,
		FullName:       fullName,
		Avatar:         avatar,
		FirstTimeLogin: true,
		RoleId:         loginInfo.RoleId,
		Position:       position,
		Email:          loginInfo.Email,
		Phone:          loginInfo.Phone,
		//SSO_ID:         &loginInfo.SsoID,
		Token:        t,
		RefreshToken: refresh_token,
	}

	if permissionGrpID != uuid.Nil {
		var permissionGrp models.PermissionGroup
		if permissionGrp, err = repo.FirstPermissionGrp(app.Database.DB.Where("id = ?", permissionGrpID)); err == nil {
			userReturn.PermissionGrp = &permissionGrp
		}
	}

	return ResponseSuccess(c, fiber.StatusOK, "Confirmed! Welcome", userReturn)
}
func Register(c *fiber.Ctx) error {
	var (
		err  error
		form models.CreateUserForm
	)
	if err = c.BodyParser(&form); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}
	form.Username = utils.GenerateUniqueUsername()
	// data validation
	var existence repo.LoginInfo
	if err = existence.First("email = ? OR username = ?", []interface{}{form.Email, form.Username}); err == nil {
		var (
			errExist     = ""
			errExistCode []int
		)
		if existence.Email == form.Email {
			errExist += "Email existed"
			return ResponseError(c, fiber.StatusBadRequest, errExist, consts.EmailDuplication)
		}
		if existence.Username == form.Username {
			errExist += "Username existed"
			return ResponseError(c, fiber.StatusBadRequest, errExist, consts.UsernameDuplication)
		}
		return ResponseError(c, fiber.StatusBadRequest, errExist, errExistCode)
	}
	form.Username = utils.GenerateUniqueUsername()
	// Create user
	isActive := true
	entry := models.User{
		FullName: form.FullName,
		Username: form.Username,
		Email:    form.Email,
		IsActive: &isActive,
		RoleId:   consts.CenterOwner,
	}
	args := map[string]interface{}{"password": form.Password}
	if err = repo.RegisterUser(&entry, args); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.RegisterFailed)
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry.ID)
}

func ResendOTP(c *fiber.Ctx) error {
	fmt.Println("ok")
	type EmailInput struct {
		Email string `json:"email"`
	}
	var (
		input EmailInput
		err   error
	)
	if err = c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Failed", consts.InvalidReqInput)
	}
	ok := utils.EmailValid(input.Email)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, "Failed2", consts.InvalidReqInput)
	}
	var loginInfo repo.LoginInfo
	if err = loginInfo.First("email = ?", []interface{}{input.Email}); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, "Failed3", consts.GetFailed)
	}
	//if _, err = repo.GetNewestOTPLogByReceiver(input.Email); err != nil {
	//	return ResponseError(c, fiber.StatusBadRequest, "Failed4", consts.DataNotFound)
	//}
	var emailInfo helpers.EmailSchema
	emailInfo.Title = helpers.EmailVerifyTitle
	emailInfo.Content = helpers.EmailVerifyContent
	emailInfo.Receivers = loginInfo.Email
	emailInfo.Sender = app.Config("SEND_MAIL_EMAIL")
	emailInfo.Provider = app.Config("PROVIDER")
	emailInfo.Password = app.Config("SEND_MAIL_PASSWORD")
	emailInfo.Port = app.Config("PORT_EMAIL")
	code, err := helpers.SendEmailOTP(emailInfo)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "verify", consts.SendOTPFailed)
	}
	var otpInfo models.OTPLog
	otpInfo.Code = code
	otpInfo.CreatedBy = loginInfo.ID
	otpInfo.Receiver = loginInfo.Email
	otpInfo.ExpiredAt = time.Now().Add(time.Minute * 5)
	err3 := repo.CreateOTPLog(&otpInfo)
	if err3 != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "verify", consts.CreateFailed)
	}
	return ResponseSuccess(c, fiber.StatusOK, "success", loginInfo.Email)
}
func CreateUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error permission denied!", consts.Forbidden)
	}
	var (
		err  error
		form models.CreateUserForm
	)
	if err = c.BodyParser(&form); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}
	// data validation
	var existence repo.LoginInfo
	if err = existence.First("email = ? OR phone = ? OR username = ?", []interface{}{form.Email, form.Phone, form.Username}); err == nil {
		var (
			errExist     = ""
			errExistCode []int
		)
		if form.Email != "" && existence.Email == form.Email {
			errExist += "Email existed"
			errExistCode = append(errExistCode, consts.EmailDuplication)
		}
		if form.Phone != "" && existence.Phone == form.Phone {
			errExist += "Số điện thoại đã tồn tại."
			errExistCode = append(errExistCode, consts.PhoneDuplication)
		}
		if form.Username != "" && existence.Username == form.Username {
			errExist += "Tên tài khoản đã tồn tại."
			errExistCode = append(errExistCode, consts.UsernameDuplication)
		}
	}
	var student repo.Student
	if err = student.First("type NOT IN ? AND (email = ? OR phone = ? OR username = ?)",
		[]interface{}{[]int64{consts.Official, consts.Trial}, form.Email, form.Phone, form.Username}); err == nil {
		var (
			errExist     = ""
			errExistCode []int
		)
		if student.Email == form.Email {
			errExist += "Email đã tồn tại."
			errExistCode = append(errExistCode, consts.EmailDuplication)
		}
		if student.Phone == form.Phone {
			errExist += "Số điện thoại đã tồn tại."
			errExistCode = append(errExistCode, consts.PhoneDuplication)
		}
		if student.Username == form.Username {
			errExist += "Tên tài khoản đã tồn tại."
			errExistCode = append(errExistCode, consts.UsernameDuplication)
		}
		return ResponseError(c, fiber.StatusConflict, errExist, errExistCode)
	}
	// Create user
	entry := models.User{
		Avatar:          form.Avatar,
		FullName:        form.FullName,
		Username:        form.Username,
		Email:           form.Email,
		Phone:           form.Phone,
		Position:        form.Position,
		BranchId:        form.BranchId,
		OrganStructId:   form.OrganStructId,
		PermissionGrpId: form.PermissionGrpId,
		Introduction:    form.Introduction,
		IsActive:        form.IsActive,
		Salary:          form.Salary,
		SalaryType:      form.SalaryType,
		CenterId:        user.CenterId,
		RoleId:          consts.CenterHR,
	}
	if len(form.SubjectIds) > 0 {
		if entry.Subjects, err = repo.GetSubjectByIdsAndCenterId(form.SubjectIds, *user.CenterId, nil); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.CreateFailed)
		}
	}
	args := map[string]interface{}{"password": form.Password}
	if err = repo.CreateUser(&entry, args); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.CreateFailed)
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry.ID)
}
func ListUsers(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied!", consts.Forbidden)
	}
	var (
		err        error
		entries    []models.User
		pagination = consts.BindRequestTable(c, "created_at")
		DB         = pagination.CustomOptions(app.Database.DB)
	)
	DB = DB.Where(map[string]interface{}{
		"center_id": user.CenterId,
		"role_id":   consts.CenterHR,
	}).Omit("phone", "email", "updated_at", "role_id").Order("full_name asc")
	if user.BranchId != nil && c.Query("allBranch") != "true" {
		DB = DB.Where("branch_id = ?", user.BranchId)
	}
	if pagination.Search != "" {
		DB = DB.Where("fullname LIKE ?", "%"+pagination.Search+"%")
	}
	if c.Query("organ_struct") != "" {
		if id, errId := uuid.Parse(c.Query("organ_struct")); errId != nil {
			logrus.Error("Failed to parse organ_struct UUID:", errId)
		} else {
			DB = DB.Where("organ_struct_id = ?", id)
		}
	}
	if c.Query("branch") != "" {
		if id, errId := uuid.Parse(c.Query("branch")); errId != nil {
			logrus.Error("Failed to parse branch UUID:", errId)
		} else {
			DB = DB.Where("branch_id = ?", id)
		}
	}
	if c.Query("permission_grp") != "" {
		if id, errId := uuid.Parse(c.Query("permission_grp")); errId != nil {
			logrus.Error("Failed to parse permission_grp UUID:", errId)
		} else {
			DB = DB.Where("permission_grp_id = ?", id)
		}
	}
	if c.Query("active") != "" {
		isActive, _ := strconv.ParseBool(c.Query("active")) //Chuyển giá trị active từ chuỗi ("true"/"false") thành bool (true/false).
		DB = DB.Where("is_active = ?", isActive)
	}
	if c.Query("position") != "" {
		position, _ := strconv.Atoi(c.Query("position")) //Chuyển giá trị position từ chuỗi thành int.
		DB = DB.Where("position = ?", position)
	}
	if entries, err = repo.FindUsers(DB); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.GetFailed)
	}
	for i := range entries {
		repo.PreloadUser(&entries[i], "organStructName", "permissionGrpName", "branchName")
	}
	pagination.Total = repo.CountUser(DB.Where("center_id = ?", user.CenterId).Offset(-1))
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}

func ReadUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", consts.Forbidden)
	}
	entry, err := repo.FirstUser("id = ? AND id <> ?",
		[]interface{}{c.Params("id"), user.ID}, "Subjects") // Truy xuất user theo id, nhưng không cho phép user tự xem thông tin của chính họ.
	switch {
	case err == nil:
		repo.PreloadUser(&entry, "organStructName", "permissionGrpName", "branchName")
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, err.Error(), consts.DataNotFound)
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.GetFailed)
	}
}

func UpdateUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", consts.Forbidden)
	}
	entry, err := repo.FirstUser("id = ? AND id <> ?", []interface{}{c.Params("id"), user.ID})
	switch {
	case err == nil:
		origin := entry
		if err = c.BodyParser(&entry); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusBadRequest, err.Error(), consts.InvalidReqInput)
		}
		// data validation
		var existence repo.LoginInfo
		if err = existence.First("id <> ? AND (email = ? OR phone = ? username = ?)",
			[]interface{}{entry.ID, entry.Email, entry.Phone, entry.Username}); err == nil {
			var (
				errExist     = ""
				errExistCode []int
			)
			if entry.Email != "" && existence.Email == entry.Email {
				errExist += "Email đã tồn tại."
				errExistCode = append(errExistCode, consts.EmailDuplication)
			}
			if entry.Phone != "" && existence.Phone == entry.Phone {
				errExist += "Phone đã tồn tại."
				errExistCode = append(errExistCode, consts.PhoneDuplication)
			}
			if entry.Username != "" && existence.Username == entry.Username {
				errExist += "Tên tài khoản đã tồn tại."
				errExistCode = append(errExistCode, consts.UsernameDuplication)
			}
			return ResponseError(c, fiber.StatusConflict, errExist, errExistCode)
		}
		var student repo.Student
		if err = student.First("id <> ? AND type NOT IN ? AND (email = ? OR phone = ? OR username = ?)",
			[]interface{}{entry.ID, []int64{consts.Official, consts.Trial}, entry.Email, entry.Phone, entry.Username}); err == nil {
			var (
				errExist     = ""
				errExistCode []int
			)
			if entry.Email != "" && student.Email == entry.Email {
				errExist += "Email đã tồn tại."
				errExistCode = append(errExistCode, consts.EmailDuplication)
			}
			if entry.Phone != "" && student.Phone == entry.Phone {
				errExist += "Phone đã tồn tại."
				errExistCode = append(errExistCode, consts.PhoneDuplication)
			}
			if entry.Username != "" && student.Username == entry.Username {
				errExist += "Tên tài khoản đã tồn tại."
				errExistCode = append(errExistCode, consts.UsernameDuplication)
			}
			return ResponseError(c, fiber.StatusConflict, errExist, errExistCode)
		}
		if origin.Position == consts.Teacher || origin.Position == consts.TeachingAssistant {
			if entry.Position != origin.Position && repo.TeacherIsArranged(origin.ID) {
				return ResponseError(c, fiber.StatusBadRequest,
					"Nhân sự đã được gán dữ liệu. Không thể chỉnh sửa.", consts.UpdateUserHasDataLinked)
			}
		}
		entry.Username = origin.Username // không cho phép đổi username
		if entry.Subjects, err = repo.GetSubjectByIdsAndCenterId(entry.SubjectIds, *user.CenterId, nil); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
		}
		if err = repo.UpdateUser(&entry, origin, "id = ?", []interface{}{c.Params("id")}); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.UpdateFailed)
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, err.Error(), consts.DataNotFound)
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.GetFailed)
	}
}
func ForgotPwd(c *fiber.Ctx) error {
	var (
		err error
	)
	type PwdForgotForm struct {
		Email  string `json:"email"`
		NewPwd string `json:"new_pwd"`
	}
	var form PwdForgotForm
	if err = c.BodyParser(&form); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}
	if form.NewPwd == "" {
		return ResponseError(c, fiber.StatusBadRequest, "Error empty password!", consts.InvalidReqInput)
	}
	var loginInfo repo.LoginInfo
	if err = loginInfo.First("email = ?", []interface{}{form.Email}); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("Error can't find the login info: %s", err.Error()), consts.DataNotFound)
	}
	if err = loginInfo.PwdChanging("", form.NewPwd); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.UpdateFail, err.Error()), consts.UpdateFailed)
	}
	return ResponseSuccess(c, fiber.StatusOK, "Đổi mật khẩu thành công!", nil)
}

func DeleteUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err    error
		reqIds models.ReqIds
	)

	if err = c.BodyParser(&reqIds); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, err.Error(), consts.InvalidReqInput)
	}

	if len(reqIds.Ids) < 1 {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, "Error length of ids < 1", consts.InvalidReqInput)
	}

	var users []models.User
	if users, err = repo.NewFindUsers("id IN (?) AND id <> ?", []interface{}{reqIds.Ids, user.ID}); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.DeletedFailed)
	} else if len(users) < 1 {
		return ResponseError(c, fiber.StatusNotFound,
			"Can't find any users by requirement", consts.DeletedFailed)
	}

	var (
		teacherIds, otherPosIds []uuid.UUID
		entry                   repo.User
	)

	for _, v := range users {
		if v.Position == consts.Teacher || v.Position == consts.TeachingAssistant {
			if repo.TeacherIsArranged(v.ID) {
				return ResponseError(c, fiber.StatusBadRequest,
					fmt.Sprintf("Nhân sự %s đã được phân công giảng dạy. Không thể xóa.", v.FullName),
					consts.UserIsArranged)
			}
			teacherIds = append(teacherIds, v.ID)
		} else {
			otherPosIds = append(otherPosIds, v.ID)
		}
	}

	if err = entry.Delete(teacherIds, otherPosIds); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, err.Error(), consts.DeletedFailed)
	}

	// jobrunner.Now(repo.DeleteCasdoorUsers{IDs: append(teacherIds, otherPosIds...)})

	//if len(users) > 0 {
	//	for _, v := range users {
	//		if _, err = DeleteCasdoorUser(v.Email, v.Phone, v.Username); err != nil {
	//			logrus.Error(err)
	//		}
	//	}
	//}

	return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, nil)
}
