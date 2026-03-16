package models

import "github.com/goravel/framework/database/orm"

// 课程集模型
type Episode struct {
	orm.Model

	CourseID    uint   `json:"course_id" form:"course_id"`
	GroupID     uint   `json:"group_id" form:"group_id"`
	Title       string `json:"title"`
	FilePath    string `json:"file_path" form:"file_path"`
	Sort        uint   `json:"sort"`
	Duration    uint   `json:"duration"`
	IsCompleted bool   `json:"is_completed"`
}
