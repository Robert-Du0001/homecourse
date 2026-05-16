package controllers

import (
	"fmt"
	"homecourse/app/http/response"
	"homecourse/app/models"
	"path/filepath"
	"strings"

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

	// 查询剧集及其所属分组和课程
	var episode models.Episode
	if err := facades.Orm().Query().Find(&episode, episodeID); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	var group models.Group
	if err := facades.Orm().Query().Find(&group, episode.GroupID); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	var course models.Course
	if err := facades.Orm().Query().Find(&course, group.CourseID); err != nil {
		return response.InternalServerError(ctx, "E1", err)
	}

	// 构建路径: attachments/课程名/分组名/剧集名/
	dirPath := filepath.Join("attachments", course.Title, group.Name, episode.Title)
	originalName := file.GetClientOriginalName()

	// 处理同名文件：如果存在则在扩展名前追加（1）、（2）...
	finalName := originalName
	fullPath := filepath.Join(dirPath, finalName)
	counter := 1
	for facades.Storage().Exists(fullPath) {
		ext := filepath.Ext(originalName)
		baseName := strings.TrimSuffix(originalName, ext)
		finalName = fmt.Sprintf("%s（%d）%s", baseName, counter, ext)
		fullPath = filepath.Join(dirPath, finalName)
		counter++
	}

	// 保存文件（PutFileAs 会自动创建目录）
	attachmentPath, err := facades.Storage().PutFileAs(dirPath, file, finalName)
	if err != nil {
		return response.InternalServerError(ctx, "E4", err)
	}

	attachment := &models.Attachment{
		Name:      finalName,
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

	if _, err := facades.Orm().Query().Model(&models.Attachment{}).Where("id", id).Delete(); err != nil {
		return response.InternalServerError(ctx, "E1", err)
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
