package models

import "github.com/goravel/framework/database/orm"

// 课程集模型
type Episode struct {
	orm.Model

	UserID      uint
	CourseID    uint   `json:"course_id"`
	Title       string `json:"title"`
	FilePath    string `json:"file_path"`
	Sort        uint   `json:"sort"`
	Duration    uint   `json:"duration"`
	IsCompleted bool   `json:"is_completed"`
}
