package controllers

import (
	"errors"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"

	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkSessionInput struct {
	Title     string     `json:"title"`
	BranchId  *uuid.UUID `json:"branch_id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
	IsActive  *bool      `json:"is_active"`
}

type WorkSessionUpdateInput struct {
	Id        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	BranchId  *uuid.UUID `json:"branch_id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
	IsActive  *bool      `json:"is_active"`
}

func ListWorkSessionForSchedule(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Permission denied!", "Permission denied!")
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusUnauthorized, consts.InvalidInput, "Permission denied!")
	}

	var (
		err     error
		entries []models.WorkSession
		DB      = app.Database.DB
	)

	if c.Query("allBranch") == "true" {
		DB = DB.Where("center_id = ? AND is_active = ? AND branch_id IS NULL", *user.CenterId, true)
	} else {
		DB = DB.Where("(center_id = ? AND is_active = ?) AND (branch_id = ? OR branch_id IS NULL)",
			user.CenterId, true, c.Query("branchId"))
	}

	if entries, err = repo.ListWorkSessions(DB.Order("TIME(start_time)")); err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, entries)
}

func GetListWorkSessions(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}

	var (
		err        error
		entries    []models.WorkSession
		pagination = consts.BindRequestTable(c, "created_at")
		DB         = pagination.CustomOptions(app.Database.DB)
	)

	DB = DB.Where("center_id = ?", *user.CenterId)

	if user.BranchId != nil {
		DB = DB.Where("branch_id = ? OR branch_id IS NULL", user.BranchId)
	}
	if c.Query("branchId") != "" {
		DB = DB.Where("branch_id = ?", c.Query("branchId"))
	}

	if pagination.Search != "" {
		DB = DB.Where("title LIKE ?", "%"+pagination.Search+"%")
	}
	if c.Query("active") != "" {
		isActive, _ := strconv.ParseBool(c.Query("active"))
		DB = DB.Where("is_active = ?", isActive)
	}
	DB = DB.Preload("Branch", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})
	entries, err = repo.ListWorkSessions(DB)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.NotFound, "Can't find data")
	}
	pagination.Total = repo.CountWorkSession(DB.Offset(-1))
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}

func GetWorkSessionDetail(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}
	id, err := uuid.Parse(c.Query("id"))
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Invalid input!")
	}
	workSession, err := repo.GetWorkSessionByIdAndCenterId(id, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.GetFail, "Get data failed!")
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, workSession)
}

func CreateWorkSession(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}
	var input WorkSessionInput
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Please check your arg")
	}
	if _, err := repo.FindSessionByName(input.Title, *user.CenterId); err == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_WORK_SESSION_NAME_EXIST)
	}
	var newWorkSession models.WorkSession
	if input.BranchId != nil {
		_, err := repo.GetBranchByIdAndCenterId(*input.BranchId, *user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Can't find branch")
		}
		newWorkSession.BranchId = input.BranchId
	}
	if input.StartTime.After(input.EndTime) {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_START_TIME_MUST_SMALLER_THAN_END_TIME)
	}
	newWorkSession.Title = input.Title
	newWorkSession.StartTime = &input.StartTime
	newWorkSession.EndTime = &input.EndTime
	newWorkSession.IsActive = input.IsActive
	newWorkSession.CenterId = user.CenterId
	newWorkSession.UserId = user.ID
	newWorkSession.ID = uuid.New()
	row, err := repo.CreateWorkSession(&newWorkSession)
	if row < 1 || err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.CreateFail, "Create record failed")
	}
	return ResponseSuccess(c, fiber.StatusCreated, consts.CreateSuccess, newWorkSession)
}

func UpdateWorkSession(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Permission denied!")
	}
	var input WorkSessionUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Please check your arg")
	}
	workSession, err := repo.GetWorkSessionByIdAndCenterId(input.Id, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, "Can't find data")
	}
	if _, err = repo.FirstWorkSession(app.Database.DB.Where(consts.NilDeletedAt).Where("title = ?", input.Title).
		Where("id <> ? AND center_id = ?", workSession.ID, *user.CenterId)); err == nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_WORK_SESSION_NAME_EXIST)
	}
	if input.BranchId != nil {
		_, err := repo.GetBranchByIdAndCenterId(*input.BranchId, *user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.DataNotFound)
		}
		workSession.BranchId = input.BranchId
	} else {
		workSession.BranchId = nil
	}
	if input.StartTime.After(input.EndTime) {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_START_TIME_MUST_SMALLER_THAN_END_TIME)
	}
	if repo.IsWorkSessionHasDataDependencies(workSession.ID, *user.CenterId) {
		if input.IsActive != nil && workSession.IsActive != nil && *workSession.IsActive != *input.IsActive {
			goto JustUpdateActive
		}
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.ERROR_WORK_SESSION_HAVE_DATA_DEPENDS)
	}
	workSession.Title = input.Title
	workSession.StartTime = &input.StartTime
	workSession.EndTime = &input.EndTime
	workSession.Branch = nil
JustUpdateActive:
	workSession.IsActive = input.IsActive
	row, err := repo.UpdateWorkSessionById(&workSession)
	if row < 1 && err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFail, "Update record failed")
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, workSession)
}

func DeleteWorkSession(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusUnauthorized, consts.InvalidInput, "Permission denied!")
	}
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusUnauthorized, consts.InvalidInput, "Permission denied!")
	}

	var (
		err   error
		id, _ = uuid.Parse(c.Params("id"))
	)

	_, err = repo.GetWorkSessionByIdAndCenterId(id, *user.CenterId)
	switch {
	case err == nil:
		var slots []models.TimeSlot
		if slots, err = repo.FindTimeSlots(app.Database.DB.Where("work_session_id = ?", id)); err != nil {
			return ResponseError(c, fiber.StatusInternalServerError, "Error get time slots!",
				err.Error())
		} else if len(slots) > 0 {
			return ResponseError(c, fiber.StatusConflict,
				"Ca làm đã được gán dữ liệu. Không thể xóa.", "Error validation data!")
		}

		if err = repo.DeleteWorkSessions(app.Database.DB.Where("id = ?", id)); err != nil {
			return ResponseError(c, fiber.StatusInternalServerError, consts.DeleteFail, err.Error())
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, nil)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ResponseError(c, fiber.StatusNotFound, consts.NotFound, err.Error())
	default:
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
}
