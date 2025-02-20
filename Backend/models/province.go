package models

type Province struct {
	ID    int64  `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Code  string `json:"code"`
	Alias string `json:"-"`
}
type Ward struct {
	ID         int64  `json:"id" gorm:"primary_key"`
	Name       string `json:"name"`
	DistrictId string `json:"district_id"`
}
