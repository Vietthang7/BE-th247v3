package models

import "github.com/google/uuid"

type TuitionFeePkgCf struct {
	Model       `gorm:"embedded"`
	OneTime     *bool      `json:"one_time" gorm:"default:true"`  // Thanh toán 1 lần
	Installment *bool      `json:"installment"`                   // Thanh toán theo từng đợt
	PerClass    *bool      `json:"per_class" gorm:"default:true"` // Thanh toán theo từng lớp
	IsTemplate  bool       `json:"-" gorm:"default:false"`
	CenterID    *uuid.UUID `json:"-"`
}
