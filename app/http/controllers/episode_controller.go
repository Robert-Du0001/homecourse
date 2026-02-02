package controllers

import (
	"homecourse/app/http/response"
	"homecourse/app/models"
	"slices"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/str"
)

type EpisodeController struct {
	// Dependent services
}

func NewEpisodeController() *EpisodeController {
	return &EpisodeController{
		// Inject services
	}
}

// 获取课程集列表
func (r *EpisodeController) Index(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"course_id": "required",
	})

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	type req struct {
		CourseID uint `form:"course_id"`
	}
	var requestData req

	if err := validator.Bind(&requestData); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	var episodes []models.Episode

	if err := facades.Orm().Query().Where("course_id", requestData.CourseID).
		Find(&episodes); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "课程集获取成功", episodes)
}

// 获取课程集信息
func (r *EpisodeController) Show(ctx http.Context) http.Response {
	episodeId := ctx.Request().Route("id")

	type course struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	}

	type episode struct {
		ID       uint   `json:"id"`
		Title    string `json:"title"`
		Duration uint   `json:"duration"`
		CourseId uint   `json:"course_id"`

		Course *course `json:"course"`
	}

	var resData episode

	if err := facades.Orm().Query().With("Course").FindOrFail(&resData, episodeId); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return response.BadRequest(ctx, "课程id不存在", nil)
		}

		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "获取课集详情成功", resData)
}

// 扫描课程
func (r *EpisodeController) Scan(ctx http.Context) http.Response {
	files, err := facades.Storage().AllFiles("courses")
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": err.Error(),
		})
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
				FilePath:   courseName + "/" + episodeName,
				CourseName: courseName,
			})
		}
	}

	// 2. 处理 Course：查重并批量插入
	type course struct {
		ID    uint
		Title string
	}

	var existingCourses []course
	courseArgs := make([]any, len(scannedCourseNames))
	for i, v := range scannedCourseNames {
		courseArgs[i] = v
	}
	if err := facades.Orm().Query().WhereIn("title", courseArgs).Get(&existingCourses); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	// 构建已有课程 Map，方便快速查找
	courseMap := make(map[string]uint)
	for _, c := range existingCourses {
		courseMap[c.Title] = c.ID
	}

	userID := ctx.Value(models.UserID).(uint)

	// 找出需要新创建的课程
	var newCourses []models.Course
	for _, name := range scannedCourseNames {
		if _, exists := courseMap[name]; !exists {
			newCourses = append(newCourses, models.Course{UserID: userID, Title: name})
		}
	}

	if len(newCourses) > 0 {
		if err := facades.Orm().Query().Create(&newCourses); err != nil {
			return response.InternalServerError(ctx, "E2", err)
		} else {
			// 插入后更新 Map，拿到新生成的 ID
			for _, nc := range newCourses {
				courseMap[nc.Title] = nc.ID
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
	if err := facades.Orm().Query().WhereIn("file_path", allPaths).Get(&existingEpisodes); err != nil {
		return response.InternalServerError(ctx, "E3", err)
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
				UserID:   userID,
				CourseID: courseMap[se.CourseName],
				Title:    se.Title,
				FilePath: se.FilePath,
			})
		}
	}

	if len(finalEpisodes) > 0 {
		if err := facades.Orm().Query().Create(&finalEpisodes); err != nil {
			return response.InternalServerError(ctx, "E4", err)
		}
	}

	return response.Ok(ctx, "扫描完成", map[string]int{
		"new_courses":  len(newCourses),
		"new_episodes": len(finalEpisodes),
	})
}

// 返回课程文件
func (r *EpisodeController) Play(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")

	var episode models.Episode

	if err := facades.Orm().Query().FindOrFail(&episode, id); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return response.BadRequest(ctx, "课程id不存在", nil)
		}

		return response.InternalServerError(ctx, "E3", err)
	}

	return ctx.Response().File(facades.Storage().Disk("course").Path(episode.FilePath))
}
