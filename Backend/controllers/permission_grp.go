package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"strconv"
)

func CreatePermissionGrp(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error permission denied", "")
	}
	var (
		entry models.PermissionGroup
		err   error
	)
	if err = c.BodyParser(&entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.InvalidInput, err.Error())
	}
	if entry.Name != "" {
		if _, err = repo.FirstPermissionGrp(app.Database.DB.Where("name= ? AND center_id = ?", entry.Name, *user.CenterId)); err == nil {
			return ResponseError(c, fiber.StatusConflict, "Tên nhóm quyền này đã tồn tại", "")
		}
	}
	entry.CenterId = user.CenterId
	if err = repo.CreatePermissionGrp(app.Database.DB, &entry); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.CreateFailed, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.CreateSuccess, entry)
}
func ListPermissionGrp(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error permission denied", "")
	}
	var (
		err        error
		entries    []models.PermissionGroup
		pagination = consts.BindRequestTable(c, "created_at")
		DB         = pagination.CustomOptions(app.Database.DB).Where(consts.NilDeletedAt)
	)
	if pagination.Search != "" {
		DB = DB.Where("name LIKE ?", "%"+pagination.Search+"%")
	}
	if c.Query("active") != "" {
		isActive, _ := strconv.ParseBool(c.Query("active"))
		DB = DB.Where("is_active = ?", isActive)
	}
	//users just can find data belong center them
	DB = DB.Where("center_id = ?", user.CenterId)
	if entries, err = repo.FindPermissionGrp(DB); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFailed, err.Error())
	}
	pagination.Total = repo.CountPermissionGrp(DB.Offset(-1))
	for i := range entries {
		repo.PreloadPermissionGrp(&entries[i], "tag")
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}
func ReadPermissionGrp(c *fiber.Ctx) error {
	entry, err := repo.FirstPermissionGrp(app.Database.DB.Where("id = ?", c.Params("id")))
	switch {
	case err == nil:
		//repo.PreloadPermissionGrp(&entry, "tag")
		return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, consts.NotFound, err.Error())
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
}
func UpdatePermissionGrp(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error permission denied", "")
	}
	DB := app.Database.DB.Where("id = ?", c.Params("id"))
	entry, err := repo.FirstPermissionGrp(DB)
	switch {
	case err == nil:
		if err = c.BodyParser(&entry); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, err.Error())
		}
		if entry.Name != "" {
			if _, err = repo.FirstPermissionGrp(app.Database.DB.Where("id <> ? AND name = ? AND center_id = ?",
				c.Params("id"), entry.Name, *user.CenterId)); err == nil {
				return ResponseError(c, 0, "Tên nhóm quyền này đã tồn tại!", "")
			}
		}
		if err = repo.UpdatePermissionGroup(DB, &entry); err != nil {
			logrus.Error(err)
			return ResponseError(c, fiber.StatusInternalServerError, consts.UpdateFail, err.Error())
		}
		return ResponseSuccess(c, fiber.StatusOK, consts.UpdateSuccess, entry)
	case errors.Is(err, gorm.ErrRecordNotFound):
		logrus.Error(err)
		return ResponseError(c, fiber.StatusNotFound, consts.GetFail, err.Error())
	default:
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
}
func DeletePermissionGroup(c *fiber.Ctx) error {
	var (
		err    error
		reqIds models.ReqIds
	)

	if err = c.BodyParser(&reqIds); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, err.Error())
	}

	var groups []models.PermissionGroup
	if groups, err = repo.FindPermissionGrp(app.Database.DB.Where(consts.NilDeletedAt).
		Where("id IN (?)", reqIds.Ids)); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.DeleteFail, err.Error())
	} else {
		for _, v := range groups {
			if _, err = repo.FirstOrganStruct(app.Database.DB.Where(consts.NilDeletedAt).
				Where("permission_grp_id = ?", v.ID)); err == nil {
				return ResponseError(c, 0, v.Name+" đã được gán cho cơ cấu tổ chức", "")
			}

			if _, err = repo.FirstUser("permission_grp_id = ?", []interface{}{v.ID}); err == nil {
				return ResponseError(c, 0, "Nhóm quyền đã gắn cho nhân sự. Không thể xóa", "")
			}
		}
	}

	if err = repo.DeletePermissionGroup(app.Database.DB.Where(consts.NilDeletedAt).
		Where("id IN (?)", reqIds.Ids)); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.DeleteFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.DeleteSuccess, nil)
}

