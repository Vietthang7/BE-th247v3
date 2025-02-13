package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"intern_247/consts"
	"intern_247/repo"
	"strconv"
)

func ListProvinces(c *fiber.Ctx) error {
	var (
		err        error
		entries    repo.Provinces
		entry      repo.Province
		pagination = consts.BindRequestTable(c, "id")
	)
	pagination.Dir = "asc"

	if entries, err = entry.Find(&pagination, "", nil); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}

func ListWards(c *fiber.Ctx) error {
	var (
		err        error
		entries    repo.Wards
		entry      repo.Ward
		pagination = consts.BindRequestTable(c, "id")
		query      = ""
		args       []interface{}
	)
	pagination.Dir = "asc"

	if c.Query("district") != "" {
		districtId, _ := strconv.Atoi(c.Query("district"))
		query += "district_id = ?"
		args = append(args, districtId)
	}

	if entries, err = entry.Find(&pagination, query, args); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}

func ListDistricts(c *fiber.Ctx) error {
	var (
		err        error
		entries    repo.Districts
		entry      repo.District
		pagination = consts.BindRequestTable(c, "id")
		query      = ""
		args       []interface{}
	)
	pagination.Dir = "asc"

	if c.Query("province") != "" {
		provinceId, _ := strconv.Atoi(c.Query("province"))
		query += "province_id = ?"
		args = append(args, provinceId)
	}

	if entries, err = entry.Find(&pagination, query, args); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError, consts.GetFail, err.Error())
	}
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}
