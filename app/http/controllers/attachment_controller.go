package controllers

import (
	"homecourse/app/http/response"
	"homecourse/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type AttachmentController struct {
	// Dependent services
}

func NewAttachmentController() *AttachmentController {
	return &AttachmentController{
		// Inject services
	}
}

// 获取附件列表
func (r *AttachmentController) Index(ctx http.Context) http.Response {
	episodeId := ctx.Request().RouteInt("id")

	var attachments []models.Attachment

	if err := facades.Orm().Query().Where("episode_id", episodeId).
		Get(&attachments); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	return response.Ok(ctx, "获取成功", attachments)
}

// 创建附件
func (r *AttachmentController) Store(ctx http.Context) http.Response {
	episodeID := uint(ctx.Request().RouteInt("id"))

	file, err := ctx.Request().File("attachment_file")
	if err != nil {
		return response.InternalServerError(ctx, "E3", nil)
	}

	attachmentPath, err := facades.Storage().PutFile("/attachments", file)
	if err != nil {
		return response.InternalServerError(ctx, "E4", nil)
	}

	// 获取上传的附件名
	file, err = ctx.Request().File("attachment_file")
	if err != nil {
		return response.InternalServerError(ctx, "E5", nil)
	}

	attachment := &models.Attachment{
		Name:      file.GetClientOriginalName(),
		EpisodeID: episodeID,
		FilePath:  attachmentPath,
	}

	if err := facades.Orm().Query().Create(&attachment); err != nil {
		return response.InternalServerError(ctx, "E5", err)
	}

	return response.Ok(ctx, "附件上传成功", nil)
}

// 删除附件
func (r *AttachmentController) Destroy(ctx http.Context) http.Response {
	id := uint(ctx.Request().RouteInt("id"))

	var attachment models.Attachment

	// 获取文件路径
	if err := facades.Orm().Query().Select("file_path").Find(&attachment, id); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	tx, err := facades.Orm().Query().BeginTransaction()
	if err != nil {
		return response.InternalServerError(ctx, "E2", err)
	}

	if _, err := tx.Model(&models.Attachment{}).Where("id", id).Delete(); err != nil {
		if err := tx.Rollback(); err != nil {
			return response.InternalServerError(ctx, "E3", err)
		}
		return response.InternalServerError(ctx, "E4", err)
	}

	// 删除对应的附件
	if err := facades.Storage().Delete(attachment.FilePath); err != nil {
		if err := tx.Rollback(); err != nil {
			return response.InternalServerError(ctx, "E5", err)
		}
		return response.InternalServerError(ctx, "E6", err)
	}

	if err := tx.Commit(); err != nil {
		return response.InternalServerError(ctx, "E7", err)
	}

	return response.Ok(ctx, "附件删除成功", nil)
}

// 获取附件
func (r *AttachmentController) Show(ctx http.Context) http.Response {
	id := ctx.Request().RouteInt("id")

	var attachment models.Attachment

	// 获取文件路径
	if err := facades.Orm().Query().Select("file_path").Find(&attachment, id); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	return ctx.Response().File(facades.Storage().Path(attachment.FilePath))
}

// 统计附件
func (r *AttachmentController) Statistic(ctx http.Context) http.Response {
	total, err := facades.Orm().Query().Model(&models.Attachment{}).Count()
	if err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	return response.Ok(ctx, "获取成功", map[string]any{
		"total": total,
	})
}
