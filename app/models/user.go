package models

import (
	"github.com/goravel/framework/database/orm"
)

// userid的上下文名
const UserID = "uid"

type User struct {
	orm.Model
	Name     string `json:"name" form:"name"`
	Password string `json:"-" form:"password"`
	Role     uint8  `json:"role"`
}
