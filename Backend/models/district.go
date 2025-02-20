package models

type District struct {
	ID         int64  `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	Prefix     string `json:"prefix"`
	ProvinceId int64  `json:"province_id,omitempty"`
	Alias      string `json:"-"`
}
