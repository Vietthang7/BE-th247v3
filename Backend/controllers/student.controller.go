package controllers

import (
	"errors"
	"fmt"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateStudent(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error unauthorized!", consts.ERROR_UNAUTHORIZED)
	}
	var (
		err   error
		entry repo.Student
	)
	if err = c.BodyParser(&entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
	}
	// data validation
	if entry.Type == consts.Official || entry.Type == consts.Trial {
		var existence repo.LoginInfo
		if err = existence.First("email = ? OR phone = ?", []interface{}{entry.Email, entry.Phone}); err == nil {
			var errExist = ""
			if entry.Email != "" && existence.Email == entry.Email {
				errExist += "Email đã tồn tại."
				return ResponseError(c, fiber.StatusConflict, errExist, consts.EmailDuplication)
			}
			if entry.Phone != "" && existence.Phone == entry.Phone {
				errExist += "Số điện thoại đã tồn tại."
				return ResponseError(c, fiber.StatusConflict, errExist, consts.PhoneDuplication)
			}
		}
	} else {
		var (
			existence repo.Student
			errExist  = ""
		)
		if err = existence.First("email = ? OR phone = ?", []interface{}{entry.Email, entry.Phone}); err == nil {
			if entry.Email != "" && existence.Email == entry.Email {
				errExist += "Email đã tồn tại"
				return ResponseError(c, fiber.StatusConflict, errExist, consts.EmailDuplication)
			}
			if entry.Phone != "" && existence.Phone == entry.Phone {
				errExist += "Số điện thoại đã tồn tại"
				return ResponseError(c, fiber.StatusConflict, errExist, consts.PhoneDuplication)
			}
		}
		var userExistence repo.User
		if err = userExistence.First("email = ? OR phone = ?", []interface{}{entry.Email, entry.Phone}); err == nil {
			if entry.Email != "" && existence.Email == entry.Email {
				errExist += "Email đã tồn tại."
				return ResponseError(c, fiber.StatusConflict, errExist, consts.EmailDuplication)
			}
			if entry.Phone != "" && existence.Phone == entry.Phone {
				errExist += "Số điện thoại đã tồn tại."
				return ResponseError(c, fiber.StatusConflict, errExist, consts.PhoneDuplication)
			}
		}
	}
	entry.CenterId = *user.CenterId
	entry.BranchId = user.BranchId
	if entry.Type == consts.Official {
		entry.IsOfficialAt = time.Now()
	}
	if err = entry.Create(); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.CreateFail, err.Error()), consts.CreateFailed)
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry.ID)
}
func ReadStudent(c *fiber.Ctx) error {
	_, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, consts.InvalidInput, "Permission denied!")
	}
	var (
		err   error
		entry repo.Student
	)
	if err = entry.First("id = ?", []interface{}{c.Params("id")}, "Province", "District", "CustomerSource", "ContactChannel"); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusForbidden, consts.GetFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		//"can_not_update": canNotUpdate,
		"entry": entry,
	})
}

func UpdateStudent(c *fiber.Ctx) error {
	_, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Permission denied!", consts.Forbidden)
	}
	var (
		err   error
		entry repo.Student
	)

	err = entry.First("id", []interface{}{c.Params("id")})
	switch {
	case err == nil:
		origin := entry
		if err = c.BodyParser(&entry); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusBadRequest,
				fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.InvalidReqInput)
		}

		// data validation
		if entry.Type == consts.Official || entry.Type == consts.Trial {
			var existence repo.LoginInfo
			if err = existence.First("id <> ? AND (email = ? OR phone = ?)",
				[]interface{}{entry.ID, entry.Email, entry.Phone}); err == nil {
				var (
					errExist     = ""
					errExistCode []int
				)
				if entry.Email != "" && existence.Email == entry.Email {
					errExist += "Email đã tồn tại."
					errExistCode = append(errExistCode, consts.EmailDuplication)
				}
				if entry.Phone != "" && existence.Phone == entry.Phone {
					errExist += "Số điện thoại đã tồn tại."
					errExistCode = append(errExistCode, consts.PhoneDuplication)
				}

				return ResponseError(c, fiber.StatusConflict, errExist, errExistCode)
			}
		} else {
			var (
				existence repo.Student
				errExist  = ""
			)
			if err = existence.First("id <> ? AND (email = ? OR phone = ?)",
				[]interface{}{entry.ID, entry.Email, entry.Phone}); err == nil {
				if entry.Email != "" && existence.Email == entry.Email {
					errExist += "Email đã tồn tại."
					return ResponseError(c, fiber.StatusConflict, errExist, consts.EmailDuplication)
				}
				if entry.Phone != "" && existence.Phone == entry.Phone {
					errExist += "Số điện thoại đã tồn tại."
					return ResponseError(c, fiber.StatusConflict, errExist, consts.PhoneDuplication)
				}
			}

			var userExistence repo.User
			if err = userExistence.First("id <> ? AND (email = ? OR phone = ?)",
				[]interface{}{entry.ID, entry.Email, entry.Phone}); err == nil {
				if entry.Email != "" && userExistence.Email == entry.Email {
					errExist += "Email đã tồn tại."
					return ResponseError(c, fiber.StatusConflict, errExist, consts.EmailDuplication)
				}
				if entry.Phone != "" && userExistence.Phone == entry.Phone {
					errExist += "Số điện thoại đã tồn tại."
					return ResponseError(c, fiber.StatusConflict, errExist, consts.PhoneDuplication)
				}
			}
		}

		if origin.Type != consts.Official && entry.Type == consts.Official {
			entry.IsOfficialAt = time.Now()
		}

		if err = entry.Update(origin, "id", []interface{}{c.Params("id")}); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError,
				fmt.Sprintf("%s: %s", consts.UpdateFail, err.Error()), consts.UpdateFailed)
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, nil)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound,
			fmt.Sprintf("%s: %s", consts.NotFound, err.Error()), consts.GetFailed)
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
}

