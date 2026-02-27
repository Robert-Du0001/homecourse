package controllers

import (
	"homecourse/app/http/response"
	"homecourse/app/models"

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
	if err := facades.Orm().Query().Get(&categories); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	return response.Ok(ctx, "分类获取成功", categories)
}

// 获取分类信息 - 管理员
func (r *CategoryController) AdminIndex(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"page":  "required|uint",
		"limit": "required|uint",
	}, validation.Filters(map[string]string{
		"page":  "int",
		"limit": "int",
	}))

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

	var categories []models.Category
	var total int64

	if err := facades.Orm().Query().Paginate(
		request.Page,
		request.Limit,
		&categories,
		&total,
	); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "分类获取成功", map[string]any{
		"categories": categories,
		"total":      total,
	})
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
