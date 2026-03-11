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
		"category_id": "uint",
	}, validation.Filters(map[string]string{
		"category_id": "uint",
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
		if requestData.CategoryID != 0 {
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

// 获取课程列表 - 管理员
func (r *CourseController) AdminIndex(ctx http.Context) http.Response {
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

	var courses []models.Course
	var total int64

	if err := facades.Orm().Query().Paginate(
		request.Page,
		request.Limit,
		&courses,
		&total,
	); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "课程获取成功", map[string]any{
		"courses": courses,
		"total":   total,
	})
}

// 删除课程 - 管理员
func (r *CourseController) Destroy(ctx http.Context) http.Response {
	courseId := ctx.Request().Route("id")

	// 如果课程下有剧集，则不能删除
	if exists, err := facades.Orm().Query().Model(&models.Episode{}).
		Where("course_id", courseId).Exists(); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	} else if exists {
		return response.BadRequest(ctx, "课程下有剧集，请先删除剧集", nil)
	}

	if _, err := facades.Orm().Query().Model(&models.Course{}).
		Where("id", courseId).Delete(); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	return response.Ok(ctx, "删除成功", nil)
}

// 添加课程 - 管理员
func (r *CourseController) Store(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"title":       "required|string|max_len:10",
		"category_id": "required|uint|exists:categories",
		"description": "string|max_len:200",
		"cover_path":  "string|max_len:255",
	}, validation.Filters(map[string]string{
		"category_id": "uint",
	}))

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	type req struct {
		Title       string `form:"title"`
		CategoryID  uint   `form:"category_id"`
		Description string `form:"description"`
		CoverPath   string `form:"cover_path"`
	}
	var requestData req

	if err := validator.Bind(&requestData); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	course := &models.Course{
		Title:       requestData.Title,
		CategoryID:  requestData.CategoryID,
		Description: requestData.Description,
		CoverPath:   requestData.CoverPath,
	}

	if err := facades.Orm().Query().Create(course); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "课程创建成功", nil)
}

// 修改课程 - 管理员
func (r *CourseController) Update(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"id":          "required|uint",
		"title":       "required|string|max_len:10",
		"category_id": "required|uint|exists:categories",
		"description": "string|max_len:200",
		"cover_path":  "string|max_len:255",
	}, validation.Filters(map[string]string{
		"id":          "uint",
		"category_id": "uint",
	}))

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	var course models.Course

	if err := validator.Bind(&course); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	if _, err := facades.Orm().Query().Model(&models.Course{}).
		Where("id", course.ID).
		Update(course); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "课程修改成功", nil)
}
