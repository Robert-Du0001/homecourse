package routes

import (
	"homecourse/app/http/controllers"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Web() {
	// 处理静态文件
	facades.Route().Static("assets", "public/assets")
	facades.Route().Static("img", "public/img")
	facades.Route().StaticFile("favicon.svg", "public/favicon.svg")

	facades.Route().Get("/media/{id}", controllers.NewEpisodeController().Play)

	facades.Route().Fallback(func(ctx http.Context) http.Response {
		return ctx.Response().File("public/index.html")
	})
}