//func ExportListGrps(c *fiber.Ctx) error {
//	var (
//		err     error
//		entries []models.PermissionGroup
//	)
//
//	DB := app.Database.DB.Order(consts.DescCreatedAt).Where(consts.NilDeletedAt).
//		Select("name", "is_active", "created_at", "tags", "select_all")
//
//	if c.Query("search") != "" {
//		DB = DB.Where("name LIKE ?", "%"+c.Query("search")+"%")
//	}
//	if c.Query("active") != "" {
//		isActive, _ := strconv.ParseBool(c.Query("active"))
//		DB = DB.Where("is_active = ?", isActive)
//	}
//
//	if entries, err = repo.FindPermissionGrp(DB); err != nil {
//		logrus.Error(err)
//		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
//	}
//	for i := range entries {
//		repo.PreloadPermissionGrp(&entries[i], "tag")
//	}
//
//	var fileName string
//	file := excelize.NewFile()
//	if file, fileName, err = ExportExcelListGrps(entries); err != nil {
//		return err
//	}
//
//	var b bytes.Buffer
//	if err = file.Write(&b); err != nil {
//		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
//	}
//	c.Set("Content-Description", "File Transfer")
//	c.Set("Content-Disposition", "attachment; filename="+fileName)
//	return c.Status(fiber.StatusOK).Type("application/octet-stream").Send(b.Bytes())
//}
//func ExportExcelListGrps(entries []models.PermissionGroup) (file *excelize.File, fileName string, err error) {
//	var entry models.PermissionGroup
//	sheetName := "Dữ liệu"
//	file = excelize.NewFile()
//	index := file.NewSheet(sheetName)
//	file.DeleteSheet("Sheet1")
//	// set header
//	file.SetCellValue(sheetName, "A1", "STT")
//	file.SetCellValue(sheetName, "B1", "Tên nhóm quyền")
//	file.SetCellValue(sheetName, "C1", "Phân quyền")
//	file.SetCellValue(sheetName, "D1", "Ngày tạo")
//	file.SetCellValue(sheetName, "E1", "Trạng thái")
//
//	// set header style : bold font and all borders
//	bolAndAlBor, _ := file.NewStyle(`{"alignment":{"horizontal":"center"}, "font":{"bold":true},
//	"border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},
//	{"type":"bottom","color":"000000","style":1},{"type":"right","color":"000000","style":1}]}`)
//	file.SetCellStyle(sheetName, "A1", "E1", bolAndAlBor)
//
//	// style: center and all borders
//	horCen, _ := file.NewStyle(`{
//        "alignment": {
//            "horizontal": "center",
//            "vertical": "center"
//        },
//        "border": [
//            {
//                "type": "left",
//                "color": "000000",
//                "style": 1
//            },
//            {
//                "type": "top",
//                "color": "000000",
//                "style": 1
//            },
//            {
//                "type": "bottom",
//                "color": "000000",
//                "style": 1
//            },
//            {
//                "type": "right",
//                "color": "000000",
//                "style": 1
//            }
//        ]
//    }`)
//
//	// style: all borders
//	alBor, _ := file.NewStyle(`{
//        "border": [
//            {
//                "type": "left",
//                "color": "000000",
//                "style": 1
//            },
//            {
//                "type": "top",
//                "color": "000000",
//                "style": 1
//            },
//            {
//                "type": "bottom",
//                "color": "000000",
//                "style": 1
//            },
//            {
//                "type": "right",
//                "color": "000000",
//                "style": 1
//            }
//        ],
//        "alignment": {
//            "wrap_text": true,
//            "vertical": "center"
//        }
//    }`)
//
//	var names []string
//
//	// Fill Data
//	countLine := 2
//	for _, entry = range entries {
//		file.SetCellValue(sheetName, "A"+strconv.Itoa(countLine), countLine-1)
//		file.SetCellValue(sheetName, "B"+strconv.Itoa(countLine), entry.Name)
//		//file.SetCellValue(sheetName, "C"+strconv.Itoa(countLine), FormatExcelGroupTags(entry.Tags))
//		file.SetCellValue(sheetName, "D"+strconv.Itoa(countLine), entry.CreatedAt.Add(7*time.Hour).
//			Format("02/01/2006"))
//		if *entry.IsActive {
//			file.SetCellValue(sheetName, "E"+strconv.Itoa(countLine), "Hoạt động")
//		} else {
//			file.SetCellValue(sheetName, "E"+strconv.Itoa(countLine), "Dừng hoạt động")
//		}
//
//		// Set all Borders
//		file.SetCellStyle(sheetName, "A"+strconv.Itoa(countLine), "E"+strconv.Itoa(countLine), alBor)
//		file.SetCellStyle(sheetName, "A"+strconv.Itoa(countLine), "A"+strconv.Itoa(countLine), horCen)
//
//		names = append(names, entry.Name)
//
//		countLine++
//	}
//
//	// set column width
//	file.SetColWidth(sheetName, "B", "B", float64(LongestElem(names))+4)
//	file.SetColWidth(sheetName, "C", "C", 50)
//	file.SetColWidth(sheetName, "D", "E", 20)
//
//	file.SetActiveSheet(index)
//	currentDate := time.Now().UTC().Format("01-02-2006")
//	currentTime := time.Now().UTC().Format("2006-01-02 15_04_05")
//	fileName = "Export_DSNhomQuyen_" + currentTime + ".xlsx"
//	directory := "./public"
//	if _, err = os.Stat(directory); os.IsNotExist(err) {
//		if err = os.Mkdir(directory, os.ModePerm); err != nil {
//			logrus.Error(err)
//		}
//	}
//	path, _ := os.Getwd()
//	_ = os.Mkdir(filepath.Join(path, "public", currentDate), 0755)
//	pathFile := filepath.Join(path, "public", currentDate, fileName)
//
//	if err = file.SaveAs(pathFile); err != nil {
//		return
//	}
//	return
//}
