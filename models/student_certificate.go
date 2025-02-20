package models

import (
	"github.com/google/uuid"
	"time"
)

type StudentCertificates struct {
	Model         `gorm:"embedded"`
	StudentId     uuid.UUID    `json:"-"`
	CertificateId uuid.UUID    `json:"-"`
	ClassId       uuid.UUID    `json:"-"`
	SubjectId     uuid.UUID    `json:"-"`
	Status        uint         `gorm:"index" json:"status"`
	CreatedBy     uuid.UUID    `json:"-"`
	Code          string       `gorm:"default:NULL" json:"-"`
	ApprovedBy    *uuid.UUID   `gorm:"default:null" json:"-"`
	ApprovedAt    *time.Time   `gorm:"default:null" json:"approved_at,omitempty"`
	Description   string       `gorm:"default:null" json:"description,omitempty"`
	Student       *Student     `gorm:"foreignKey:StudentId" json:"student,omitempty"`
	Certificate   *Certificate `gorm:"foreignKey:CertificateId" json:"certificate,omitempty"`
	Class         *Class       `gorm:"foreignKey:ClassId" json:"class,omitempty"`
	Subject       *Subject     `gorm:"foreignKey:SubjectId" json:"subject,omitempty"`
	Creator       *User        `gorm:"foreignKey:CreatedBy" json:"created_by,omitempty"`
	Approver      *User        `gorm:"foreignKey:ApprovedBy" json:"approved_by,omitempty"`
	CenterId      uuid.UUID    `json:"-"`
	Center        *Center      `gorm:"foreignKey:CenterId" json:"center,omitempty"`
}
