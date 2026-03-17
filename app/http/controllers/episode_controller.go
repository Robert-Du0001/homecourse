package controllers

import (
	"homecourse/app/http/response"
	"homecourse/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"
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
	groupID := ctx.Request().RouteInt("group_id")

	var episodes []models.Episode

	if err := facades.Orm().Query().Where("group_id", groupID).
		Order("sort").
		Get(&episodes); err != nil {
		return response.InternalServerError(ctx, "E1", err)
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

	return ctx.Response().File(facades.Storage().Path(episode.FilePath))
}
