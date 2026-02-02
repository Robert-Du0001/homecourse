package models

import "github.com/goravel/framework/database/orm"

type Course struct {
	orm.Model

	UserID      uint
	CategoryID  uint   `json:"category_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverPath   string `json:"cover_path"`
	Status      uint8  `json:"status"`
}
