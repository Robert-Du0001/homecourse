package controllers

import (
	"fmt"
	"homecourse/app/http/response"
	"homecourse/app/models"
	"slices"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/str"
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
		return response.InternalServerError(ctx, "E6", err)
	}

	// 更新课程内容
	if _, err := facades.Orm().Query().Model(&models.Course{}).
		Where("id", course.ID).
		Update(course); err != nil {
		return response.InternalServerError(ctx, "E7", err)
	}

	// 删除旧的封面
	if oldCoverPath.CoverPath != "" {
		if err := facades.Storage().Delete(oldCoverPath.CoverPath); err != nil {
			return response.InternalServerError(ctx, "E8", err)
		}
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

// 扫描课程
func (r *CourseController) Scan(ctx http.Context) http.Response {
	files, err := facades.Storage().AllFiles("courses")
	if err != nil {
		return response.InternalServerError(ctx, "E0", err)
	}

	// 1. 预处理扫描到的文件数据
	scannedCourseNames := make([]string, 0)
	type tempEpisode struct {
		Title      string
		FilePath   string
		CourseName string
	}
	scannedEpisodes := make([]tempEpisode, 0, len(files))

	for _, file := range files {
		parts := str.Of(file).Replace("\\", "/").Split("/")

		courseName := parts[0]
		episodeName := parts[1]

		// 创建课程
		if !slices.Contains(scannedCourseNames, courseName) {
			scannedCourseNames = append(scannedCourseNames, courseName)
		}

		if len(parts) >= 2 {
			scannedEpisodes = append(scannedEpisodes, tempEpisode{
				Title:      episodeName,
				FilePath:   "/courses/" + courseName + "/" + episodeName,
				CourseName: courseName,
			})
		}
	}

	// 2. 处理 Group：查重并批量插入

	type course struct {
		ID      uint
		Title   string
		GroupID uint
	}

	var existingCourses []course
	courseArgs := make([]any, len(scannedCourseNames))
	for i, v := range scannedCourseNames {
		courseArgs[i] = v
	}

	if err := facades.Orm().Query().Table("courses AS c").
		Select("c.id", "c.title", "g.id AS group_id").
		Join("LEFT JOIN groups AS g on g.course_id = c.id").
		WhereIn("c.title", courseArgs).
		Where("g.is_default", true).
		Get(&existingCourses); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	// 构建已有课程 Map，方便快速查找
	// 课程与分组ID的映射
	CourseGroupIDMap := make(map[string]uint)
	for _, c := range existingCourses {
		CourseGroupIDMap[c.Title] = c.GroupID
	}

	// 获取默认课程分类
	type category struct {
		ID uint
	}
	var defaultCategory category
	if err := facades.Orm().Query().Where("is_default", true).FirstOrFail(&defaultCategory); err != nil {
		if !errors.Is(err, errors.OrmRecordNotFound) {
			return response.InternalServerError(ctx, "E2", err)
		}

		// 如果默认课程分类不存在，则无分类
		defaultCategory.ID = 0
	}

	// 找出需要新创建的课程
	var newGroups []models.Group
	var newCourses []models.Course
	for _, name := range scannedCourseNames {
		if _, exists := CourseGroupIDMap[name]; !exists {
			newCourses = append(
				newCourses,
				models.Course{
					CategoryID: defaultCategory.ID,
					Title:      name,
				})
		}
	}

	tx, _ := facades.Orm().Query().BeginTransaction()

	// 同时创建个默认分组
	if len(newCourses) > 0 {
		if err := tx.Create(&newCourses); err != nil {
			if err := tx.Rollback(); err != nil {
				return response.InternalServerError(ctx, "E6", err)
			}
			return response.InternalServerError(ctx, "E3", err)
		} else {
			// 根据新创建的课程ID，设置默认分组
			for _, nc := range newCourses {
				newGroups = append(newGroups, models.Group{
					CourseID:  nc.ID,
					IsDefault: true,
					Name:      "默认分组",
				})
			}

			// 创建默认分组
			if err := tx.Create(&newGroups); err != nil {
				if err := tx.Rollback(); err != nil {
					return response.InternalServerError(ctx, "E7", err)
				}
				return response.InternalServerError(ctx, "E8", err)
			} else {
				// 填充课程到剧集分组的映射
				for _, nc := range newCourses {
					for _, ng := range newGroups {
						if nc.ID == ng.CourseID {
							CourseGroupIDMap[nc.Title] = ng.ID
							break
						}
					}
				}
			}

		}
	}

	// 3. 处理 Episode：查重并批量插入
	// 提取所有扫描到的路径
	allPaths := make([]any, len(scannedEpisodes))
	for i, e := range scannedEpisodes {
		allPaths[i] = e.FilePath
	}

	var existingEpisodes []models.Episode
	if err := tx.WhereIn("file_path", allPaths).Get(&existingEpisodes); err != nil {
		if err := tx.Rollback(); err != nil {
			return response.InternalServerError(ctx, "E9", err)
		}
		return response.InternalServerError(ctx, "E4", err)
	}

	// 构建已有 Episode 的 Map
	existingPathMap := make(map[string]bool)
	for _, e := range existingEpisodes {
		existingPathMap[e.FilePath] = true
	}

	// 过滤出真正需要插入的单集
	var finalEpisodes []models.Episode
	for _, se := range scannedEpisodes {
		if _, exists := existingPathMap[se.FilePath]; !exists {
			finalEpisodes = append(finalEpisodes, models.Episode{
				GroupID:  CourseGroupIDMap[se.CourseName],
				Title:    se.Title,
				FilePath: se.FilePath,
			})
		}
	}

	if len(finalEpisodes) > 0 {
		if err := tx.Create(&finalEpisodes); err != nil {
			if err := tx.Rollback(); err != nil {
				return response.InternalServerError(ctx, "E10", err)
			}
			return response.InternalServerError(ctx, "E5", err)
		} else {
			if err := tx.Commit(); err != nil {
				return response.InternalServerError(ctx, "E11", err)
			}
		}
	}

	return response.Ok(ctx, "扫描完成", map[string]int{
		"new_courses":  len(newCourses),
		"new_episodes": len(finalEpisodes),
	})
}

// 统计课程
func (r *CourseController) Statistic(ctx http.Context) http.Response {
	total, err := facades.Orm().Query().Model(&models.Course{}).Count()
	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	return response.Ok(ctx, "获取成功", map[string]any{
		"total": total,
	})
}
