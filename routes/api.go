package routes

import (
	"github.com/goravel/framework/contracts/route"

	"homecourse/app/facades"
	"homecourse/app/http/controllers"
	"homecourse/app/http/middleware"
)

func Api() {
	userController := controllers.NewUserController()
	categoryController := controllers.NewCategoryController()
	courseController := controllers.NewCourseController()
	episodeController := controllers.NewEpisodeController()

	facades.Route().Prefix("api").Group(func(router route.Router) {
		router.Post("/user", userController.Store)
		router.Post("/user/token", userController.Login)

		router.Middleware(middleware.Auth()).Group(func(router route.Router) {
			router.Get("/user", userController.Show)

			// 课程分类相关路由
			router.Get("/categories", categoryController.Index)

			// 课程相关路由
			router.Get("/courses", courseController.Index)
			router.Get("/courses/{id}", courseController.Show)

			// 课程集相关路由
			router.Get("/episodes", episodeController.Index)
			router.Get("/episodes/{id}", episodeController.Show)
			router.Put("/episodes/scan", episodeController.Scan)

			// 管理员
			router.Middleware(middleware.Admin()).Group(func(router route.Router) {
				router.Get("/users", userController.Index)
			})
		})
	})
}
