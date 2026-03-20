package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"

	"homecourse/app/facades"
	"homecourse/app/http/controllers"
	"homecourse/app/http/middleware"
)

func Web() {
	// 处理静态文件
	facades.Route().Static("assets", "public/assets")
	facades.Route().Static("img", "public/img")
	facades.Route().StaticFile("favicon.svg", "public/favicon.svg")

	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		// 课程封面
		facades.Route().Get("covers/{path}", controllers.NewCourseController().ShowCover)
		// 课程视频
		facades.Route().Get("videos/{id}", controllers.NewEpisodeController().Play)
		// 剧集附件
		facades.Route().Get("attachments/{id}", controllers.NewAttachmentController().Show)
	})

	facades.Route().Fallback(func(ctx http.Context) http.Response {
		return ctx.Response().File("public/index.html")
	})
}
