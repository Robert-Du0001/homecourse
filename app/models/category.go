package models

import "github.com/goravel/framework/database/orm"

type Category struct {
	orm.Model

	UserID   uint
	Name     string `json:"name"`
	ParentID uint   `json:"parent_id"`
	Sort     uint   `json:"sort"`
}
