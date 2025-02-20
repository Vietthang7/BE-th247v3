package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Certificate struct {
	Model     `gorm:"embedded"`
	Name      string         `gorm:"size:250;index" json:"name"`
	IsActive  *bool          `gorm:"default:true" json:"is_active,omitempty"`
	ImageUrl  string         `json:"image_url"`
	CreatedBy uuid.UUID      `json:"-"`
	Metadata  datatypes.JSON `json:"metadata"`
	CenterId  uuid.UUID      `json:"-"`
	Creator   *User          `gorm:"foreignKey:CreatedBy" json:"created_by,omitempty"`
	Subjects  *[]Subject     `gorm:"foreignKey:CertificateId" json:"subjects,omitempty"`
	Center    Center         `gorm:"foreignKey:CenterId" json:"-"`
}
