package controllers

import (
	"fmt"
	"homecourse/app/http/response"
	"homecourse/app/models"
	netHttp "net/http"
	"path/filepath"
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
	if exists, err := facades.Orm().Query().Table("episodes").
		Join("JOIN groups ON episodes.group_id = groups.id").
		Where("groups.course_id", courseId).
		Exists(); err != nil {
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
		"title":       "required|string|max_len:30",
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
	var coverPath string
	file, err := ctx.Request().File("cover_file")
	if err != nil {
		// 如果不是文件没上传的错误，则需要报错
		if !errors.Is(err, netHttp.ErrMissingFile) {
			return response.InternalServerError(ctx, "E3", err)
		}
	} else {
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
		coverPath, err = facades.Storage().PutFile("/covers", file)
		if err != nil {
			return response.InternalServerError(ctx, "E5", nil)
		}
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
		"title":       "required|string|max_len:30",
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
		// 如果不是文件没上传的错误，则需要报错
		if !errors.Is(err, netHttp.ErrMissingFile) {
			return response.InternalServerError(ctx, "E3", err)
		}
	} else { // 有文件上传
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
			return response.InternalServerError(ctx, "E5", err)
		}

		course.CoverPath = coverPath
	}

	// 获取旧的封面路径
	type oldCover struct {
		CoverPath string
	}
	var oldCoverPath oldCover
	if course.CoverPath != "" {
		if err := facades.Orm().Query().Model(&models.Course{}).
			Select("cover_path").
			Where("id", course.ID).
			Get(&oldCoverPath); err != nil {
			return response.InternalServerError(ctx, "E6", err)
		}
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
		return response.InternalServerError(ctx, "E1", err)
	}

	// 1. 预处理数据结构
	type tempEpisode struct {
		Title      string
		FilePath   string
		GroupName  string // 临时存储分组名
		CourseName string
	}

	// 记录每个课程下有哪些分组，以及是否有直接文件
	type courseMeta struct {
		Groups         []string
		HasDirectFiles bool
		DefaultGroup   string
	}

	scannedEpisodes := make([]tempEpisode, 0, len(files))
	courseMap := make(map[string]*courseMeta)

	for _, file := range files {
		// 标准化路径并分割
		pathStr := str.Of(file).Replace("\\", "/").LTrim("/").String()
		parts := strings.Split(pathStr, "/")

		if len(parts) < 2 {
			continue
		} // 忽略根目录下的文件（必须在根目录下创建子目录，否则不会被扫描）

		courseName := parts[0]
		var groupName string
		var episodeTitle string

		if _, ok := courseMap[courseName]; !ok {
			courseMap[courseName] = &courseMeta{Groups: make([]string, 0)}
		}

		if len(parts) == 2 {
			// 情况1: /courses/吉他/6.mp4 -> 直接文件
			groupName = "__COURSE_DIRECT_FILE__"
			episodeTitle = strings.TrimSuffix(parts[1], filepath.Ext(parts[1]))
			courseMap[courseName].HasDirectFiles = true
		} else {
			// 情况2: /courses/吉他/进阶/测试/4.mp4 -> 子目录
			groupName = parts[1]
			// 剩余部分平铺作为标题: "测试/4"
			subParts := parts[2:]
			fullTitle := strings.Join(subParts, "/")
			episodeTitle = strings.TrimSuffix(fullTitle, filepath.Ext(fullTitle))
		}

		// 记录分组（查重添加）
		if groupName != "__COURSE_DIRECT_FILE__" && !slices.Contains(courseMap[courseName].Groups, groupName) {
			courseMap[courseName].Groups = append(courseMap[courseName].Groups, groupName)
		}

		scannedEpisodes = append(scannedEpisodes, tempEpisode{
			Title:      episodeTitle,
			FilePath:   "/courses/" + pathStr,
			GroupName:  groupName,
			CourseName: courseName,
		})
	}

	// 确定每个课程的默认分组
	for _, meta := range courseMap {
		if meta.HasDirectFiles {
			meta.DefaultGroup = "默认分组"
		} else if len(meta.Groups) > 0 {
			// 按照字母顺序或扫描顺序取第一个子目录
			meta.DefaultGroup = meta.Groups[0]
		}
	}

	// 2. 数据库操作
	tx, _ := facades.Orm().Query().BeginTransaction()

	// 获取默认分类
	var defaultCategory models.Category
	if err := facades.Orm().Query().Where("is_default", true).
		First(&defaultCategory); err != nil {
		if !errors.Is(err, errors.OrmRecordNotFound) {
			return response.InternalServerError(ctx, "E2", err)
		}

		defaultCategory.ID = 0
	}

	// 处理课程
	allCourseNames := make([]any, 0, len(courseMap))
	for name := range courseMap {
		allCourseNames = append(allCourseNames, name)
	}

	var existingCourses []models.Course
	if err := tx.Table("courses").WhereIn("title", allCourseNames).
		Get(&existingCourses); err != nil {
		if err := tx.Rollback(); err != nil {
			return response.InternalServerError(ctx, "E3", err)
		}
		return response.InternalServerError(ctx, "E4", err)
	}

	existingCourseMap := make(map[string]uint)
	for _, c := range existingCourses {
		existingCourseMap[c.Title] = c.ID
	}

	// 创建新课程
	for name := range courseMap {
		if _, exists := existingCourseMap[name]; !exists {
			newC := models.Course{Title: name, CategoryID: defaultCategory.ID}
			if err := tx.Create(&newC); err != nil {
				if err := tx.Rollback(); err != nil {
					return response.InternalServerError(ctx, "E5", err)
				}
				return response.InternalServerError(ctx, "E6", err)
			}
			existingCourseMap[name] = newC.ID
		}
	}

	// 处理分组 (Group)
	// key: CourseID_GroupName -> value: GroupID
	groupMapping := make(map[string]uint)
	// 记录 GroupID -> 现有剧集数量
	groupCurrentCount := make(map[uint]uint)

	for courseName, meta := range courseMap {
		courseID := existingCourseMap[courseName]

		// 查出该课程下现有的所有分组
		var existingGroups []models.Group
		if err := tx.Where("course_id", courseID).Get(&existingGroups); err != nil {
			if err := tx.Rollback(); err != nil {
				return response.InternalServerError(ctx, "E7", err)
			}
			return response.InternalServerError(ctx, "E8", err)
		}

		var currentDefaultGroupID uint
		hasDefault := false
		for _, g := range existingGroups {
			if g.IsDefault {
				currentDefaultGroupID = g.ID
				hasDefault = true
				break
			}
		}

		// 如果数据库里已经有分组了但没设默认，或者这是个完全的新课程
		// 需要从 meta.Groups 里选一个作为未来的默认组（如果数据库没有的话）
		existingGroupNameMap := make(map[string]uint)
		for _, g := range existingGroups {
			existingGroupNameMap[g.Name] = g.ID
			groupMapping[fmt.Sprintf("%d_%s", courseID, g.Name)] = g.ID

			// 查询该分组下已有的剧集数量，作为 sort 的起点
			var count int64
			if count, err = tx.Table("episodes").Where("group_id", g.ID).Count(); err != nil {
				if err := tx.Rollback(); err != nil {
					return response.InternalServerError(ctx, "E5", err)
				}
				return response.InternalServerError(ctx, "E6", err)
			}
			groupCurrentCount[g.ID] = uint(count)
		}

		// 创建缺失的分组
		for _, gName := range meta.Groups {
			if _, exists := existingGroupNameMap[gName]; !exists {
				// 只有当该课程目前【完全没有】默认分组时，才把新创建的第一个分组设为默认
				isDefault := false
				if !hasDefault {
					isDefault = true
					hasDefault = true // 标记已有，后面的新组就不再设为默认了
				}

				newG := models.Group{
					CourseID:  courseID,
					Name:      gName,
					IsDefault: isDefault,
				}
				if err := tx.Create(&newG); err != nil {
					if err := tx.Rollback(); err != nil {
						return response.InternalServerError(ctx, "E9", err)
					}
					return response.InternalServerError(ctx, "E10", err)
				}
				groupMapping[fmt.Sprintf("%d_%s", courseID, gName)] = newG.ID

				// 新创建的分组，起始数量为 0
				groupCurrentCount[newG.ID] = 0
			}
		}

		// --- 关键：处理那些直接放在课程根目录下的文件 ---
		// 如果没有子目录，也没有现成的分组，此时必须强行创建一个
		if meta.HasDirectFiles && !hasDefault {
			newG := models.Group{
				CourseID:  courseID,
				Name:      "默认分组",
				IsDefault: true,
			}
			if err := tx.Create(&newG); err != nil {
				return response.InternalServerError(ctx, "E16", err)
			}
			currentDefaultGroupID = newG.ID
			groupMapping[fmt.Sprintf("%d_%s", courseID, "默认分组")] = newG.ID
		}

		// 将占位符映射到最终的默认分组 ID 上
		groupMapping[fmt.Sprintf("%d___COURSE_DIRECT_FILE__", courseID)] = currentDefaultGroupID
	}

	// 3. 处理剧集 (Episode)
	allPaths := make([]any, len(scannedEpisodes))
	for i, e := range scannedEpisodes {
		allPaths[i] = e.FilePath
	}

	var existingEpisodes []models.Episode
	if err := tx.WhereIn("file_path", allPaths).Get(&existingEpisodes); err != nil {
		if err := tx.Rollback(); err != nil {
			return response.InternalServerError(ctx, "E11", err)
		}
		return response.InternalServerError(ctx, "E12", err)
	}

	existingPathMap := make(map[string]bool)
	for _, e := range existingEpisodes {
		existingPathMap[e.FilePath] = true
	}

	var finalEpisodes []models.Episode
	for _, se := range scannedEpisodes {
		if _, exists := existingPathMap[se.FilePath]; !exists {
			cID := existingCourseMap[se.CourseName]
			gID := groupMapping[fmt.Sprintf("%d_%s", cID, se.GroupName)]

			// 获取该分组当前的计数值
			currentSort := groupCurrentCount[gID]

			finalEpisodes = append(finalEpisodes, models.Episode{
				GroupID:  gID,
				Title:    se.Title,
				FilePath: se.FilePath,
				Sort:     currentSort,
			})

			// 更新计数值，保证同一次扫描中的下一集 sort 会递增
			groupCurrentCount[gID]++
		}
	}

	if len(finalEpisodes) > 0 {
		if err := tx.Create(&finalEpisodes); err != nil {
			if err := tx.Rollback(); err != nil {
				return response.InternalServerError(ctx, "E13", err)
			}
			return response.InternalServerError(ctx, "E14", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return response.InternalServerError(ctx, "E15", err)
	}

	return response.Ok(ctx, "扫描成功", map[string]int{
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
