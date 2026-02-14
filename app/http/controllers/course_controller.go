package controllers

import (
	"homecourse/app/http/response"
	"homecourse/app/models"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/validation"
)

type CourseController struct {
	// Dependent services
}

func NewCourseController() *CourseController {
	return &CourseController{
		// Inject services
	}
}

// 获取课程列表
func (r *CourseController) Index(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"category_id": "int:-1",
	}, validation.Filters(map[string]string{
		"category_id": "int",
	}))

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	type req struct {
		CategoryID int `form:"category_id"`
	}
	var requestData req

	if err := validator.Bind(&requestData); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	var courses []models.Course

	if err := facades.Orm().Query().Where(func(query orm.Query) orm.Query {
		if requestData.CategoryID != -1 {
			return query.Where("category_id", requestData.CategoryID)
		}
		return query
	}).
		Find(&courses); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "课程获取成功", courses)
}

// 获取课程详情
func (r *CourseController) Show(ctx http.Context) http.Response {
	courseId := ctx.Request().Route("id")

	var course models.Course

	if err := facades.Orm().Query().FindOrFail(&course, courseId); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return response.BadRequest(ctx, "课程不存在", nil)
		}

		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "课程详情获取成功", course)
}
