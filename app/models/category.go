package models

import (
	"github.com/goravel/framework/database/orm"
)

type Category struct {
	orm.Model

	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	Sort      uint   `json:"sort"`
}
