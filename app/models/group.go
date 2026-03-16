package models

import (
	"github.com/goravel/framework/database/orm"
)

type Group struct {
	orm.Model

	ID       uint   `gorm:"primaryKey" json:"id"`
	CourseID uint   `json:"course_id" form:"course_id"`
	Name     string `json:"name"`
	Sort     uint   `json:"sort"`
}
