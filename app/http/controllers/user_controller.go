package controllers

import (
	"homecourse/app/http/response"
	"homecourse/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"
)

type UserController struct {
	// Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		// Inject services
	}
}

// 获取用户信息
func (r *UserController) Show(ctx http.Context) http.Response {
	cuser := ctx.Value(models.Cuser).(models.User)

	var user models.User

	if err := facades.Orm().Query().Find(&user, cuser.ID); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	return response.Ok(ctx, "获取成功", user)
}

// 注册用户
func (r *UserController) Store(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"name":     "required|max_len:10",
		"password": "required|min_len:8|max_len:20",
	})

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	var user models.User

	if err := validator.Bind(&user); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	if exists, err := facades.Orm().Query().Model(&models.User{}).
		Where("name", user.Name).Exists(); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	} else if exists {
		return response.BadRequest(ctx, "已存在相同账号名", nil)
	}

	var password string
	if password, err = facades.Hash().Make(user.Password); err != nil {
		return response.InternalServerError(ctx, "E4", err)
	}
	user.Password = password

	// 判断是否为第一个用户，第一个用户自动会成为管理员
	if exists, err := facades.Orm().Query().Model(&models.User{}).Exists(); err != nil {
		return response.InternalServerError(ctx, "E6", err)
	} else if !exists {
		// 不存在数据，则设置为管理员
		user.Role = models.RoleAdmin
	}

	if err := facades.Orm().Query().Create(&user); err != nil {
		return response.InternalServerError(ctx, "E5", err)
	}

	return response.Ok(ctx, "注册成功", nil)
}

// 登录
func (r *UserController) Login(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"name":     "required|max_len:10",
		"password": "required|min_len:8|max_len:20",
	})

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	type req struct {
		Name     string `form:"name"`
		Password string `form:"password"`
	}
	var visitor req

	if err := validator.Bind(&visitor); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	type userInfo struct {
		ID       uint
		Role     models.UserRole
		Password string
	}
	var user userInfo

	// 根据用户名查询用户
	if err := facades.Orm().Query().Model(&models.User{}).
		Select("id, password, role").
		Where("name", visitor.Name).
		FirstOrFail(&user); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return response.BadRequest(ctx, "用户不存在", nil)
		}

		return response.InternalServerError(ctx, "E3", err)
	}

	if !facades.Hash().Check(visitor.Password, user.Password) {
		return response.BadRequest(ctx, "密码错误", nil)
	}

	// 生成jwt
	var token string
	if token, err = facades.Auth(ctx).LoginUsingID(user.ID); err != nil {
		return response.InternalServerError(ctx, "E4", err)
	}

	return response.Ok(ctx, "登录成功", map[string]any{
		"id":    user.ID,
		"name":  visitor.Name,
		"role":  user.Role,
		"token": token,
	})
}

// 获取用户列表 - 管理员
func (r *UserController) Index(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"page":  "required|uint",
		"limit": "required|uint",
	})

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	type req struct {
		Page  int `form:"page"`
		Limit int `form:"limit"`
	}
	var request req

	if err := validator.Bind(&request); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	var users []models.User
	var total int64

	if err := facades.Orm().Query().Paginate(
		request.Page,
		request.Limit,
		&users,
		&total,
	); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "获取成功", map[string]any{
		"users": users,
		"total": total,
	})
}