type Status struct {
	All      int64 `json:"all"`
	Studying int64 `json:"studying"`
	Pending  int64 `json:"pending"`
	Reserved int64 `json:"reserved"`
}

func ListStudents(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, "Error Unauthorized!", consts.ERROR_UNAUTHORIZED)
	}

	var (
		err                error
		entry              repo.Student
		entries            repo.Students
		pagination         = consts.BindRequestTable(c, "created_at")
		DB                 = app.Database.DB.Where("students.center_id = ?", user.CenterId)
		startDate, endDate time.Time
		startTimeSearch    = false
		endTimeSearch      = false
		status             Status
		preload            = []string{"StudyNeeds"}
		studentType        int
	)

	if c.Query("type") != "" {
		if studentType, err = strconv.Atoi(c.Query("type")); err != nil {
			logrus.Error("Error 'type' param: ", err.Error())
			return ResponseError(c, fiber.StatusBadRequest, "Error 'type' param: "+err.Error(), consts.InvalidReqInput)
		}

		switch studentType {
		case consts.Trial:
			if user.RoleId == consts.CenterHR && user.Position != consts.CareAssignee {
				DB = DB.Where("caregiver_id = ?", user.ID)
			}
		case consts.Potential:
			preload = append(preload, "Source")
			if user.RoleId == consts.CenterHR && user.Position != consts.CareAssignee {
				DB = DB.Where("caregiver_id = ?", user.ID)
			} else {
				if c.Query("assigned") != "" {
					var assgined bool
					if assgined, err = strconv.ParseBool(c.Query("assigned")); err != nil {
						logrus.Error("Error 'assigned' query: ", err.Error())
						return ResponseError(c, fiber.StatusBadRequest,
							"Error 'assigned' query: "+err.Error(), consts.InvalidReqInput)
					}
					if assgined {
						DB = DB.Where("caregiver_id IS NOT NULL")
					} else {
						DB = DB.Where("caregiver_id IS NULL")
					}
				}
			}

			if c.Query("source") != "" {
				DB = DB.Where("customer_source_id = ?", c.Query("source"))
			}
		}

		DB = DB.Where("type = ?", c.Query("type"))
	}

	preload = append(preload, "Caregiver")

	if c.Query("start_at") != "" {
		startTimeSearch = true
	}
	if c.Query("end_at") != "" {
		endTimeSearch = true
	}
	if c.Query("enroll_id") != "" {
		DB = DB.Joins("JOIN study_needs sn ON sn.student_id = students.id").
			Where("sn.enrollment_id = ?", c.Query("enroll_id"))
	}

	if endTimeSearch || startTimeSearch {
		if startTimeSearch {
			startDate, _ = time.Parse("2006-01-02", c.Query("start_at"))
			if !endTimeSearch {
				endDate = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
					23, 59, 59, 0, time.UTC)
			} else {
				endDate, _ = time.Parse("2006-01-02", c.Query("end_at"))
			}
		} else {
			startDate = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
				0, 0, 0, 0, time.UTC)
			endDate, _ = time.Parse("2006-01-02", c.Query("end_at"))
		}

		switch studentType {
		case consts.Official:
			DB = DB.Where("is_official_at BETWEEN ? AND ?", startDate, endDate)
		case consts.Trial:
			DB = DB.Where("is_trial_at BETWEEN ? AND ?", startDate, endDate)
		default:
			DB = DB.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}
	}

	if user.BranchId != nil {
		DB = DB.Where("branch_id = ?", *user.BranchId)
	}

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		DB = DB.Where("full_name LIKE ? OR phone LIKE ?", search, search)
	}

	status.All = entry.Count(DB.Session(&gorm.Session{}))
	status.Studying = entry.Count(DB.Session(&gorm.Session{}).Where("status = ?", consts.Studying))

	status.Pending = entry.Count(DB.Session(&gorm.Session{}).Where("status = ?", consts.Pending))

	status.Reserved = entry.Count(DB.Session(&gorm.Session{}).Where("status = ?", consts.Reserved))
	DB = pagination.CustomOptions(DB)
	if c.Query("status") != "" {
		status, _ := strconv.Atoi(c.Query("status"))
		DB = DB.Where("status = ?", status)
	}

	if entries, err = entry.Find(DB, preload...); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}

	switch studentType {
	case consts.Trial:
		for i := range entries {
			repo.PreloadTotalTrialSession(entries[i])
		}
	case consts.Potential:
		for i := range entries {
			entries[i].CareResult = repo.LoadCareResult(entries[i].ID)
		}
	case consts.Official:
		for i := range entries {
			entries[i].CompletedSubject = entry.PreloadCompletedSubject(entries[i].ID)
		}
	}

	pagination.Total = entry.Count(DB.Offset(-1))
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
		"status":     status,
	})
}
