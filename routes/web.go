package routes

import (
	"github.com/goravel/framework/contracts/http"

	"homecourse/app/facades"
	"homecourse/app/http/controllers"
)

func Web() {
	// 处理静态文件
	facades.Route().Static("assets", "public/assets")
	facades.Route().Static("img", "public/img")
	facades.Route().StaticFile("favicon.svg", "public/favicon.svg")

	// 课程封面
	facades.Route().Get("covers/{path}", controllers.NewCourseController().ShowCover)
	// 课程视频
	facades.Route().Get("videos/{id}", controllers.NewEpisodeController().Play)

	facades.Route().Fallback(func(ctx http.Context) http.Response {
		return ctx.Response().File("public/index.html")
	})
}
