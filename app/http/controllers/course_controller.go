package controllers

import (
	"fmt"
	"homecourse/app/http/response"
	"homecourse/app/models"
	"strings"

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

	if err := facades.Orm().Query().OrderBy("sort").Where(func(query orm.Query) orm.Query {
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

	var courses []models.Course

	if err := facades.Orm().Query().OrderBy("sort").Get(&courses); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "课程获取成功", courses)
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
		"category_id": "uint|exists:categories",
		"description": "string|max_len:200",
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
	}
	var requestData req

	if err := validator.Bind(&requestData); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	// 获取封面文件，不通过Bind中的获取
	// 不光是获取不到，验证也有问题
	file, err := ctx.Request().File("cover_file")
	if err != nil {
		return response.InternalServerError(ctx, "E3", nil)
	}

	// 验证 MIME 类型
	mime, err := file.MimeType()
	if err != nil {
		return response.InternalServerError(ctx, "E4", err)
	}

	// 定义允许的类型
	allowMimes := facades.Config().Get("app.allow_img_mimes")

	if !allowMimes.(map[string]bool)[mime] {
		return response.BadRequest(ctx, "不支持的文件格式: "+mime, nil)
	}

	// 保存封面文件
	coverPath, err := facades.Storage().PutFile("/covers", file)
	if err != nil {
		return response.InternalServerError(ctx, "E5", nil)
	}

	course := &models.Course{
		Title:       requestData.Title,
		CategoryID:  requestData.CategoryID,
		Description: requestData.Description,
		CoverPath:   coverPath,
	}

	if err := facades.Orm().Query().Create(course); err != nil {
		return response.InternalServerError(ctx, "E6", err)
	}

	return response.Ok(ctx, "课程创建成功", nil)
}

// 修改课程 - 管理员
func (r *CourseController) Update(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"id":          "required|uint",
		"title":       "required|string|max_len:10",
		"category_id": "uint|exists:categories",
		"description": "string|max_len:200",
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

	file, err := ctx.Request().File("cover_file")
	if err != nil {
		return response.InternalServerError(ctx, "E3", nil)
	}

	// 验证 MIME 类型
	mime, err := file.MimeType()
	if err != nil {
		return response.InternalServerError(ctx, "E4", err)
	}

	// 定义允许的类型
	allowMimes := facades.Config().Get("app.allow_img_mimes")

	if !allowMimes.(map[string]bool)[mime] {
		return response.BadRequest(ctx, "不支持的文件格式: "+mime, nil)
	}

	// 保存封面文件
	coverPath, err := facades.Storage().PutFile("/covers", file)
	if err != nil {
		return response.InternalServerError(ctx, "E5", nil)
	}

	course.CoverPath = coverPath

	// 获取旧的封面路径
	type oldCover struct {
		CoverPath string
	}
	var oldCoverPath oldCover
	if err := facades.Orm().Query().Model(&models.Course{}).
		Select("cover_path").
		Where("id", course.ID).
		Get(&oldCoverPath); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	// 更新课程内容
	if _, err := facades.Orm().Query().Model(&models.Course{}).
		Where("id", course.ID).
		Update(course); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	// 删除旧的封面
	if err := facades.Storage().Delete(oldCoverPath.CoverPath); err != nil {
		return response.InternalServerError(ctx, "E4", err)
	}

	return response.Ok(ctx, "课程修改成功", nil)
}

// 获取课程封面
func (r *CourseController) ShowCover(ctx http.Context) http.Response {
	coverPath := ctx.Request().Route("path")

	return ctx.Response().File(facades.Storage().Path("covers/" + coverPath))
}

// 修改分类排序
func (r *CourseController) UpdateSort(ctx http.Context) http.Response {
	fmt.Println(ctx.Request().All())

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

	caseSql.WriteString("UPDATE courses SET sort = CASE id ")
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
