package controllers

import (
	"fmt"
	"homecourse/app/http/response"
	"homecourse/app/models"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/validation"
)

type EpisodeController struct {
	// Dependent services
}

func NewEpisodeController() *EpisodeController {
	return &EpisodeController{
		// Inject services
	}
}

// 获取剧集列表
func (r *EpisodeController) Index(ctx http.Context) http.Response {
	groupID := ctx.Request().RouteInt("group_id")

	var episodes []models.Episode

	if err := facades.Orm().Query().Where("group_id", groupID).
		Order("sort").
		Get(&episodes); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	return response.Ok(ctx, "课程集获取成功", episodes)
}

// 获取剧集信息
func (r *EpisodeController) Show(ctx http.Context) http.Response {
	episodeId := ctx.Request().Route("id")

	type course struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	}

	type group struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		CourseID uint   `json:"course_id"`

		Course *course `json:"course"`
	}

	type attachment struct {
		ID        uint   `json:"id"`
		EpisodeID uint   `json:"episode_id"`
		Name      string `json:"name"`
	}

	type episode struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
		// Duration uint   `json:"duration"`
		GroupID uint `json:"group_id"`

		Group       *group       `json:"group"`
		Attachments []attachment `json:"attachments"`
	}

	var resData episode

	if err := facades.Orm().Query().With("Group.Course").With("Attachments").FindOrFail(&resData, episodeId); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return response.BadRequest(ctx, "课程id不存在", nil)
		}

		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "获取课集详情成功", resData)
}

// 返回剧集文件
func (r *EpisodeController) Play(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")

	var episode models.Episode

	if err := facades.Orm().Query().FindOrFail(&episode, id); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return response.BadRequest(ctx, "课程id不存在", nil)
		}

		return response.InternalServerError(ctx, "E3", err)
	}

	return ctx.Response().File(facades.Storage().Path(episode.FilePath))
}

// 删除剧集
func (r *EpisodeController) Destroy(ctx http.Context) http.Response {
	episodeId := ctx.Request().RouteInt("id")

	tx, err := facades.Orm().Query().BeginTransaction()
	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if _, err := tx.Model(&models.Episode{}).Where("id", episodeId).
		Delete(); err != nil {
		if err := tx.Rollback(); err != nil {
			return response.InternalServerError(ctx, "E2", err)
		}
		return response.InternalServerError(ctx, "E3", err)
	}

	// 获取所有属于的附件
	type file struct {
		FilePath string
	}

	var files []file

	if err := tx.Model(&models.Attachment{}).
		Select("file_path").
		Where("episode_id", episodeId).
		Get(&files); err != nil {
		if err := tx.Rollback(); err != nil {
			return response.InternalServerError(ctx, "E4", err)
		}
		return response.InternalServerError(ctx, "E5", err)
	}

	// 删除这些附件
	if len(files) > 0 {
		for _, file := range files {
			if err := facades.Storage().Delete(file.FilePath); err != nil {
				if err := tx.Rollback(); err != nil {
					return response.InternalServerError(ctx, "E6", err)
				}
				return response.InternalServerError(ctx, "E7", err)
			}
		}
	}

	// 删除附件数据
	if _, err := tx.Model(&models.Attachment{}).
		Where("episode_id", episodeId).
		Delete(); err != nil {
		if err := tx.Rollback(); err != nil {
			return response.InternalServerError(ctx, "E8", err)
		}
		return response.InternalServerError(ctx, "E9", err)
	}

	if err := tx.Commit(); err != nil {
		return response.InternalServerError(ctx, "E10", err)
	}

	return response.Ok(ctx, "删除成功", nil)
}

// 修改剧集排序
func (r *EpisodeController) UpdateSort(ctx http.Context) http.Response {
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

	caseSql.WriteString("UPDATE episodes SET sort = CASE id ")
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

// 创建剧集
func (r *EpisodeController) Store(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"title":     "required|string|max_len:20",
		"file_path": "required|string",
		"group_id":  "required|uint",
	}, validation.Filters(map[string]string{
		"group_id": "uint",
	}))

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	var episode models.Episode

	if err := validator.Bind(&episode); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	// 获取分组下面的剧集总数，新增的剧集排序在最后
	if count, err := facades.Orm().Query().Model(&models.Episode{}).
		Where("group_id", episode.GroupID).
		Count(); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	} else {
		episode.Sort = uint(count)
	}

	if err := facades.Orm().Query().Create(&episode); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "剧集创建成功", nil)
}

// 修改剧集
func (r *EpisodeController) Update(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, ctx.Request().All(), map[string]string{
		"id":        "required|uint",
		"title":     "required|string|max_len:20",
		"file_path": "required|string",
		"group_id":  "required|uint",
	}, validation.Filters(map[string]string{
		"id":       "uint",
		"group_id": "uint",
	}))

	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	if validator.Fails() {
		return response.BadRequest(ctx, "参数错误", validator.Errors().All())
	}

	var episode models.Episode

	if err := validator.Bind(&episode); err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	if _, err := facades.Orm().Query().Model(&models.Episode{}).
		Where("id", episode.ID).
		Update(episode); err != nil {
		return response.InternalServerError(ctx, "E3", err)
	}

	return response.Ok(ctx, "剧集更新成功", nil)
}

// 统计剧集
func (r *EpisodeController) Statistic(ctx http.Context) http.Response {
	total, err := facades.Orm().Query().Model(&models.Episode{}).Count()
	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	return response.Ok(ctx, "获取成功", map[string]any{
		"total": total,
	})
}
