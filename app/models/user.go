package models

import (
	"github.com/goravel/framework/database/orm"
)

// 上下文user
const Cuser = "cuser"

// 用户角色
type UserRole uint8

const (
	// 普通用户
	RoleGuest UserRole = iota
	// 管理员
	RoleAdmin
)

type User struct {
	orm.Model
	Name     string   `json:"name" form:"name"`
	Password string   `json:"-" form:"password"`
	Role     UserRole `json:"role"`
}
