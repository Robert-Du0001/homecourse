package controllers

import (
	"fmt"
	"homecourse/app/http/response"
	"homecourse/app/models"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/validation"
)

type CategoryController struct {
	// Dependent services
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		// Inject services
	}
}

// 获取分类信息
func (r *CategoryController) Index(ctx http.Context) http.Response {
	var categories []models.Category

	if err := facades.Orm().Query().OrderBy("sort").Get(&categories); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "分类获取成功", categories)
}

// 修改默认课程分类
func (r *CategoryController) UpdateDefault(ctx http.Context) http.Response {
	categoryId := ctx.Request().Route("id")

	tx, err := facades.Orm().Query().BeginTransaction()

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	// 判断课程分类是否存在
	if exists, err := facades.Orm().Query().
		Model(&models.Category{}).
		Where("id", categoryId).
		Exists(); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	} else if !exists {
		return response.BadRequest(ctx, "课程分类不存在", nil)
	}

	if _, err := tx.Model(&models.Category{}).
		Where("id", categoryId).
		Update("is_default", true); err != nil {
		if err := tx.Rollback(); err != nil {
			return response.InternalServerError(ctx, "E2", err)
		}
		return response.InternalServerError(ctx, "E3", err)
	}

	if _, err := tx.Model(&models.Category{}).
		Where("id <> ?", categoryId).
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

// 新增分类
func (r *CategoryController) Store(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"name": "required|string|max_len:10",
	})

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	var category models.Category

	if err := validator.Bind(&category); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	if err := facades.Orm().Query().Create(&category); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "分类创建成功", category)
}

// 删除分类
func (r *CategoryController) Destroy(ctx http.Context) http.Response {
	categoryId := ctx.Request().Route("id")

	// 判断此分类是否有课程存在
	if exists, err := facades.Orm().Query().Model(&models.Course{}).
		Where("category_id", categoryId).
		Exists(); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	} else if exists {
		return response.BadRequest(ctx, "此分类下有课程存在，请先删除课程", nil)
	}

	if _, err := facades.Orm().Query().Model(&models.Category{}).
		Where("id", categoryId).
		Delete(); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	return response.Ok(ctx, "分类删除成功", nil)
}

// 修改分类
func (r *CategoryController) Update(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"id":   "required|int",
		"name": "required|string|max_len:10",
	}, validation.Filters(map[string]string{
		"id": "int",
	}))

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	var category models.Category

	if err := validator.Bind(&category); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	if _, err := facades.Orm().Query().
		Where("id", category.ID).
		Update(&category); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "分类修改成功", map[string]any{
		"name": category.Name,
	})
}

// 修改分类排序
func (r *CategoryController) UpdateSort(ctx http.Context) http.Response {
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

	caseSql.WriteString("UPDATE categories SET sort = CASE id ")
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
