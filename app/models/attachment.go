package models

import (
	"github.com/goravel/framework/database/orm"
)

type Attachment struct {
	orm.Model

	EpisodeID uint   `json:"episode_id" form:"episode_id"`
	Name      string `json:"name"`
	FilePath  string `json:"file_path" form:"file_path"`
}
