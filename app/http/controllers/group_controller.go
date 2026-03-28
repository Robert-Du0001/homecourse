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
	courseID := ctx.Request().RouteInt("id")

	var groups []models.Group

	if err := facades.Orm().Query().Where("course_id", courseID).
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
	groupId := ctx.Request().RouteInt("id")

	var group models.Group
	if err := facades.Orm().Query().Select("id", "is_default").
		Find(&group, groupId); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	// 默认分类不能被删除
	if group.IsDefault {
		return response.BadRequest(ctx, "默认分组不能被删除", nil)
	}

	if exists, err := facades.Orm().Query().Model(&models.Episode{}).
		Where("group_id", groupId).
		Exists(); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	} else if exists {
		return response.BadRequest(ctx, "此分组下有剧集存在，请先删除剧集", nil)
	}

	if _, err := facades.Orm().Query().Model(&models.Group{}).
		Where("id", groupId).
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

// 设置默认分组
func (r *GroupController) UpdateDefault(ctx http.Context) http.Response {
	groupId := ctx.Request().Route("id")

	tx, err := facades.Orm().Query().BeginTransaction()

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	// 判断剧集分组是否存在
	if exists, err := facades.Orm().Query().
		Model(&models.Group{}).
		Where("id", groupId).
		Exists(); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	} else if !exists {
		return response.BadRequest(ctx, "课程分组不存在", nil)
	}

	// 把当前分组ID设置成默认
	if _, err := tx.Model(&models.Group{}).
		Where("id", groupId).
		Update("is_default", true); err != nil {
		if err := tx.Rollback(); err != nil {
			return response.InternalServerError(ctx, "E2", err)
		}
		return response.InternalServerError(ctx, "E3", err)
	}

	// 把其他分组ID取消默认
	if _, err := tx.Model(&models.Group{}).
		Where("id <> ?", groupId).
		Update("is_default", false); err != nil {
		if err := tx.Rollback(); err != nil {
			return response.InternalServerError(ctx, "E4", err)
		}
		return response.InternalServerError(ctx, "E5", err)
	}

	if err := tx.Commit(); err != nil {
		return response.InternalServerError(ctx, "E6", err)
	}

	return response.Ok(ctx, "设置默认课程分类成功", nil)
}
