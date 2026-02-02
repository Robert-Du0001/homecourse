package controllers

import (
	"homecourse/app/http/response"
	"homecourse/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
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
