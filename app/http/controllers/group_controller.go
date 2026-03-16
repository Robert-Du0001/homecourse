package controllers

import (
	"fmt"
	"homecourse/app/facades"
	"homecourse/app/http/response"
	"homecourse/app/models"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/validation"
)

type GroupController struct {
	// Dependent services
}

func NewGroupController() *GroupController {
	return &GroupController{
		// Inject services
	}
}

// 获取分组列表
func (r *GroupController) Index(ctx http.Context) http.Response {
	courseId := ctx.Request().Route("id")

	var groups []models.Group

	if err := facades.Orm().Query().Model(&models.Group{}).
		Where("course_id", courseId).
		Order("sort").
		Get(&groups); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	return response.Ok(ctx, "获取成功", groups)
}

// 获取添加详情
func (r *GroupController) Store(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"course_id": "required|uint|exists:courses",
		"name":      "required|string|max_len:10",
	}, validation.Filters(map[string]string{
		"course_id": "uint",
	}))

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	var group models.Group

	if err := validator.Bind(&group); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	if err := facades.Orm().Query().Create(&group); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "分组创建成功", nil)
}

// 修改分组
func (r *GroupController) Update(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"id":        "required|int",
		"course_id": "required|uint|exists:courses",
		"name":      "required|string|max_len:10",
	}, validation.Filters(map[string]string{
		"id":        "int",
		"course_id": "uint",
	}))

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	var group models.Group

	if err := validator.Bind(&group); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	if _, err := facades.Orm().Query().Model(&models.Group{}).
		Where("id", group.ID).
		Update(group); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "分组修改成功", nil)
}

// 删除分组
func (r *GroupController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().RouteInt("id")

	if exists, err := facades.Orm().Query().Model(&models.Episode{}).
		Where("group_id", id).
		Exists(); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	} else if exists {
		return response.BadRequest(ctx, "此分组下有课程存在，请先删除课程", nil)
	}

	if _, err := facades.Orm().Query().Model(&models.Group{}).
		Where("id", id).
		Delete(); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	return response.Ok(ctx, "分组删除成功", nil)
}

// 修改分组排序
func (r *GroupController) UpdateSort(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"ids": "required",
	})

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	type req struct {
		Ids []int `json:"ids"`
	}
	var request req

	if err := validator.Bind(&request); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	// 1. 构建 CASE WHEN 字符串
	var caseSql strings.Builder
	var idsStr []string

	caseSql.WriteString("UPDATE groups SET sort = CASE id ")
	for index, id := range request.Ids {
		// 安全起见，手动拼入数字（uint 类型不涉及 SQL 注入）
		caseSql.WriteString(fmt.Sprintf("WHEN %d THEN %d ", id, index))
		idsStr = append(idsStr, fmt.Sprintf("%d", id))
	}
	caseSql.WriteString("END ")

	caseSql.WriteString(fmt.Sprintf("WHERE id IN (%s)", strings.Join(idsStr, ",")))

	// 3. 执行原生 SQL
	if _, err := facades.Orm().Query().Exec(caseSql.String()); err != nil {
		facades.Log().Error("批量排序失败", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"msg": "数据库操作失败"})
	}

	return response.Ok(ctx, "排序同步成功", nil)
}
