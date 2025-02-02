package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/helpers"
	"intern_247/models"
	"intern_247/repo"
	"intern_247/utils"
	"strings"
	"time"
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
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
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
		Code  string `json:"password"`
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
		logrus.Error("Invalid OTP !")
		return ResponseError(c, fiber.StatusBadRequest, "Failed 2", consts.InvalidReqInput)
	}
	if time.Now().After(otpInfo.ExpiredAt) {
		logrus.Println("Expired OTP !")
		return ResponseError(c, fiber.StatusBadRequest, "Failed 3", consts.ERROR_EXPIRED_TIME)
	}
	if otpInfo.IsConfirmed {
		logrus.Println("Confirmed OTP !")
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
			return ResponseError(c, fiber.StatusBadRequest, "Failed 8", consts.InvalidReqInput)
		}
		var student repo.Student
		fullName = loginInfo.Student.FullName
		avatar = loginInfo.Student.Avatar

		if err = student.VerifyEmail(input.Email); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, "Failed 9", consts.UpdateFailed)
		}
	} else {
		if loginInfo.User.EmailVerified {
			return ResponseError(c, fiber.StatusBadRequest, "Failed 10", consts.InvalidReqInput)
		}
		if err = repo.VerifyUserEmail(input.Email); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, "Failed 11", consts.UpdateFailed)
		}
		if !utils.IsActiveData(loginInfo.User.IsActive) {
			return c.JSON(fiber.Map{"status": false, "message": "Login failed", "error": consts.UserIsInactive})
		}
		if loginInfo.User.PermissionGrpId != nil && !utils.Contains([]int64{consts.Root, consts.Student, consts.CenterOwner}, loginInfo.User.RoleId) {
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
		//"sso_id" : loginInfo.
		"status": true,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}
	refreshClaims := jwt.MapClaims{
		"user_id":   loginInfo.ID,
		"role_id":   loginInfo.RoleId,
		"full_name": fullName,
		"role":      roleData,
		"site_id":   loginInfo.CenterID,
		//"sso_id" : loginInfo.
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
		return ResponseError(c, fiber.StatusInternalServerError, "Failed errs1", consts.UpdateFailed)
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
		Token:          t,
		RefreshToken:   refresh_token,
	}
	if permissionGrpID != uuid.Nil {
		var permissionGrp models.PermissionGroup
		if permissionGrp, err = repo.FirstPermissionGrp(app.Database.DB.Where("id = ?", permissionGrpID)); err == nil {
			userReturn.PermissionGrp = &permissionGrp
		}
	}
	return ResponseSuccess(c, fiber.StatusOK, "Confirmed! Welcome", userReturn)
}

//	func VerifyToken(c *fiber.Ctx) error {
//		type Input struct {
//			Token string `json:"token"`
//		}
//		type MyCustomClaims struct {
//			UserID string `json:"user_id"`
//			jwt.RegisteredClaims
//		}
//		var (
//			input                     Input
//			claims                    *MyCustomClaims
//			fullName, avatar          string
//			branchId                  *uuid.UUID
//			position                  int64
//			loginInfo                 repo.LoginInfo
//			userLogin, isVerifiedMail bool
//		)
//		if err := c.BodyParser(&input); err != nil {
//			return c.JSON(fiber.Map{"status": false, "message": "Review your input", "error": err.Error(), "user": nil})
//		}
//		token, err := jwt.ParseWithClaims(input.Token, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//			return []byte(app.Config("SECRET_KEY")), nil
//		})
//		if err != nil {
//			return ResponseError(c, fiber.StatusInternalServerError, "verify", consts.ERROR_INTERNAL_SERVER_ERROR)
//		}
//		ok := false
//		if claims, ok := token.Claims.(*MyCustomClaims); !ok && !token.Valid {
//			return ResponseError(c, fiber.StatusInternalServerError, "verify", consts.ERROR_INTERNAL_SERVER_ERROR)
//		}
//		if errInfo := loginInfo.First("id = ?", []interface{}{claims.UserID}, "Student", "User"); errInfo != nil {
//			return ResponseError(c,fiber.StatusInternalServerError, "verify", consts.ERROR_INTERNAL_SERVER_ERROR)
//		}
//		switch expr {
//
//		}
//
// }
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
		return ResponseError(c, fiber.StatusBadRequest, "Failed", consts.InvalidReqInput)
	}
	var loginInfo repo.LoginInfo
	if err = loginInfo.First("email = ?", []interface{}{input.Email}); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, "Failed", consts.GetFailed)
	}
	if _, err = repo.GetNewestOTPLogByReceiver(input.Email); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Failed", consts.DataNotFound)
	}
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
