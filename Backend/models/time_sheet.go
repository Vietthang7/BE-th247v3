package models

import "github.com/google/uuid"

type SalaryHistory struct {
	Model      `gorm:"embedded"`
	Salary     int64     `json:"salary"`                // Số tiền lương
	SalaryType int64     `json:"salary_type,omitempty"` // Cách tính lương
	UserID     uuid.UUID `json:"-"`
	//User       *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OrganID  *uuid.UUID   `json:"-"`
	Organ    *OrganStruct `json:"organ,omitempty" gorm:"foreignKey:OrganID"`
	BranchID *uuid.UUID   `json:"-"`
	Branch   *Branch      `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
	CenterID *uuid.UUID   `json:"-"`
}
