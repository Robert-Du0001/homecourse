package models

import "github.com/goravel/framework/database/orm"

// 课程集模型
type Episode struct {
	orm.Model

	GroupID         uint   `json:"group_id" form:"group_id"`
	Title           string `json:"title"`
	FilePath        string `json:"file_path" form:"file_path"`
	Sort            uint   `json:"sort"`
	Duration        uint   `json:"duration"`
	WatchedDuration bool   `json:"watched_duration"`
}
