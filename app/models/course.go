package models

import (
	"github.com/goravel/framework/database/orm"
)

type Course struct {
	orm.Model

	CategoryID  uint   `json:"category_id" form:"category_id"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	CoverPath   string `json:"cover_path" form:"cover_path"`
	Status      uint8  `json:"status" form:"status"`
}
