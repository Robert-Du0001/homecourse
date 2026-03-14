package models

type Category struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	Sort      uint   `json:"sort"`
}
