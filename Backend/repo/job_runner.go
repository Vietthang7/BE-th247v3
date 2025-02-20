package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"intern_247/app"
	"intern_247/models"
	"time"
)

type UpdateSalary struct {
	UserId     uuid.UUID
	Salary     int64
	SalaryType int64
	CenterId   *uuid.UUID
	BranchId   *uuid.UUID
	OrganID    *uuid.UUID
}

func (u UpdateSalary) Run() {
	var (
		err         error
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		tx          = app.Database.DB.WithContext(ctx).Begin()
	)
	defer cancel()
	if err = tx.Model(&models.User{}).Where("id = ?", u.UserId).Update("salary", u.Salary).Update("salary_type", u.SalaryType).Error; err != nil {
		logrus.Error(err)
		tx.Rollback()
	} else {
		history := SalaryHistory{
			UserID:     u.UserId,
			SalaryType: u.SalaryType,
			Salary:     u.Salary,
			CenterID:   u.CenterId,
			BranchID:   u.BranchId,
			OrganID:    u.OrganID,
		}
		now := time.Now()
		history.CreatedAt = &now
		if err = tx.Create(&history).Error; err != nil {
			logrus.Error("Error while create salary history: ", err)
			tx.Rollback()
		}
	}
	tx.Commit()
}
